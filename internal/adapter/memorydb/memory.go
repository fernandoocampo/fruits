package memorydb

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
)

// UpdateEntityError if there is an error
// updating an entity this is the error
// to return.
type UpdateEntityError struct {
	Message string
}

// Repository is the repository handler for the map storage.
type Repository struct {
	lastID  int64
	nextID  int64
	ids     []int64
	storage map[int64]interface{}
	locker  sync.Mutex
	logger  *loggers.Logger
}

var errNotFoundEntity = errors.New("given entity doesn't exist")

// NewRepository creates a new  repository that will use a rdb.
func NewRepository(logger *loggers.Logger) *Repository {
	newRepo := Repository{
		lastID:  0,
		nextID:  1,
		storage: make(map[int64]interface{}),
		logger:  logger,
	}

	return &newRepo
}

func (u UpdateEntityError) Error() string {
	return u.Message
}

// Save store the given entity and return its id.
func (u *Repository) Save(ctx context.Context, entityID int64, entity interface{}) error {
	u.logger.Debug(
		"storing entity",
		loggers.Fields{
			"method": "memory.Repository.Save",
			"entity": entity,
		},
	)

	u.locker.Lock()
	{
		u.storage[entityID] = entity
		u.ids = append(u.ids, entityID)
	}
	u.locker.Unlock()

	return nil
}

// Update update the given entity.
func (u *Repository) Update(ctx context.Context, entityID int64, entity interface{}) error {
	u.logger.Debug(
		"updating entity",
		loggers.Fields{
			"method": "memory.Repository.Update",
			"entity": entity,
		},
	)

	if entityID == 0 {
		u.logger.Error(
			"cannot update given entity, because it doesn't contain a valid id",
			loggers.Fields{
				"method": "memory.Repository.Update",
				"entity": entity,
			},
		)

		return UpdateEntityError{
			Message: fmt.Sprintf(
				"cannot update given entity %v, because it doesn't contain a valid id",
				entity,
			),
		}
	}

	u.locker.Lock()
	defer u.locker.Unlock()

	if u.entityNotExist(entityID) {
		return errNotFoundEntity
	}

	u.storage[entityID] = entity

	return nil
}

// FindByID finds a entity with the given id in the memory storate of this dry run database.
func (u *Repository) FindByID(ctx context.Context, entityID int64) (interface{}, error) {
	u.logger.Debug(
		"reading entity",
		loggers.Fields{
			"method":    "memory.Repository.FindByID",
			"entity_id": entityID,
		},
	)

	var entity interface{}

	var recordExist bool

	u.locker.Lock()
	{
		entity, recordExist = u.storage[entityID]
	}
	u.locker.Unlock()

	if !recordExist {
		return entity, nil
	}

	u.logger.Debug(
		"entity found",
		loggers.Fields{
			"method":    "memory.Repository.FindByID",
			"entity_id": entityID,
			"entity":    entity,
		},
	)

	return entity, nil
}

// FindAll return all entities.
func (u *Repository) FindAll(ctx context.Context, start, count int) ([]interface{}, error) {
	u.logger.Debug(
		"reading all entities",
		loggers.Fields{
			"method": "memory.Repository.FindAll",
			"start":  start,
			"count":  count,
		},
	)

	result := make([]interface{}, 0)

	u.locker.Lock()
	{
		if len(u.ids) == 0 || start > len(u.ids) {
			start = 1
			count = 0
		}

		newcount := start + count - 1

		if len(u.ids) < count {
			newcount = len(u.ids)
		}

		ids := u.ids[start-1 : newcount]
		for _, id := range ids {
			v, ok := u.storage[id]
			if !ok {
				continue
			}
			result = append(result, v)
		}
	}

	u.locker.Unlock()

	u.logger.Debug(
		"found entities",
		loggers.Fields{
			"method": "memory.Repository.FindAll",
			"start":  start,
			"count":  count,
			"result": result,
		},
	)

	return result, nil
}

// Count counts records in the memory repo.
func (u *Repository) Count() int {
	var count int

	u.locker.Lock()
	{
		count = len(u.storage)
	}
	u.locker.Unlock()

	return count
}

// entityNotExist return true if the given entity id does not exist in the repository.
func (u *Repository) entityNotExist(entityID int64) bool {
	_, ok := u.storage[entityID]

	return !ok
}

// NewID create and return a new id for this repository.
func (u *Repository) NewID() int64 {
	u.locker.Lock()
	defer u.locker.Unlock()
	newID := u.nextID
	u.lastID = newID
	u.nextID++

	return newID
}

// UpdateID update the fruit id sequence on this repository.
func (u *Repository) UpdateID(newID int64) {
	u.locker.Lock()
	{
		u.lastID = newID
		u.nextID++
	}
	u.locker.Unlock()
}
