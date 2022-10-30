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

const lifestyleTable = "lifestyles"

var _ storage.LifestyleRepository = (*lifestyleRepository)(nil)

type lifestyleRepository struct {
	storage *postgresStorage
}

func NewLifestyleRepository(storage *postgresStorage) *lifestyleRepository {
	return &lifestyleRepository{
		storage: storage,
	}
}

func (r *lifestyleRepository) Update(lifestyleData models.Lifestyle) error {
	query := fmt.Sprintf(`
		UPDATE %v
        SET 
            family_status=$2,
            events_participation=$3,
            physical_activity=$4,
            work_status=$5,
            significant_value_high=$6,
            significant_value_medium=$7,
            significant_value_low=$8
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
	)
	return err
}

func (r *lifestyleRepository) Get(userID uint64) (*models.Lifestyle, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
		       family_status,
		       events_participation,
		       physical_activity,
		       work_status,
		       significant_value_high,
		       significant_value_medium,
		       significant_value_low
		FROM %v WHERE user_id=$1`,
		lifestyleTable,
	)
	queryCtx := context.Background()

	var lifestyleData models.Lifestyle
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
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &lifestyleData, nil
}
