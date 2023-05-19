package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
)

const regionUsersTable = "region_users"

type RegionUsersRepository struct {
	storage *Storage
}

func NewRegionUsersRepository(storage *Storage) *RegionUsersRepository {
	return &RegionUsersRepository{
		storage: storage,
	}
}

func (r *RegionUsersRepository) Increment(region string) error {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		INSERT INTO %[1]v (region) VALUES ($1)
		ON CONFLICT (region) 
		    DO UPDATE SET 
		        users_counter = %[1]v.users_counter + 1 
		    WHERE %[1]v.region=$1`,
		regionUsersTable,
	)

	_, err := r.storage.conn.Exec(queryCtx, query, region)
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		return sql.ErrNoRows
	}
	return err
}

type regionUsers struct {
	region       string `db:"region"`
	usersCounter int64  `db:"users_counter"`
}

func (r *RegionUsersRepository) All() (map[string]int64, error) {
	queryCtx := context.Background()

	query := fmt.Sprintf(`
		SELECT region, users_counter FROM %v`,
		regionUsersTable,
	)

	rows, err := r.storage.conn.Query(queryCtx, query)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	defer rows.Close()

	usersByRegions := make(map[string]int64)
	for rows.Next() {
		var tmp regionUsers
		if err = rows.Scan(&tmp.region, &tmp.usersCounter); err != nil {
			return nil, err
		}

		if tmp.region != "" {
			usersByRegions[tmp.region] = tmp.usersCounter
		}
	}

	if rows.Err() != nil {
		return nil, err
	}

	return usersByRegions, nil
}
