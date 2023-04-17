package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

const lifestyleTable = "lifestyles"

var _ storage.LifestyleRepository = (*LifestyleRepository)(nil)

type LifestyleRepository struct {
	storage *Storage
}

func NewLifestyleRepository(storage *Storage) *LifestyleRepository {
	return &LifestyleRepository{
		storage: storage,
	}
}

func (r *LifestyleRepository) Update(lifestyleData model.Lifestyle) error {
	query := fmt.Sprintf(`
		UPDATE %v
        SET 
            family_status=$2,
            events_participation=$3,
            physical_activity=$4,
            work_status=$5,
            significant_value_high=$6,
            significant_value_medium=$7,
            significant_value_low=$8,
            angina_score=$9,
            adherence_drug_therapy=$10,
            adherence_medical_support=$11,
            adherence_lifestyle_mod=$12
        WHERE user_id=$1`,
		lifestyleTable,
	)
	queryCtx := context.Background()

	_, err := r.storage.conn.Exec(queryCtx, query,
		lifestyleData.UserID,
		lifestyleData.FamilyStatus,
		lifestyleData.EventsParticipation,
		lifestyleData.PhysicalActivity,
		lifestyleData.WorkStatus,
		lifestyleData.SignificantValueHigh,
		lifestyleData.SignificantValueMedium,
		lifestyleData.SignificantValueLow,
		lifestyleData.AnginaScore,
		lifestyleData.AdherenceDrugTherapy,
		lifestyleData.AdherenceMedicalSupport,
		lifestyleData.AdherenceLifestyleMod,
	)
	return err
}

func (r *LifestyleRepository) Get(userID uint64) (*model.Lifestyle, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
		       family_status,
		       events_participation,
		       physical_activity,
		       work_status,
		       significant_value_high,
		       significant_value_medium,
		       significant_value_low,
               angina_score,
               adherence_drug_therapy,
               adherence_medical_support,
               adherence_lifestyle_mod
		FROM %v WHERE user_id=$1`,
		lifestyleTable,
	)
	queryCtx := context.Background()

	var lifestyleData model.Lifestyle
	if err := r.storage.conn.QueryRow(
		queryCtx, query, userID,
	).Scan(
		&lifestyleData.UserID,
		&lifestyleData.FamilyStatus,
		&lifestyleData.EventsParticipation,
		&lifestyleData.PhysicalActivity,
		&lifestyleData.WorkStatus,
		&lifestyleData.SignificantValueHigh,
		&lifestyleData.SignificantValueMedium,
		&lifestyleData.SignificantValueLow,
		&lifestyleData.AnginaScore,
		&lifestyleData.AdherenceDrugTherapy,
		&lifestyleData.AdherenceMedicalSupport,
		&lifestyleData.AdherenceLifestyleMod,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			query = fmt.Sprintf(
				`INSERT INTO %v (user_id) VALUES ($1) RETURNING *`,
				lifestyleTable,
			)

			if err = r.storage.conn.QueryRow(queryCtx, query, userID).Scan(
				&lifestyleData.UserID,
				&lifestyleData.FamilyStatus,
				&lifestyleData.EventsParticipation,
				&lifestyleData.PhysicalActivity,
				&lifestyleData.WorkStatus,
				&lifestyleData.SignificantValueHigh,
				&lifestyleData.SignificantValueMedium,
				&lifestyleData.SignificantValueLow,
				&lifestyleData.AnginaScore,
				&lifestyleData.AdherenceDrugTherapy,
				&lifestyleData.AdherenceMedicalSupport,
				&lifestyleData.AdherenceLifestyleMod,
			); err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	return &lifestyleData, nil
}
