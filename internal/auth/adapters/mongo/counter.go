package mongo

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Counter struct {
	ID    string `bson:"_id,omitempty"`
	Value uint64 `json:"value"`
}

func initCounter(ctx context.Context, counters *mongo.Collection, counterName string) error {
	filter := bson.D{
		{"_id", bson.M{"$eq": counterName}},
	}

	var c Counter
	err := counters.FindOne(ctx, filter).Decode(&c)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return nil
		} else if errors.Is(mongo.ErrNoDocuments, err) {
			_, err = counters.InsertOne(ctx, bson.M{
				"_id":   counterName,
				"value": 0,
			})
		}
	}
	return err
}

func (s *Storage) getNextValue(ctx context.Context, counterName string) (uint64, error) {
	filter := bson.D{
		{"_id", bson.M{"$eq": counterName}},
	}
	update := bson.D{
		{"$inc", bson.D{{"value", 1}}},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var c Counter
	if err := s.counters.FindOneAndUpdate(ctx, filter, update, opts).Decode(&c); err != nil {
		return 0, err
	}

	return c.Value, nil
}
