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

const sessionTable = "sessions"

// check whether Database structure implements the storage.SessionStorage interface
var _ storage.SessionStorage = (*Database)(nil)

func (d *Database) SaveSession(sessionData models.Session) error {
	createSessionQuery := fmt.Sprintf(`
		INSERT INTO %[1]v (user_id,
						   refresh_token,
						   whitelist)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) 
		    DO UPDATE SET 
		        refresh_token=$2,
		        whitelist=$3 
		    WHERE %[1]v.user_id=$1`,
		sessionTable,
	)
	queryCtx := context.Background()

	_, err := d.db.Exec(queryCtx, createSessionQuery,
		sessionData.UserID,
		sessionData.RefreshToken,
		sessionData.Whitelist,
	)
	return err
}

func (d *Database) GetSession(userID uint64) (*models.Session, error) {
	query := fmt.Sprintf(
		`
		SELECT id,
			   user_id,
			   refresh_token,
			   whitelist
		FROM %v WHERE user_id=$1`,
		sessionTable,
	)
	queryCtx := context.Background()

	var session models.Session
	if err := d.db.QueryRow(
		queryCtx, query, userID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.Whitelist,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &session, nil
}

func (d *Database) FindSession(userID uint64) (*models.Session, error) {
	query := fmt.Sprintf(
		`
		SELECT id,
			   user_id,
			   refresh_token,
			   whitelist
		FROM %v WHERE user_id=$1`,
		sessionTable,
	)
	queryCtx := context.Background()

	var session models.Session
	if err := d.db.QueryRow(
		queryCtx, query, userID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.RefreshToken,
		&session.Whitelist,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &session, nil
}
