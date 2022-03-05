package web

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func makeDecodeGetFruitWithIDRequest(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		v := mux.Vars(r)
		fruitIDParam, ok := v["id"]
		if !ok {
			return nil, errors.New("fruit ID was not provided")
		}
		fruitID, err := strconv.ParseInt(fruitIDParam, 10, 64)
		if err != nil {
			logger.Error(
				"invalid fruit id",
				loggers.Fields{
					"method":   "decodeGetFruitWithIDRequest",
					"received": fruitIDParam,
					"error":    err,
				},
			)
			return nil, errors.New("fruit ID must be a valid integer")
		}
		return fruitID, nil
	}
}

func makeDecodeSearchFruitsRequest(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		filterRequest := SearchFruitFilter{
			Start: 1,
			Count: 10,
		}

		filters := r.URL.Query()

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
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		logger.Debug(
			"decoding create fruit request",
			loggers.Fields{
				"method": "decodeCreateFruitRequest",
			},
		)
		var req NewFruit
		defer r.Body.Close()

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(body, &req)
		if err != nil {
			logger.Error(
				"new fruit request could not be decoded",
				loggers.Fields{
					"method":  "decodeCreateFruitRequest",
					"request": string(body),
					"error":   err,
				},
			)
			return nil, err
		}

		logger.Debug(
			"fruit request was decoded",
			loggers.Fields{
				"method":  "decodeCreateFruitRequest",
				"request": req,
			},
		)

		domainFruit := req.toFruit()

		return domainFruit, nil
	}
}

func makeEmptyDecoder(logger *loggers.Logger) httptransport.DecodeRequestFunc {
	return func(ctx context.Context, r *http.Request) (interface{}, error) {
		logger.Debug("calling empty decoder", loggers.Fields{})
		return nil, nil
	}
}
