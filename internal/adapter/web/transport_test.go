package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/web"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/go-kit/kit/endpoint"
	"github.com/stretchr/testify/assert"
)

type webResultGetFruit struct {
	Success bool       `json:"success"`
	Data    *web.Fruit `json:"data"`
	Errors  []string   `json:"errors"`
}

type webResultSearchFruits struct {
	Success bool                    `json:"success"`
	Data    *web.SearchFruitsResult `json:"data"`
	Errors  []string                `json:"errors"`
}

type webResultCreateFruit struct {
	Success bool     `json:"success"`
	Data    string   `json:"data"`
	Errors  []string `json:"errors"`
}

type webResultGetStatus struct {
	Status    string `json:"status"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

var errAnyError = errors.New("any error")

func TestGetFruitSuccessfully(t *testing.T) {
	t.Parallel()

	fruitID := "1234"
	expectedResponse := webResultGetFruit{
		Success: true,
		Data: &web.Fruit{
			ID:             "1234",
			Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
			Variety:        "White Blend",
			Vault:          "Nicosia",
			Year:           87,
			Country:        "Italy",
			Province:       "Sicily & Sardinia",
			Region:         "Etna",
			Description:    "brisk acidity",
			Classification: "Vulka Bianco",
			LocalName:      "Kerin OoKeefe",
			WikiPage:       "@kerinokeefe",
		},
		Errors: nil,
	}
	fruitToReturn := fruits.Fruit{
		ID:             "1234",
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OoKeefe",
		WikiPage:       "@kerinokeefe",
	}
	fruitEndpoints := fruits.Endpoints{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, &fruitToReturn, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndpoints, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	ctx := context.TODO()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL+"/fruit/"+fruitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetFruit

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestSearchFruitsSuccessfully(t *testing.T) {
	t.Parallel()

	queryParams := "?page=1&pagesize=10"
	expectedFilter := fruits.SearchFruitFilter{
		Start: 1,
		Count: 10,
	}
	expectedResponse := webResultSearchFruits{
		Success: true,
		Data: &web.SearchFruitsResult{
			Fruits: []web.FruitItemResult{
				{
					ID:   "1234",
					Name: "Alicia",
				},
				{
					ID:   "1240",
					Name: "Oliver",
				},
			},
			Total: 2,
			Start: 1,
			Count: 10,
		},
		Errors: nil,
	}

	serviceResult := fruits.SearchFruitsResult{
		Fruits: []fruits.FruitItem{
			{
				ID:   "1234",
				Name: "Alicia",
			},
			{
				ID:   "1240",
				Name: "Oliver",
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	fruitEndpoints := fruits.Endpoints{
		SearchFruitsEndpoint: makeDummySearchFruitsSuccessfullyEndpoint(t, expectedFilter, &serviceResult, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndpoints, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	ctx := context.TODO()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL+"/fruit"+queryParams, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultSearchFruits

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestGetFruitNotFound(t *testing.T) {
	t.Parallel()

	fruitID := "ad1a4350-978b-4a08-83f7-20199dc8f21a"
	expectedResponse := webResultGetFruit{
		Success: true,
		Data:    nil,
		Errors:  nil,
	}
	fruitEndpoints := fruits.Endpoints{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, nil, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndpoints, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	ctx := context.TODO()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL+"/fruit/"+fruitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetFruit

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestGetFruitWithError(t *testing.T) {
	t.Parallel()

	fruitID := "1234"
	expectedResponse := webResultGetFruit{
		Success: false,
		Data:    nil,
		Errors:  []string{"any error"},
	}
	errorToReturn := errAnyError
	fruitEndpoints := fruits.Endpoints{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, nil, errorToReturn),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndpoints, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	ctx := context.TODO()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL+"/fruit/"+fruitID, nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetFruit

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestPostFruitSuccessfully(t *testing.T) {
	t.Parallel()

	newFruit := web.NewFruit{
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OoKeefe",
		WikiPage:       "@kerinokeefe",
	}
	newFruitJson, err := json.Marshal(newFruit)
	if err != nil {
		t.Errorf("unexpected error marshalling new fruit: %s", err)
		t.FailNow()
	}
	fruitEndpoints := fruits.Endpoints{
		CreateFruitEndpoint: makeDummyCreateFruitSuccessfullyEndpoint(t, "1234", nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitHandler := web.NewHTTPServer(fruitEndpoints, logger)

	dummyServer := httptest.NewServer(fruitHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateFruit{
		Success: true,
		Data:    "1234",
		Errors:  nil,
	}

	ctx := context.TODO()
	createRequest, err := http.NewRequestWithContext(ctx, http.MethodPut, dummyServer.URL+"/fruit", bytes.NewBuffer(newFruitJson))
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(createRequest)
	if err != nil {
		t.Errorf("unexected error creating new fruit request: %s", err)
	}
	defer response.Body.Close()

	var result webResultCreateFruit

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestPostFruitWithError(t *testing.T) {
	t.Parallel()

	newFruit := web.NewFruit{
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OoKeefe",
		WikiPage:       "@kerinokeefe",
	}
	newFruitJson, err := json.Marshal(newFruit)
	if err != nil {
		t.Errorf("unexpected error marshalling new fruit: %s", err)
		t.FailNow()
	}
	fruitEndpoints := fruits.Endpoints{
		CreateFruitEndpoint: makeDummyCreateFruitSuccessfullyEndpoint(t, "", errAnyError),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitHandler := web.NewHTTPServer(fruitEndpoints, logger)

	dummyServer := httptest.NewServer(fruitHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateFruit{
		Success: false,
		Data:    "",
		Errors:  []string{"any error"},
	}

	ctx := context.TODO()
	createRequest, err := http.NewRequestWithContext(ctx, http.MethodPut, dummyServer.URL+"/fruit", bytes.NewBuffer(newFruitJson))
	if err != nil {
		t.Errorf("unexected error creating put request: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(createRequest)
	if err != nil {
		t.Errorf("unexected error creating new fruit request: %s", err)
	}
	defer response.Body.Close()

	var result webResultCreateFruit

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func TestStatusSuccessfully(t *testing.T) {
	t.Parallel()

	expectedResponse := webResultGetStatus{
		Status:    "ok",
		Message:   "",
		Timestamp: 1234,
	}
	status := fruits.DatasetStateOK
	errorMessage := ""
	timeStamp := int64(1234)
	fruitEndpoints := fruits.Endpoints{
		GetStatusEndpoint: makeDummyGetStatusEndpoint(status, errorMessage, timeStamp),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndpoints, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	ctx := context.TODO()
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL+"/status", nil)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}
	defer response.Body.Close()

	var result webResultGetStatus

	err = json.NewDecoder(response.Body).Decode(&result)
	if err != nil {
		t.Error("unexpected error", err)
		t.FailNow()
	}

	assert.Equal(t, expectedResponse, result)
}

func makeDummyGetFruitWithIDSuccessfullyEndpoint(t *testing.T, fruitToReturn *fruits.Fruit, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		fruitID, ok := request.(string)
		if !ok {
			t.Errorf("fruit id parameter is not valid: %+v", fruitID)
			t.FailNow()
		}

		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := fruits.GetFruitWithIDResult{
			Fruit: fruitToReturn,
			Err:   errMessage,
		}
		return result, nil
	}
}

func makeDummySearchFruitsSuccessfullyEndpoint(t *testing.T, expectedFilter fruits.SearchFruitFilter, resultToReturn *fruits.SearchFruitsResult, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
		filter, ok := request.(fruits.SearchFruitFilter)
		if !ok {
			t.Errorf("fruit filter parameter is not valid: %T", request)
			t.FailNow()
		}

		assert.Equal(t, expectedFilter, filter)

		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := fruits.SearchFruitsDataResult{
			SearchResult: resultToReturn,
			Err:          errMessage,
		}
		return result, nil
	}
}

func makeDummyCreateFruitSuccessfullyEndpoint(t *testing.T, newFruitID string, err error) endpoint.Endpoint {
	t.Helper()
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		_, ok := request.(*fruits.NewFruit)
		if !ok {
			t.Errorf("fruit parameter is not valid: %T", request)
			t.FailNow()
		}
		t.Log("using newFruitID", newFruitID)
		var errMessage string
		if err != nil {
			errMessage = err.Error()
		}
		result := fruits.CreateFruitResult{
			ID:  newFruitID,
			Err: errMessage,
		}
		return result, nil
	}
}

func makeDummyGetStatusEndpoint(statusToReturn fruits.DatasetState, errMessage string, ts int64) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		result := fruits.DatasetStatus{
			Status:    statusToReturn,
			Message:   errMessage,
			Timestamp: ts,
		}
		return result, nil
	}
}
