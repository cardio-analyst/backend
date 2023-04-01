package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
	"github.com/cardio-analyst/backend/internal/gateway/ports/storage"

	"github.com/jackc/pgx/v4"
)

const sessionTable = "sessions"

// check whether sessionRepository structure implements the storage.SessionRepository interface
var _ storage.SessionRepository = (*sessionRepository)(nil)

// sessionRepository implements storage.SessionRepository interface.
type sessionRepository struct {
	storage *postgresStorage
}

func NewSessionRepository(storage *postgresStorage) *sessionRepository {
	return &sessionRepository{
		storage: storage,
	}
}

func (r *sessionRepository) Save(sessionData models.Session) error {
	query := fmt.Sprintf(`
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

	_, err := r.storage.conn.Exec(queryCtx, query,
		sessionData.UserID,
		sessionData.RefreshToken,
		sessionData.Whitelist,
	)
	return err
}

func (r *sessionRepository) Get(userID uint64) (*models.Session, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
			   refresh_token,
			   whitelist
		FROM %v WHERE user_id=$1`,
		sessionTable,
	)
	queryCtx := context.Background()

	var session models.Session
	if err := r.storage.conn.QueryRow(
		queryCtx, query, userID,
	).Scan(
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

func (r *sessionRepository) Find(userID uint64) (*models.Session, error) {
	query := fmt.Sprintf(
		`
		SELECT user_id,
			   refresh_token,
			   whitelist
		FROM %v WHERE user_id=$1`,
		sessionTable,
	)
	queryCtx := context.Background()

	var session models.Session
	if err := r.storage.conn.QueryRow(
		queryCtx, query, userID,
	).Scan(
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
