package store

import "database/sql"

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) Create() error {
	return nil
}
