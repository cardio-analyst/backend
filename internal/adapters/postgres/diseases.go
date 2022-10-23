package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

const diseasesTable = "diseases"

// check whether diseasesRepository structure implements the storage.DiseasesRepository interface
var _ storage.DiseasesRepository = (*diseasesRepository)(nil)

// diseasesRepository implements storage.DiseasesRepository interface.
type diseasesRepository struct {
	storage *postgresStorage
}

func NewDiseasesRepository(storage *postgresStorage) *diseasesRepository {
	return &diseasesRepository{
		storage: storage,
	}
}

func (r *diseasesRepository) Update(diseasesData models.Diseases) error {
	query := fmt.Sprintf(`
		UPDATE %v
        SET 
            cvd_predisposed=$2,
            takes_statins=$3,
            has_chronic_kidney_disease=$4,
            has_arterial_hypertension=$5,
            has_ischemic_heart_disease=$6,
            has_type_two_diabetes=$7,
            had_infarction_or_stroke=$8,
            has_atherosclerosis=$9,
            has_other_cvd=$10
        WHERE user_id=$1`,
		diseasesTable,
	)
	queryCtx := context.Background()

	_, err := r.storage.conn.Exec(queryCtx, query,
		diseasesData.UserID,
		diseasesData.CVDPredisposed,
		diseasesData.TakesStatins,
		diseasesData.HasChronicKidneyDisease,
		diseasesData.HasArterialHypertension,
		diseasesData.HasIschemicHeartDisease,
		diseasesData.HasTypeTwoDiabetes,
		diseasesData.HadInfarctionOrStroke,
		diseasesData.HasAtherosclerosis,
		diseasesData.HasOtherCVD,
	)
	return err
}

func (r *diseasesRepository) Get(userID uint64) (*models.Diseases, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
		       cvd_predisposed,
		       takes_statins,
		       has_chronic_kidney_disease,
		       has_arterial_hypertension,
		       has_ischemic_heart_disease,
		       has_type_two_diabetes,
		       had_infarction_or_stroke,
		       has_atherosclerosis,
		       has_other_cvd
		FROM %v WHERE user_id=$1`,
		diseasesTable,
	)
	queryCtx := context.Background()

	var diseasesData models.Diseases
	if err := r.storage.conn.QueryRow(
		queryCtx, query, userID,
	).Scan(
		&diseasesData.UserID,
		&diseasesData.CVDPredisposed,
		&diseasesData.TakesStatins,
		&diseasesData.HasChronicKidneyDisease,
		&diseasesData.HasArterialHypertension,
		&diseasesData.HasIschemicHeartDisease,
		&diseasesData.HasTypeTwoDiabetes,
		&diseasesData.HadInfarctionOrStroke,
		&diseasesData.HasAtherosclerosis,
		&diseasesData.HasOtherCVD,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &diseasesData, nil
}
