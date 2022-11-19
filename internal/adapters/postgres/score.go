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

var _ storage.ScoreRepository = (*scoreRepository)(nil)

// scoreRepository implements storage.ScoreRepository interface.
type scoreRepository struct {
	storage *postgresStorage
}

func NewScoreRepository(storage *postgresStorage) *scoreRepository {
	return &scoreRepository{
		storage: storage,
	}
}

func (s *scoreRepository) GetCveRisk(cveRiskData models.CveRiskData) (uint64, error) {
	statusSmoking := convertStatusSmoking(cveRiskData.Smoking)
	gender := convertGender(cveRiskData.Gender)

	cveRiskTable := fmt.Sprintf(`very_high_risk_%v_%v`, gender, statusSmoking)

	query := fmt.Sprintf(`
         SELECT risk_value FROM %v
                WHERE %v BETWEEN systolic_blood_pressure_min and systolic_blood_pressure_max and 
                %v BETWEEN non_hdl_cholesterol_min and non_hdl_cholesterol_max and 
                %v BETWEEN age_min and age_max`,
		cveRiskTable, cveRiskData.SbpLevel, cveRiskData.TotalCholesterolLevel, cveRiskData.Age,
	)
	queryCtx := context.Background()

	var riskValue uint64
	if err := s.storage.conn.QueryRow(
		queryCtx, query,
	).Scan(&riskValue); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, sql.ErrNoRows
		}
		return 0, err
	}

	return riskValue, nil
}

func convertStatusSmoking(smoking bool) string {
	if smoking {
		return "smoking"
	}
	return "not_smoking"
}

func convertGender(gender string) string {
	if gender == "мужской" {
		return "male"
	}
	return "female"
}
