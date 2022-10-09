package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/domain/models"
	"github.com/cardio-analyst/backend/internal/ports/storage"
)

const userTable = "users"

// check whether Database structure implements the storage.UserStorage interface
var _ storage.UserStorage = (*Database)(nil)

func (d *Database) Save(userData models.User) error {
	userIDPlaceholder := "DEFAULT"
	if userData.ID != 0 {
		userIDPlaceholder = "$1"
	}

	updateSetStmtArgs := `
        first_name=$2,
		last_name=$3,
		middle_name=$4,
		region=$5,
		birth_date=$6,
        login=$7,
		email=$8`
	if userData.Password != "" {
		updateSetStmtArgs += `,
		password_hash=$9`
	}

	createUserQuery := fmt.Sprintf(`
		INSERT INTO %[1]v (id,
		                first_name,
						last_name,
						middle_name,
						region,
						birth_date,
						login,
						email,
						password_hash)
		VALUES (%[2]v, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (id) 
		    DO UPDATE SET 
		        %[3]v 
		    WHERE %[1]v.id=$1`,
		userTable, userIDPlaceholder, updateSetStmtArgs,
	)
	queryCtx := context.Background()

	// cast birthDate to query format
	birthDateCasted := pgtype.Date{Status: pgtype.Null}
	if err := birthDateCasted.Set(userData.BirthDate.Time); err != nil {
		return err
	}

	_, err := d.db.Exec(queryCtx, createUserQuery,
		userData.ID,
		userData.FirstName,
		userData.LastName,
		userData.MiddleName,
		userData.Region,
		birthDateCasted,
		userData.Login,
		userData.Email,
		userData.Password,
	)
	return err
}

func (d *Database) GetOneByCriteria(criteria models.UserCriteria) (*models.User, error) {
	whereStmt, whereStmtArgs := criteria.GetWhereStmtAndArgs()

	query := fmt.Sprintf(
		`
		SELECT id,
			   first_name,
			   last_name,
			   middle_name,
			   region,
			   birth_date,
			   login,
			   email,
			   password_hash
		FROM %v WHERE %v`,
		userTable, whereStmt,
	)
	queryCtx := context.Background()

	var user models.User
	if err := d.db.QueryRow(
		queryCtx, query, whereStmtArgs...,
	).Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.MiddleName,
		&user.Region,
		&user.BirthDate.Time,
		&user.Login,
		&user.Email,
		&user.Password,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (d *Database) FindByCriteria(criteria models.UserCriteria) ([]*models.User, error) {
	whereStmt, whereStmtArgs := criteria.GetWhereStmtAndArgs()

	query := fmt.Sprintf(
		`
		SELECT id,
			   first_name,
			   last_name,
			   middle_name,
			   region,
			   birth_date,
			   login,
			   email,
			   password_hash
		FROM %v WHERE %v`,
		userTable, whereStmt,
	)
	queryCtx := context.Background()

	rows, err := d.db.Query(queryCtx, query, whereStmtArgs...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	users := make([]*models.User, 0, 3)
	for rows.Next() {
		var user models.User

		if err = rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.MiddleName,
			&user.Region,
			&user.BirthDate.Time,
			&user.Login,
			&user.Email,
			&user.Password,
		); err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return users, nil
}
