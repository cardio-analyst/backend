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

var _ storage.DiseaseStorage = (*Database)(nil)

func (d *Database) SaveDisease(diseaseData models.Disease) (err error) {

	diseaseIDPlaceholder := "DEFAULT"
	if diseaseData.ID != 0 {
		diseaseIDPlaceholder = "$1"
	}

	updateSetStmtArgs := `
        user_id=$2,
		cvds_predisposition=$3,
		take_statins=$4,
		ckd=$5,
		arterial_hypertension=$6,
        cardiac_ischemia=$7,
		type_two_diabets=$8,
        infarction_or_stroke=$9,
        atherosclerosis=$10,
        other_cvds_diseases=$11`

	createDiseaseQuery := fmt.Sprintf(`
		INSERT INTO %[1]v (id,
		                user_id,
						cvds_predisposition,
						take_statins,
						ckd,
						arterial_hypertension,
                        cardiac_ischemia,
						type_two_diabets,
						infarction_or_stroke,
						atherosclerosis,
                        other_cvds_diseases)
		VALUES (%[2]v, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		ON CONFLICT (id) 
		    DO UPDATE SET 
		        %[3]v 
		    WHERE %[1]v.user_id=$2`,
		userTable, diseaseIDPlaceholder, updateSetStmtArgs,
	)
	queryCtx := context.Background()

	_, err = d.db.Exec(queryCtx, createDiseaseQuery,
		diseaseData.ID,
		diseaseData.UserID,
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

func (d *Database) GetDiseaseByUserId(userId uint) (*models.Disease, error) {
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
		FROM %v WHERE %[1]v.user_id =%v`,
		diseaseTable, userId,
	)
	queryCtx := context.Background()

	var disease models.Disease
	if err := d.db.QueryRow(queryCtx, query, userId).Scan(
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
