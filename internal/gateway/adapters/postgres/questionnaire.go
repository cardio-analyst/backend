package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

const questionnaireTable = "questionnaire"

var _ storage.QuestionnaireRepository = (*QuestionnaireRepository)(nil)

type QuestionnaireRepository struct {
	storage *Storage
}

func NewQuestionnaireRepository(storage *Storage) *QuestionnaireRepository {
	return &QuestionnaireRepository{
		storage: storage,
	}
}

func (r *QuestionnaireRepository) Get(userID uint64) (*model.Questionnaire, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
		       angina_score,
               adherence_drug_therapy,
               adherence_medical_support,
               adherence_lifestyle_mod
		FROM %v WHERE user_id=$1`,
		questionnaireTable,
	)
	queryCtx := context.Background()

	var questionnaire model.Questionnaire
	if err := r.storage.conn.QueryRow(
		queryCtx, query, userID,
	).Scan(
		&questionnaire.UserID,
		&questionnaire.AnginaScore,
		&questionnaire.AdherenceDrugTherapy,
		&questionnaire.AdherenceMedicalSupport,
		&questionnaire.AdherenceLifestyleMod,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			query = fmt.Sprintf(
				`INSERT INTO %v (user_id) VALUES ($1) RETURNING *`,
				questionnaireTable,
			)

			if err = r.storage.conn.QueryRow(queryCtx, query, userID).Scan(
				&questionnaire.UserID,
				&questionnaire.AnginaScore,
				&questionnaire.AdherenceDrugTherapy,
				&questionnaire.AdherenceMedicalSupport,
				&questionnaire.AdherenceLifestyleMod,
			); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &questionnaire, nil
}

func (r *QuestionnaireRepository) UpdateAnginaRose(questionnaire model.Questionnaire) error {
	query := fmt.Sprintf(`
		UPDATE %v
        SET 
            angina_score=$2
        WHERE user_id=$1`,
		questionnaireTable,
	)
	queryCtx := context.Background()

	_, err := r.storage.conn.Exec(queryCtx, query,
		questionnaire.UserID,
		questionnaire.AnginaScore,
	)
	return err
}

func (r *QuestionnaireRepository) UpdateTreatmentAdherence(questionnaire model.Questionnaire) error {
	query := fmt.Sprintf(`
		UPDATE %v
        SET 
            adherence_drug_therapy=$2,
            adherence_medical_support=$3,
            adherence_lifestyle_mod=$4
        WHERE user_id=$1`,
		questionnaireTable,
	)
	queryCtx := context.Background()

	_, err := r.storage.conn.Exec(queryCtx, query,
		questionnaire.UserID,
		questionnaire.AdherenceDrugTherapy,
		questionnaire.AdherenceMedicalSupport,
		questionnaire.AdherenceLifestyleMod,
	)
	return err
}
