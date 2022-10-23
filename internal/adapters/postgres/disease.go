package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/storage"
	"github.com/jackc/pgx/v4"
)

const diseaseTable = "diseases"

var _ storage.DiseaseRepository = (*diseaseRepository)(nil)

// userRepository implements storage.UserRepository interface.
type diseaseRepository struct {
	storage *postgresStorage
}

func NewDiseaseRepository(storage *postgresStorage) *diseaseRepository {
	return &diseaseRepository{
		storage: storage,
	}
}

func (r *diseaseRepository) Save(diseaseData models.Disease) (err error) {

	updateSetStmtArgs := `
		cvds_predisposition=$1,
		take_statins=$2,
		ckd=$3,
		arterial_hypertension=$4,
        cardiac_ischemia=$5,
		type_two_diabets=$6,
        infarction_or_stroke=$7,
        atherosclerosis=$8,
        other_cvds_diseases=$9`

	createDiseaseQuery := fmt.Sprintf(`
		UPDATE %v SET %v WHERE user_id = %v`,
		diseaseTable, updateSetStmtArgs, diseaseData.UserID,
	)
	queryCtx := context.Background()

	_, err = r.storage.conn.Exec(queryCtx, createDiseaseQuery,
		diseaseData.CvdsPredisposition,
		diseaseData.TakeStatins,
		diseaseData.Ckd,
		diseaseData.ArterialHypertension,
		diseaseData.CardiacIschemia,
		diseaseData.TypeTwoDiabets,
		diseaseData.InfarctionOrStroke,
		diseaseData.Atherosclerosis,
		diseaseData.OtherCvdsDiseases,
	)
	return err
}

func (r *diseaseRepository) GetByUserId(userId uint64) (*models.Disease, error) {
	query := fmt.Sprintf(
		`
		SELECT id,
               user_id,
               cvds_predisposition,
               take_statins,
               ckd,
               arterial_hypertension,
               cardiac_ischemia,
	           type_two_diabets,
			   infarction_or_stroke,
               atherosclerosis,
               other_cvds_diseases
		FROM %v WHERE user_id = %v`,
		diseaseTable, userId,
	)
	queryCtx := context.Background()

	var disease models.Disease
	if err := r.storage.conn.QueryRow(queryCtx, query).Scan(
		&disease.ID,
		&disease.UserID,
		&disease.CvdsPredisposition,
		&disease.TakeStatins,
		&disease.Ckd,
		&disease.ArterialHypertension,
		&disease.CardiacIschemia,
		&disease.TypeTwoDiabets,
		&disease.InfarctionOrStroke,
		&disease.Atherosclerosis,
		&disease.OtherCvdsDiseases,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &disease, nil
}
