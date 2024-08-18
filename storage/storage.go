package storage

import (
	"errors"
	task "http_server/models"
)

// type Task struct {
// 	readiness string
// 	result    string
// }

type Database struct {
	data map[string]task.Task
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]task.Task),
	}
}

func (db *Database) Get(key string) (*task.Task, error) {
	value, exists := db.data[key]

	if !exists {
		return nil, errors.New("key not found")
	}

	return &value, nil
}

func (db *Database) Post(key string, value task.Task) error {
	_, exists := db.data[key]

	if exists {
		return errors.New("key already exists")
	}

	db.data[key] = value

	return nil
}

func (db *Database) Put(key string, value task.Task) error {
	_, exists := db.data[key]

	if !exists {
		return errors.New("key not found")
	}

	db.data[key] = value

	return nil
}
