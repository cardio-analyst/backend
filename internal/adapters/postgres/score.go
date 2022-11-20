package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/cardio-analyst/backend/internal/domain/common"
	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/storage"
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

func (s *scoreRepository) GetCVERisk(cveRiskData models.CVERiskData) (uint64, error) {
	// TODO: replace 27-47 with custom CVE risk table builder, accepting CVE risk data
	var tableNameBuilder strings.Builder

	// table prefix depends on geo scope value (now it is only Russian scope, others will become later)
	// TODO: generate dynamically by scope
	tableNameBuilder.WriteString("very_high_risk")

	tableNameBuilder.WriteString("_")

	switch cveRiskData.Gender {
	case common.UserGenderMale:
		tableNameBuilder.WriteString("male")
	case common.UserGenderFemale:
		tableNameBuilder.WriteString("female")
	default:
		return 0, fmt.Errorf("unknown user gender: %v", cveRiskData.Gender)
	}

	tableNameBuilder.WriteString("_")

	if cveRiskData.Smoking {
		tableNameBuilder.WriteString("smoking")
	} else {
		tableNameBuilder.WriteString("not_smoking")
	}

	query := fmt.Sprintf(`
         SELECT risk_value 
         FROM %v 
         WHERE $1 BETWEEN systolic_blood_pressure_min AND systolic_blood_pressure_max AND 
               $2 BETWEEN non_hdl_cholesterol_min AND non_hdl_cholesterol_max AND 
               $3 BETWEEN age_min AND age_max`,
		tableNameBuilder.String(),
	)
	queryCtx := context.Background()

	var riskValue uint64
	if err := s.storage.conn.QueryRow(
		queryCtx, query, cveRiskData.SBPLevel, cveRiskData.TotalCholesterolLevel, cveRiskData.Age,
	).Scan(&riskValue); err != nil {
		return 0, err
	}

	return riskValue, nil
}
