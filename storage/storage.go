package storage

import "errors"

type Database struct {
	data map[string]map[string]string
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]map[string]string),
	}
}

func (db *Database) Get(key string) (*map[string]string, error) {
	value, exists := db.data[key]

	if !exists {
		return nil, errors.New("key not found")
	}

	return &value, nil
}

func (db *Database) Post(key string, value map[string]string) error {
	_, exists := db.data[key]

	if exists {
		return errors.New("key already exists")
	}

	db.data[key] = value

	return nil
}

func (db *Database) Put(key string, value map[string]string) error {
	_, exists := db.data[key]

	if !exists {
		return errors.New("key not found")
	}

	db.data[key] = value

	return nil
}
