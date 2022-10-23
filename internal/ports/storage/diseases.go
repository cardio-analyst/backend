package storage

import "github.com/cardio-analyst/backend/internal/domain/models"

// DiseasesRepository encapsulates the logic of manipulations on the entity "Diseases" in the database.
type DiseasesRepository interface {
	// Update updates the data of the existing user diseases.
	Update(diseasesData models.Diseases) (err error)
	// Get searches for the user diseases information in the database according to the user id.
	//
	// By the time the method is used, it is assumed that the user diseases information definitely exists in the database,
	// so if it is not found, then the method returns an error.
	Get(userID uint64) (diseases *models.Diseases, err error)
}
