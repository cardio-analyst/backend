package storage

import domain "github.com/cardio-analyst/backend/internal/gateway/domain/model"

// BasicIndicatorsRepository encapsulates the logic of manipulations on the entity "Basic Indicators" in the database.
type BasicIndicatorsRepository interface {
	// Save is a symbiosis of update and insert methods (upsert).
	//
	// If the basic indicators data is in the database, then the data of the existing basic indicators is updated, otherwise the data
	// of the new basic indicators is inserted.
	Save(basicIndicatorsData domain.BasicIndicators) (err error)
	// Get searches for the user basic indicators information in the database according to the basic indicators id.
	//
	// By the time the method is used, it is assumed that the user basic indicators information definitely exists in the database,
	// so if it is not found, then the method returns an error.
	Get(id, userID uint64) (basicIndicatorsData *domain.BasicIndicators, err error)
	// FindAll searches for user analyses in the database according to the user id. If user analyses are not found,
	// the method returns nil.
	FindAll(userID uint64) (basicIndicatorsDataList []*domain.BasicIndicators, err error)
}
