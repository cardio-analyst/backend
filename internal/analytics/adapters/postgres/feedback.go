package postgres

import (
	"context"
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
		INSERT INTO %[1]v (id,
		                user_id,
						user_first_name,
						user_last_name,
						user_middle_name,
						user_login,
						user_email,
						mark,
						message,
						version)
		VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8, $9)`,
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

func (r *FeedbackRepository) FindAll() ([]model.Feedback, error) {
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
			created_at
		FROM %v`,
		feedbackTable,
	)

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
