package storage

import (
	"github.com/cardio-analyst/backend/internal/gateway/domain/models"
)

// AnalysisRepository encapsulates the logic of manipulations on the entity "Analysis" in the database.
type AnalysisRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the analysis data is in the database, then the data of the existing analysis is updated, otherwise the data
	// of the new analysis is inserted.
	Save(analysisData models.Analysis) (err error)
	// Get searches for the user analysis information in the database according to the analysis id.
	//
	// By the time the method is used, it is assumed that the user analysis information definitely exists in the database,
	// so if it is not found, then the method returns an error.
	Get(id, userID uint64) (analysisData *models.Analysis, err error)
	// FindAll searches for user analyses in the database according to the user id. If user analyses are not found,
	// the method returns nil.
	FindAll(userID uint64) (analysisDataList []*models.Analysis, err error)
}
