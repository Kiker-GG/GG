package storage

import (
	"errors"
	model "http_server/models"
)

// type Task struct {
//  id string
// 	readiness string
// 	result    string
// }

type Database struct {
	data_tasks   map[string]model.Task
	data_users   map[string]model.User
	data_session map[string]model.Session
}

func NewDatabase() *Database {
	return &Database{
		data_tasks:   make(map[string]model.Task),
		data_users:   make(map[string]model.User),
		data_session: make(map[string]model.Session),
	}
}

func (db *Database) Get(key string) (*model.Task, error) {
	value, exists := db.data_tasks[key]

	if !exists {
		return nil, errors.New("key not found")
	}

	return &value, nil
}

func (db *Database) Post(value model.Task) error {
	key := value.ID
	_, exists := db.data_tasks[key]

	if exists {
		return errors.New("key already exists")
	}

	db.data_tasks[key] = value

	return nil
}

func (db *Database) Put(value model.Task) error {
	key := value.ID
	_, exists := db.data_tasks[key]

	if !exists {
		return errors.New("key not found")
	}

	db.data_tasks[key] = value

	return nil
}

func (db *Database) Post_user(value model.User) error {
	key := value.ID

	_, exists := db.data_users[key]

	if !exists {
		return errors.New("key not found")
	}

	db.data_users[key] = value

	return nil
}

func (db *Database) Get_user(value model.User) (*model.User, error) {
	for _, elem := range db.data_users {
		if elem.Login == value.Login && elem.Password == value.Password {
			return &elem, nil
		}
	}

	return nil, errors.New("user not found")

}

func (db *Database) Post_session(value model.Session) error {
	key := value.Session_id

	_, exists := db.data_session[key]

	if !exists {
		return errors.New("key not found")
	}

	db.data_session[key] = value

	return nil
}

func (db *Database) Get_session(key string) error {
	_, exists := db.data_session[key]

	if !exists {
		return errors.New("key not found")
	}

	return nil
}
