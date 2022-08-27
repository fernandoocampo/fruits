package fruits

import (
	"context"
	"errors"
	"fmt"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/go-kit/kit/endpoint"
)

// Endyear is a wrapper for endyear.
type Endyear struct {
	GetFruitWithIDEndpoint endpoint.Endpoint
	CreateFruitEndpoint    endpoint.Endpoint
	SearchFruitsEndpoint   endpoint.Endpoint
	GetStatusEndpoint      endpoint.Endpoint
}

var (
	errInvalidFruitID      = errors.New("invalid fruit id")
	errInvalidFruitFilters = errors.New("invalid fruit filters")
	errInvalidNewFruitType = errors.New("invalid new fruit type")
)

// NewEndyear Create the endyear for fruits-micro application.
func NewEndyear(service FruitService, logger *loggers.Logger) Endyear {
	return Endyear{
		GetFruitWithIDEndpoint: MakeGetFruitWithIDEndpoint(service, logger),
		CreateFruitEndpoint:    MakeCreateFruitEndpoint(service, logger),
		SearchFruitsEndpoint:   MakeSearchFruitsEndpoint(service, logger),
		GetStatusEndpoint:      MakeGetStatusEndpoint(service, logger),
	}
}

// MakeGetFruitWithIDEndpoint create endpoint for get a fruit with ID service.
func MakeGetFruitWithIDEndpoint(srv FruitService, logger *loggers.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fruitID, ok := request.(int64)
		if !ok {
			logger.Error(
				"invalid fruit id",
				loggers.Fields{
					"method":   "GetFruitWithIDEndpoint",
					"received": fmt.Sprintf("%t", request),
				},
			)

			return nil, errInvalidFruitID
		}

		fruitFound, err := srv.GetFruitWithID(ctx, fruitID)
		if err != nil {
			logger.Error(
				"could not get a fruit with the given id",
				loggers.Fields{
					"method": "GetFruitWithIDEndpoint",
					"error":  err,
				},
			)
		}

		logger.Debug(
			"find fruit by id endpoint",
			loggers.Fields{
				"method": "GetFruitWithIDEndpoint",
				"result": fruitFound,
			},
		)

		return newGetFruitWithIDResult(fruitFound, err), nil
	}
}

// MakeCreateFruitEndpoint create endpoint for create fruit service.
func MakeCreateFruitEndpoint(srv FruitService, logger *loggers.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		newFruit, ok := request.(*NewFruit)
		if !ok {
			logger.Error(
				"invalid new fruit type",
				loggers.Fields{
					"method":   "CreateFruitEndpoint",
					"received": fmt.Sprintf("%t", request),
				},
			)

			return nil, errInvalidNewFruitType
		}

		newid, err := srv.Create(ctx, *newFruit)
		if err != nil {
			logger.Error(
				"something went wrong trying to create an fruit with the given id",
				loggers.Fields{
					"method": "CreateFruitEndpoint",
					"error":  err,
				},
			)
		}

		return newCreateFruitResult(newid, err), nil
	}
}

// MakeSearchFruitsEndpoint fruit endpoint to search fruits with filters.
func MakeSearchFruitsEndpoint(srv FruitService, logger *loggers.Logger) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		fruitFilters, ok := request.(SearchFruitFilter)
		if !ok {
			logger.Error(
				"invalid fruit filters",
				loggers.Fields{
					"method":   "SearchFruitsEndpoint",
					"received": fmt.Sprintf("%t", request),
				},
			)

			return nil, errInvalidFruitFilters
		}

		searchResult, err := srv.SearchFruits(ctx, fruitFilters)
		if err != nil {
			logger.Error(
				"something went wrong trying to search fruits with the given filter",
				loggers.Fields{
					"method": "SearchFruitsEndpoint",
					"error":  err,
				},
			)
		}

		return newSearchFruitsDataResult(searchResult, err), nil
	}
}

// MakeGetStatusEndpoint fruit endpoint that shows if the fruit dataset was loaded successfully.
func MakeGetStatusEndpoint(srv FruitService, logger *loggers.Logger) endpoint.Endpoint {
	return func(ctx context.Context, _ interface{}) (interface{}, error) {
		dataSetStatus := srv.DatasetStatus(ctx)

		logger.Debug(
			"get fruit dataset status",
			loggers.Fields{
				"method": "GetStatusEndpoint",
				"result": dataSetStatus,
			},
		)

		return dataSetStatus, nil
	}
}
