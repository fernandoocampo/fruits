package fruits

import (
	"context"
	"errors"
	"time"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
)

// Repository defines portout behavior to send fruit data to external platforms.
type Repository interface {
	FindByID(ctx context.Context, fruitID repository.FruitID) (*repository.Fruit, error)
	Save(ctx context.Context, fruit repository.NewFruit) (repository.FruitID, error)
	SearchWithFilters(ctx context.Context, filter repository.FruitFilter) (repository.FindFruitsResult, error)
	DatasetStatus(ctx context.Context) (repository.FruitDatasetStatus, error)
}

// Service implements fruit management logic.
type Service struct {
	fruitRepository Repository
	logger          *loggers.Logger
}

var ErrDataAccess = errors.New("something went wrong accessing db")

// NewService creates a new application service.
func NewService(fruitRepository Repository, logger *loggers.Logger) *Service {
	return &Service{
		fruitRepository: fruitRepository,
		logger:          logger,
	}
}

// GetFruitWithID get the fruit with the given id.
func (s *Service) GetFruitWithID(ctx context.Context, fruitID int64) (*Fruit, error) {
	s.logger.Debug(
		"getting fruit with id",
		loggers.Fields{
			"method":  "Service.GetFruitWithID",
			"fruitID": fruitID,
		},
	)

	result, err := s.fruitRepository.FindByID(ctx, repository.FruitID(fruitID))
	if err != nil {
		s.logger.Error(
			"something went wrong trying to get a fruit",
			loggers.Fields{
				"method":  "Service.GetFruitWithID",
				"fruitID": fruitID,
				"error":   err,
			},
		)

		return nil, ErrDataAccess
	}

	fruit := transformFruitPortOuttoFruit(result)

	s.logger.Debug(
		"fruit result",
		loggers.Fields{
			"method": "Service.GetFruitWithID",
			"fruit":  fruit,
		},
	)

	return fruit, nil
}

// Create creates a fruit.
func (s *Service) Create(ctx context.Context, newfruit NewFruit) (int64, error) {
	s.logger.Debug(
		"creating fruit",
		loggers.Fields{
			"method":   "Service.Create",
			"newfruit": newfruit,
		},
	)

	fruitid, err := s.fruitRepository.Save(ctx, newfruit.ToFruitPortOut())
	if err != nil {
		s.logger.Error(
			"something goes wrong creating a new fruit",
			loggers.Fields{
				"method": "Service.Create",
				"fruit":  newfruit,
			},
		)

		return 0, ErrDataAccess
	}

	s.logger.Info(
		"fruit was created successfully",
		loggers.Fields{
			"method": "Service.Create",
			"fruit":  newfruit,
		},
	)

	return repository.FruitIDValue(fruitid), nil
}

// SearchFruits search fruits who match the given filters.
func (s *Service) SearchFruits(ctx context.Context, givenFilter SearchFruitFilter) (*SearchFruitsResult, error) {
	s.logger.Debug(
		"searching fruits",
		loggers.Fields{
			"method": "Service.SearchFruits",
			"filter": givenFilter,
		},
	)

	filters := givenFilter.toRepositoryFilters()

	repoResult, err := s.fruitRepository.SearchWithFilters(ctx, filters)
	if err != nil {
		s.logger.Error(
			"something goes wrong searching fruits",
			loggers.Fields{
				"method": "Service.SearchFruits",
				"filter": givenFilter,
			},
		)

		return nil, ErrDataAccess
	}

	result := toSearchFruitsResult(repoResult)

	return &result, nil
}

// DatasetStatus check the status of the fruit dataset.
func (s *Service) DatasetStatus(ctx context.Context) DatasetStatus {
	s.logger.Debug(
		"checking dataset status",
		loggers.Fields{
			"method": "Service.DatasetStatus",
		},
	)

	currentState, err := s.fruitRepository.DatasetStatus(ctx)
	if err != nil {
		s.logger.Error(
			"something goes wrong checking the fruit dataset",
			loggers.Fields{
				"method": "Service.DatasetStatus",
				"error":  err,
			},
		)

		return DatasetStatus{
			Timestamp: time.Now().Unix(),
			Status:    DatasetStateError,
			Message:   "fruit repository is not available",
		}
	}

	result := toDatasetStatus(currentState)

	return result
}
