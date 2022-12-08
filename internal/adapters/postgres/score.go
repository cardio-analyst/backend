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

func (s *scoreRepository) GetCVERisk(data models.ScoreData) (uint64, error) {
	// TODO: replace 27-47 with custom CVE risk table builder, accepting CVE risk data
	var tableNameBuilder strings.Builder

	// table prefix depends on geo scope value (now it is only Russian scope, others will become later)
	// TODO: generate dynamically by scope
	tableNameBuilder.WriteString("very_high_risk")

	tableNameBuilder.WriteString("_")

	switch data.Gender {
	case common.UserGenderMale:
		tableNameBuilder.WriteString("male")
	case common.UserGenderFemale:
		tableNameBuilder.WriteString("female")
	default:
		return 0, fmt.Errorf("unknown user gender: %v", data.Gender)
	}

	tableNameBuilder.WriteString("_")

	if data.Smoking {
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
		queryCtx, query, data.SBPLevel, data.TotalCholesterolLevel, data.Age,
	).Scan(&riskValue); err != nil {
		return 0, err
	}

	return riskValue, nil
}

func (s *scoreRepository) GetIdealAge(cveRiskValue uint64, data models.ScoreData) (uint64, uint64, error) {
	// TODO: replace 27-47 with custom CVE risk table builder, accepting CVE risk data
	var tableNameBuilder strings.Builder

	// table prefix depends on geo scope value (now it is only Russian scope, others will become later)
	// TODO: generate dynamically by scope
	tableNameBuilder.WriteString("very_high_risk")

	tableNameBuilder.WriteString("_")

	switch data.Gender {
	case common.UserGenderMale:
		tableNameBuilder.WriteString("male")
	case common.UserGenderFemale:
		tableNameBuilder.WriteString("female")
	default:
		return 0, 0, fmt.Errorf("unknown user gender: %v", data.Gender)
	}

	tablesPrefix := tableNameBuilder.String()

	query := fmt.Sprintf(`
         WITH %[1]v AS (SELECT age_min, age_max, risk_value
                             FROM %[2]v
                             UNION
                             SELECT age_min, age_max, risk_value
                             FROM %[3]v) 
         SELECT age_min, age_max 
         FROM %[1]v 
         WHERE risk_value=$1 
         ORDER BY age_min DESC, age_max DESC 
         LIMIT 1`,
		tablesPrefix, fmt.Sprintf("%v_smoking", tablesPrefix), fmt.Sprintf("%v_not_smoking", tablesPrefix),
	)
	queryCtx := context.Background()

	var ageMin, ageMax uint64
	if err := s.storage.conn.QueryRow(
		queryCtx, query, cveRiskValue,
	).Scan(&ageMin, &ageMax); err != nil {
		return 0, 0, err
	}

	return ageMin, ageMax, nil
}
