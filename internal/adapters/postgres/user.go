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

var _ storage.UserStorage = (*Database)(nil)

func (d *Database) Create(userData models.User) error {
	createUserQuery := fmt.Sprintf(
		`
		INSERT INTO %v (first_name,
						   last_name,
						   middle_name,
						   region,
						   birth_date,
						   gender,
						   login,
						   email,
						   password_hash)
		VALUES ($1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9)`,
		userTable,
	)
	queryCtx := context.Background()

	// cast birthDate to query format
	birthDateCasted := pgtype.Date{Status: pgtype.Null}
	if err := birthDateCasted.Set(userData.BirthDate.Time); err != nil {
		return err
	}

	_, err := d.db.Exec(queryCtx, createUserQuery,
		userData.FirstName,
		userData.LastName,
		userData.MiddleName,
		userData.Region,
		birthDateCasted,
		userData.Gender,
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
			   gender,
			   login,
			   email
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
		&user.Gender,
		&user.Login,
		&user.Email,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &user, nil
}

func (d *Database) FindOneByCriteria(criteria models.UserCriteria) (*models.User, error) {
	whereStmt, whereStmtArgs := criteria.GetWhereStmtAndArgs()

	query := fmt.Sprintf(
		`
		SELECT id,
			   first_name,
			   last_name,
			   middle_name,
			   region,
			   birth_date,
			   gender,
			   login,
			   email
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
		&user.Gender,
		&user.Login,
		&user.Email,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
