package storage

type RegionUsersRepository interface {
	Increment(region string) error
	All() (map[string]int64, error)
}
