package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const feedbackTable = "feedback"

type FeedbackRepository struct {
	storage *Storage
}

func NewFeedbackRepository(storage *Storage) *FeedbackRepository {
	return &FeedbackRepository{
		storage: storage,
	}
}

func (r *FeedbackRepository) Create(feedback model.Feedback) error {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		INSERT INTO %[1]v (user_id,
						user_first_name,
						user_last_name,
						user_middle_name,
						user_login,
						user_email,
						mark,
						message,
						version)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		feedbackTable,
	)

	_, err := r.storage.conn.Exec(queryCtx, query,
		feedback.UserID,
		feedback.UserFirstName,
		feedback.UserLastName,
		feedback.UserMiddleName,
		feedback.UserLogin,
		feedback.UserEmail,
		feedback.Mark,
		feedback.Message,
		feedback.Version,
	)
	return err
}

func (r *FeedbackRepository) FindAll(criteria model.FeedbackCriteria) ([]model.Feedback, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			user_first_name,
			user_last_name,
			user_middle_name,
			user_login,
			user_email,
			mark,
			message,
            version,
            viewed,
			created_at
		FROM %v`,
		feedbackTable,
	)

	if criteria.Limit > 0 {
		if criteria.Page == 0 {
			criteria.Page = 1
		}

		offset := (criteria.Page - 1) * criteria.Limit
		limit := criteria.Limit

		query += fmt.Sprintf(" LIMIT %v OFFSET %v", limit, offset)
	}

	rows, err := r.storage.conn.Query(queryCtx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	feedbacks := make([]model.Feedback, 0)
	for rows.Next() {
		var feedback model.Feedback

		if err = rows.Scan(
			&feedback.ID,
			&feedback.UserID,
			&feedback.UserFirstName,
			&feedback.UserLastName,
			&feedback.UserMiddleName,
			&feedback.UserLogin,
			&feedback.UserEmail,
			&feedback.Mark,
			&feedback.Message,
			&feedback.Version,
			&feedback.Viewed,
			&feedback.CreatedAt.Time,
		); err != nil {
			return nil, err
		}

		feedbacks = append(feedbacks, feedback)
	}

	if rows.Err() != nil {
		return nil, err
	}

	return feedbacks, nil
}

func (r *FeedbackRepository) Count(_ model.FeedbackCriteria) (int64, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`SELECT count(*) FROM %v`, feedbackTable)

	var feedbacksNum int64
	if err := r.storage.conn.QueryRow(
		queryCtx, query,
	).Scan(&feedbacksNum); err != nil {
		return 0, err
	}

	return feedbacksNum, nil
}

func (r *FeedbackRepository) One(id uint64) (*model.Feedback, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		SELECT 
			id,
			user_id,
			user_first_name,
			user_last_name,
			user_middle_name,
			user_login,
			user_email,
			mark,
			message,
            version,
            viewed,
			created_at
		FROM %v 
		WHERE id=$1`,
		feedbackTable,
	)

	var feedback model.Feedback
	if err := r.storage.conn.QueryRow(
		queryCtx, query, id,
	).Scan(
		&feedback.ID,
		&feedback.UserID,
		&feedback.UserFirstName,
		&feedback.UserLastName,
		&feedback.UserMiddleName,
		&feedback.UserLogin,
		&feedback.UserEmail,
		&feedback.Mark,
		&feedback.Message,
		&feedback.Version,
		&feedback.Viewed,
		&feedback.CreatedAt.Time,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return &feedback, nil
}

func (r *FeedbackRepository) UpdateViewed(id uint64, viewed bool) error {
	query := fmt.Sprintf(`
		UPDATE %v SET viewed=$2 WHERE id=$1`,
		feedbackTable,
	)
	queryCtx := context.Background()

	_, err := r.storage.conn.Exec(queryCtx, query, id, viewed)
	return err
}
