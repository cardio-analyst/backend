package mongo

import (
	"context"

	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/cardio-analyst/backend/internal/pkg/model"
)

const counterNameUserID = "userID"

// UserRepository implements storage.UserRepository interface.
type UserRepository struct {
	storage *Storage
}

func NewUserRepository(storage *Storage) *UserRepository {
	return &UserRepository{
		storage: storage,
	}
}

func (r *UserRepository) Save(ctx context.Context, user model.User) error {
	if user.ID == 0 {
		var err error
		user.ID, err = r.storage.getNextValue(ctx, counterNameUserID)
		if err != nil {
			return err
		}
	}

	if user.Role == "" {
		user.Role = model.UserRoleCustomer
	}

	filter := bson.M{"id": user.ID}
	update := bson.D{{"$set", user}}
	opts := options.Update().SetUpsert(true)

	result, err := r.storage.users.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if result.MatchedCount != 0 {
		log.Debug("matched and replaced an existing user")
	} else if result.UpsertedCount != 0 {
		log.Debugf("inserted a new user with ID %v (mongo ID %v)", user.ID, result.UpsertedID)
	}

	return nil
}

func (r *UserRepository) GetOneByCriteria(ctx context.Context, criteria model.UserCriteria) (model.User, error) {
	filter := userFilterFromCriteria(criteria)

	var user model.User
	if err := r.storage.users.FindOne(ctx, filter).Decode(&user); err != nil {
		if err == mongo.ErrNoDocuments {
			return user, model.ErrUserNotFound
		}
		return user, err
	}

	return user, nil
}

func (r *UserRepository) FindAllByCriteria(ctx context.Context, criteria model.UserCriteria) ([]model.User, error) {
	filter := userFilterFromCriteria(criteria)

	var findOpts *options.FindOptions
	if criteria.Limit > 0 {
		if criteria.Page == 0 {
			criteria.Page = 1
		}

		skip := (criteria.Page - 1) * criteria.Limit
		limit := criteria.Limit

		findOpts = &options.FindOptions{
			Skip:  &skip,
			Limit: &limit,
		}
	}

	cursor, err := r.storage.users.Find(ctx, filter, findOpts)
	if err != nil {
		return nil, err
	}

	var users []model.User
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}

func userFilterFromCriteria(criteria model.UserCriteria) bson.M {
	filter := make([]bson.M, 0)

	if criteria.ID != 0 {
		filter = append(filter, bson.M{"id": criteria.ID})
	}
	if criteria.Login != "" {
		filter = append(filter, bson.M{"login": criteria.Login})
	}
	if criteria.Email != "" {
		filter = append(filter, bson.M{"email": criteria.Email})
	}
	if criteria.PasswordHash != "" {
		filter = append(filter, bson.M{"password_hash": criteria.PasswordHash})
	}
	if criteria.Region != "" {
		filter = append(filter, bson.M{"region": criteria.Region})
	}
	if !criteria.BirthDateFrom.IsZero() {
		filter = append(filter, bson.M{"birth_date": bson.M{
			"$gte": primitive.NewDateTimeFromTime(criteria.BirthDateFrom.Time),
		}})
	}
	if !criteria.BirthDateTo.IsZero() {
		filter = append(filter, bson.M{"birth_date": bson.M{
			"$lte": primitive.NewDateTimeFromTime(criteria.BirthDateTo.Time),
		}})
	}

	switch len(filter) {
	case 0:
		return bson.M{}
	case 1:
		return filter[0]
	default:
		operator := "$and"
		if criteria.CriteriaSeparator == model.CriteriaSeparatorOR {
			operator = "$or"
		}

		return bson.M{operator: filter}
	}
}

func (r *UserRepository) Count(ctx context.Context, criteria model.UserCriteria) (int64, error) {
	filter := userFilterFromCriteria(criteria)

	return r.storage.users.CountDocuments(ctx, filter)
}
