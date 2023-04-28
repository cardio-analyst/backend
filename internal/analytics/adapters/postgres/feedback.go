package postgres

import (
	"context"
	"fmt"

	"github.com/cardio-analyst/backend/internal/analytics/domain/model"
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
						message)
		VALUES (DEFAULT, $1, $2, $3, $4, $5, $6, $7, $8)`,
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
	)
	return err
}
