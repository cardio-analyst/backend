package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/gateway/domain/model"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"
)

const analysisTable = "analyses"

// check whether AnalysisRepository structure implements the storage.AnalysisRepository interface
var _ storage.AnalysisRepository = (*AnalysisRepository)(nil)

// AnalysisRepository implements storage.AnalysisRepository interface.
type AnalysisRepository struct {
	storage *Storage
}

func NewAnalysisRepository(storage *Storage) *AnalysisRepository {
	return &AnalysisRepository{
		storage: storage,
	}
}

func (r *AnalysisRepository) Save(analysisData model.Analysis) error {
	queryCtx := context.Background()

	analysisIDPlaceholder := "DEFAULT"
	if analysisData.ID != 0 {
		analysisIDPlaceholder = "$1"
	}

	updateSetStmtArgs := `
        high_density_cholesterol=$3,
		low_density_cholesterol=$4,
		triglycerides=$5,
		lipoprotein=$6,
		highly_sensitive_c_reactive_protein=$7,
        atherogenicity_coefficient=$8,
        creatinine=$9,
        atherosclerotic_plaques_presence=$10`

	query := fmt.Sprintf(`
		INSERT INTO %[1]v (id,
		                user_id,
						high_density_cholesterol,
						low_density_cholesterol,
						triglycerides,
						lipoprotein,
						highly_sensitive_c_reactive_protein,
						atherogenicity_coefficient,
						creatinine,
		                atherosclerotic_plaques_presence)
		VALUES (%[2]v, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (id) 
		    DO UPDATE SET 
		        %[3]v 
		    WHERE %[1]v.id=$1 AND %[1]v.user_id=$2`,
		analysisTable, analysisIDPlaceholder, updateSetStmtArgs,
	)

	_, err := r.storage.conn.Exec(queryCtx, query,
		analysisData.ID,
		analysisData.UserID,
		analysisData.HighDensityCholesterol,
		analysisData.LowDensityCholesterol,
		analysisData.Triglycerides,
		analysisData.Lipoprotein,
		analysisData.HighlySensitiveCReactiveProtein,
		analysisData.AtherogenicityCoefficient,
		analysisData.Creatinine,
		analysisData.AtheroscleroticPlaquesPresence,
	)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return sql.ErrNoRows
	}
	return err
}

func (r *AnalysisRepository) Get(id, userID uint64) (*model.Analysis, error) {
	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			high_density_cholesterol,
			low_density_cholesterol,
			triglycerides,
			lipoprotein,
			highly_sensitive_c_reactive_protein,
			atherogenicity_coefficient,
			creatinine,
			atherosclerotic_plaques_presence,
			created_at
		FROM %v
		WHERE id=$1 AND user_id=$2`,
		analysisTable,
	)
	queryCtx := context.Background()

	var analysisData model.Analysis
	if err := r.storage.conn.QueryRow(
		queryCtx, query, id, userID,
	).Scan(
		&analysisData.ID,
		&analysisData.UserID,
		&analysisData.HighDensityCholesterol,
		&analysisData.LowDensityCholesterol,
		&analysisData.Triglycerides,
		&analysisData.Lipoprotein,
		&analysisData.HighlySensitiveCReactiveProtein,
		&analysisData.AtherogenicityCoefficient,
		&analysisData.Creatinine,
		&analysisData.AtheroscleroticPlaquesPresence,
		&analysisData.CreatedAt.Time,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &analysisData, nil
}

func (r *AnalysisRepository) FindAll(userID uint64) ([]*model.Analysis, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			high_density_cholesterol,
			low_density_cholesterol,
			triglycerides,
			lipoprotein,
			highly_sensitive_c_reactive_protein,
			atherogenicity_coefficient,
			creatinine,
			atherosclerotic_plaques_presence,
			created_at
		FROM %v
		WHERE user_id=$1
		ORDER BY id DESC`,
		analysisTable,
	)

	rows, err := r.storage.conn.Query(queryCtx, query, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	analyses := make([]*model.Analysis, 0, 3)
	for rows.Next() {
		var analysis model.Analysis

		if err = rows.Scan(
			&analysis.ID,
			&analysis.UserID,
			&analysis.HighDensityCholesterol,
			&analysis.LowDensityCholesterol,
			&analysis.Triglycerides,
			&analysis.Lipoprotein,
			&analysis.HighlySensitiveCReactiveProtein,
			&analysis.AtherogenicityCoefficient,
			&analysis.Creatinine,
			&analysis.AtheroscleroticPlaquesPresence,
			&analysis.CreatedAt.Time,
		); err != nil {
			return nil, err
		}

		analyses = append(analyses, &analysis)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return analyses, nil
}
