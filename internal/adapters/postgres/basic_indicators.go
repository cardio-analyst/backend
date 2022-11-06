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

const basicIndicatorsTable = "basic_indicators"

// check whether basicIndicatorsRepository structure implements the storage.BasicIndicatorsRepository interface
var _ storage.BasicIndicatorsRepository = (*basicIndicatorsRepository)(nil)

// basicIndicatorsRepository implements storage.BasicIndicatorsRepository interface.
type basicIndicatorsRepository struct {
	storage *postgresStorage
}

func NewBasicIndicatorsRepository(storage *postgresStorage) *basicIndicatorsRepository {
	return &basicIndicatorsRepository{
		storage: storage,
	}
}

func (r *basicIndicatorsRepository) Save(basicIndicatorsData models.BasicIndicators) error {
	queryCtx := context.Background()

	basicIndicatorsIDPlaceholder := "DEFAULT"
	if basicIndicatorsData.ID != 0 {
		basicIndicatorsIDPlaceholder = "$1"
	}

	updateSetStmtArgs := `
        weight=$3,
		height=$4,
		body_mass_index=$5,
		waist_size=$6,
		gender=$7,
        sbp_level=$8,
        smoking=$9,
        total_cholesterol_level=$10,
        cv_events_risk_value=$11,
        ideal_cardiovascular_age=$12`

	query := fmt.Sprintf(`
		INSERT INTO %[1]v (id,
		                user_id,
						weight,
						height,
						body_mass_index,
						waist_size,
						gender,
						sbp_level,
						smoking,
						total_cholesterol_level,
						cv_events_risk_value,
						ideal_cardiovascular_age)
		VALUES (%[2]v, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) 
		    DO UPDATE SET 
		        %[3]v 
		    WHERE %[1]v.id=$1 AND %[1]v.user_id=$2`,
		basicIndicatorsTable, basicIndicatorsIDPlaceholder, updateSetStmtArgs,
	)

	_, err := r.storage.conn.Exec(queryCtx, query,
		basicIndicatorsData.ID,
		basicIndicatorsData.UserID,
		basicIndicatorsData.Weight,
		basicIndicatorsData.Height,
		basicIndicatorsData.BodyMassIndex,
		basicIndicatorsData.WaistSize,
		basicIndicatorsData.Gender,
		basicIndicatorsData.SBPLevel,
		basicIndicatorsData.Smoking,
		basicIndicatorsData.TotalCholesterolLevel,
		basicIndicatorsData.CVEventsRiskValue,
		basicIndicatorsData.IdealCardiovascularAge,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return sql.ErrNoRows
	}
	return err
}

func (r *basicIndicatorsRepository) Get(id, userID uint64) (*models.BasicIndicators, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			weight,
			height,
			body_mass_index,
			waist_size,
			gender,
			sbp_level,
			smoking,
			total_cholesterol_level,
			cv_events_risk_value,
			ideal_cardiovascular_age,
			created_at
		FROM %v
		WHERE id=$1 AND user_id=$2`,
		basicIndicatorsTable,
	)
	queryCtx := context.Background()

	var basicIndicatorsData models.BasicIndicators
	if err := r.storage.conn.QueryRow(
		queryCtx, query, id, userID,
	).Scan(
		&basicIndicatorsData.ID,
		&basicIndicatorsData.UserID,
		&basicIndicatorsData.Weight,
		&basicIndicatorsData.Height,
		&basicIndicatorsData.BodyMassIndex,
		&basicIndicatorsData.WaistSize,
		&basicIndicatorsData.Gender,
		&basicIndicatorsData.SBPLevel,
		&basicIndicatorsData.Smoking,
		&basicIndicatorsData.TotalCholesterolLevel,
		&basicIndicatorsData.CVEventsRiskValue,
		&basicIndicatorsData.IdealCardiovascularAge,
		&basicIndicatorsData.CreatedAt.Time,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &basicIndicatorsData, nil
}

func (r *basicIndicatorsRepository) FindAll(userID uint64) ([]*models.BasicIndicators, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			weight,
			height,
			body_mass_index,
			waist_size,
			gender,
			sbp_level,
			smoking,
			total_cholesterol_level,
			cv_events_risk_value,
			ideal_cardiovascular_age,
			created_at
		FROM %v
		WHERE user_id=$1
		ORDER BY id DESC`,
		basicIndicatorsTable,
	)

	rows, err := r.storage.conn.Query(queryCtx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	analyses := make([]*models.BasicIndicators, 0, 3)
	for rows.Next() {
		var basicIndicators models.BasicIndicators

		if err = rows.Scan(
			&basicIndicators.ID,
			&basicIndicators.UserID,
			&basicIndicators.Weight,
			&basicIndicators.Height,
			&basicIndicators.BodyMassIndex,
			&basicIndicators.WaistSize,
			&basicIndicators.Gender,
			&basicIndicators.SBPLevel,
			&basicIndicators.Smoking,
			&basicIndicators.TotalCholesterolLevel,
			&basicIndicators.CVEventsRiskValue,
			&basicIndicators.IdealCardiovascularAge,
			&basicIndicators.CreatedAt.Time,
		); err != nil {
			return nil, err
		}

		analyses = append(analyses, &basicIndicators)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return analyses, nil
}
