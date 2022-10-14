package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

const (
	startRecordPosition = 1
	rowPerPage          = 10
)

var (
	errFruitIDNoInt         = errors.New("fruit ID must be a valid integer")
	errDecodingRequest      = errors.New("something went wrong decoding request")
	errNoFruitIDWasProvided = errors.New("fruit ID was not provided")
)

func makeDecodeGetFruitWithIDRequest(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		v := mux.Vars(req)

		fruitID, ok := v["id"]
		if !ok {
			return nil, errNoFruitIDWasProvided
		}

		if fruitID == "" {
			logger.Error(
				"fruit id cannot be empty",
				loggers.Fields{
					"method": "decodeGetFruitWithIDRequest",
				},
			)

			return nil, errFruitIDNoInt
		}

		return fruitID, nil
	}
}

func makeDecodeSearchFruitsRequest(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		filterRequest := SearchFruitFilter{
			Start: startRecordPosition,
			Count: rowPerPage,
		}

		filters := req.URL.Query()

		if v, ok := filters["start"]; ok {
			start, err := strconv.Atoi(v[0])
			if err != nil {
				logger.Error(
					"invalid page parameter, must be an integer",
					loggers.Fields{
						"method": "decodeSearchFruitsRequest",
						"error":  err,
					},
				)

				start = 1
			}

			filterRequest.Start = start
		}

		if v, ok := filters["count"]; ok {
			count, err := strconv.Atoi(v[0])
			if err != nil {
				logger.Error(
					"invalid page size parameter, must be an integer",
					loggers.Fields{
						"method": "decodeSearchFruitsRequest",
						"error":  err,
					},
				)

				count = 10
			}

			filterRequest.Count = count
		}

		filter := filterRequest.toSearchFruitFilter()

		return filter, nil
	}
}

func makeDecodeCreateFruitRequest(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		logger.Debug(
			"decoding create fruit request",
			loggers.Fields{
				"method": "decodeCreateFruitRequest",
			},
		)

		defer req.Body.Close()

		var newFruitRequest NewFruit

		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, fmt.Errorf("something went wrong decoding create fruit request: %w", err)
		}

		err = json.Unmarshal(body, &newFruitRequest)
		if err != nil {
			logger.Error(
				"new fruit request could not be decoded",
				loggers.Fields{
					"method":  "decodeCreateFruitRequest",
					"request": string(body),
					"error":   err,
				},
			)

			return nil, errDecodingRequest
		}

		logger.Debug(
			"fruit request was decoded",
			loggers.Fields{
				"method":  "decodeCreateFruitRequest",
				"request": newFruitRequest,
			},
		)

		domainFruit := newFruitRequest.toFruit()

		return domainFruit, nil
	}
}

func makeEmptyDecoder(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, req *http.Request) (interface{}, error) {
		logger.Debug("calling empty decoder", loggers.Fields{})

		return nil, nil
	}
}
