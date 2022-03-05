package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"gitbucket.com/fernandoocampo/fruits/internal/adapter/loggers"
	"gitbucket.com/fernandoocampo/fruits/internal/fruits"
	httptransport "github.com/go-kit/kit/transport/http"
)

func makeEncodeCreateFruitRequest(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.CreateFruitResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.CreateFruitResult",
				loggers.Fields{
					"received": fmt.Sprintf("%+v", response),
					"method":   "encodeCreateFruitRequest",
				},
			)
			return errors.New("cannot build create fruit response")
		}
		w.Header().Set("Content-Type", "application/json")
		message := toCreateFruitResponse(result)
		return json.NewEncoder(w).Encode(message)
	}
}

func makeEncodeGetFruitWithIDResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.GetFruitWithIDResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.GetFruitWithIDResult",
				loggers.Fields{
					"received": fmt.Sprintf("%+v", response),
					"method":   "encodeGetFruitWithIDResponse",
				},
			)
			return errors.New("cannot build get fruit response")
		}
		w.Header().Set("Content-Type", "application/json")
		message := toGetFruitWithIDResponse(result)
		return json.NewEncoder(w).Encode(message)
	}
}

func makeEncodeSearchFruitsResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.SearchFruitsDataResult)
		if !ok {
			logger.Error(
				"cannot transform to fruits.SearchFruitsDataResult",
				loggers.Fields{
					"received": fmt.Sprintf("%T", response),
					"method":   "encodeSearchFruitsResponse",
				},
			)
			return errors.New("cannot build search fruits response")
		}
		w.Header().Set("Content-Type", "application/json")
		message := toSearchFruitsResponse(result)
		return json.NewEncoder(w).Encode(message)
	}
}

func makeEncodeGetStatusResponse(logger *loggers.Logger) httptransport.EncodeResponseFunc {
	return func(ctx context.Context, w http.ResponseWriter, response interface{}) error {
		result, ok := response.(fruits.DatasetStatus)
		if !ok {
			logger.Error(
				"cannot transform to fruits.DatasetStatus",
				loggers.Fields{
					"received": fmt.Sprintf("%T", response),
					"method":   "encodeGetStatusResponse",
				},
			)
			return errors.New("cannot build fruit dataset status response")
		}
		w.Header().Set("Content-Type", "application/json")
		message := toFruitDatasetStatusResponse(result)
		return json.NewEncoder(w).Encode(message)
	}
}
