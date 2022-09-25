package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/fruits"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	errBuildingGetFruitResponse    = errors.New("cannot build get fruit response")
	errBuildingFruitDatasetStatus  = errors.New("cannot build fruit dataset status response")
	errBuildingCreateFruitResponse = errors.New("cannot build create fruit response")
	errEncodingResultResponse      = errors.New("cannot encode result")
	errBuildingSearchFruitResponse = errors.New("cannot build search fruits response")
)

func makeEncodeCreateFruitRequest(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, res http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.CreateFruitResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.CreateFruitResult",
				loggers.Fields{
					"received": fmt.Sprintf("%+v", response),
					"method":   "encodeCreateFruitRequest",
				},
			)

			return errBuildingCreateFruitResponse
		}

		res.Header().Set("Content-Type", "application/json")

		message := toCreateFruitResponse(result)

		err := json.NewEncoder(res).Encode(message)
		if err != nil {
			logger.Error(
				"cannot encode Result",
				loggers.Fields{
					"result": fmt.Sprintf("%+v", message),
					"method": "encodeCreateFruitRequest",
				},
			)

			return errEncodingResultResponse
		}

		return nil
	}
}

func makeEncodeGetFruitWithIDResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, res http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.GetFruitWithIDResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.GetFruitWithIDResult",
				loggers.Fields{
					"received": fmt.Sprintf("%+v", response),
					"method":   "encodeGetFruitWithIDResponse",
				},
			)

			return errBuildingGetFruitResponse
		}

		res.Header().Set("Content-Type", "application/json")

		message := toGetFruitWithIDResponse(result)

		err := json.NewEncoder(res).Encode(message)
		if err != nil {
			logger.Error(
				"cannot encode Result",
				loggers.Fields{
					"result": fmt.Sprintf("%+v", message),
					"method": "encodeGetFruitWithIDResponse",
				},
			)

			return errEncodingResultResponse
		}

		return nil
	}
}

func makeEncodeSearchFruitsResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, res http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.SearchFruitsDataResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.SearchFruitsDataResult",
				loggers.Fields{
					"received": fmt.Sprintf("%T", response),
					"method":   "encodeSearchFruitsResponse",
				},
			)

			return errBuildingSearchFruitResponse
		}

		res.Header().Set("Content-Type", "application/json")

		message := toSearchFruitsResponse(result)

		err := json.NewEncoder(res).Encode(message)
		if err != nil {
			logger.Error(
				"cannot encode Result",
				loggers.Fields{
					"result": fmt.Sprintf("%+v", message),
					"method": "encodeSearchFruitsResponse",
				},
			)

			return errEncodingResultResponse
		}

		return nil
	}
}

func makeEncodeGetStatusResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, res http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.DatasetStatus)
		if !ok {
			logger.Error(
				"cannot transform to fruits.DatasetStatus",
				loggers.Fields{
					"received": fmt.Sprintf("%T", response),
					"method":   "encodeGetStatusResponse",
				},
			)

			return errBuildingFruitDatasetStatus
		}

		res.Header().Set("Content-Type", "application/json")

		message := toFruitDatasetStatusResponse(result)

		err := json.NewEncoder(res).Encode(message)
		if err != nil {
			logger.Error(
				"cannot encode Result",
				loggers.Fields{
					"result": fmt.Sprintf("%+v", message),
					"method": "encodeGetStatusResponse",
				},
			)

			return errEncodingResultResponse
		}

		return nil
	}
}
