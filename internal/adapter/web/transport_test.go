package web_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitbucket.com/fernandoocampo/fruits/internal/adapter/loggers"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/web"
	"gitbucket.com/fernandoocampo/fruits/internal/fruits"
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
	Data    int64    `json:"data"`
	Errors  []string `json:"errors"`
}

type webResultGetStatus struct {
	Status    string `json:"status"`
	Message   string `json:"msg"`
	Timestamp int64  `json:"ts"`
}

func TestGetFruitSuccessfully(t *testing.T) {
	fruitID := "1234"
	expectedResponse := webResultGetFruit{
		Success: true,
		Data: &web.Fruit{
			ID:             1234,
			Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
			Variety:        "White Blend",
			Vault:          "Nicosia",
			Year:           87,
			Country:        "Italy",
			Province:       "Sicily & Sardinia",
			Region:         "Etna",
			Description:    "brisk acidity",
			Classification: "Vulk√† Bianco",
			LocalName:      "Kerin OÄôKeefe",
			WikiPage:       "@kerinokeefe",
		},
		Errors: nil,
	}
	fruitToReturn := fruits.Fruit{
		ID:             1234,
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}
	fruitEndyear := fruits.Endyear{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, &fruitToReturn, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndyear, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/fruit/" + fruitID)
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
					ID:   1234,
					Name: "Alicia",
				},
				{
					ID:   1240,
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
				ID:   1234,
				Name: "Alicia",
			},
			{
				ID:   1240,
				Name: "Oliver",
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	fruitEndyear := fruits.Endyear{
		SearchFruitsEndpoint: makeDummySearchFruitsSuccessfullyEndpoint(t, expectedFilter, &serviceResult, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndyear, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/fruit" + queryParams)
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
	fruitID := "1234"
	expectedResponse := webResultGetFruit{
		Success: true,
		Data:    nil,
		Errors:  nil,
	}
	fruitEndyear := fruits.Endyear{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, nil, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndyear, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/fruit/" + fruitID)
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
	fruitID := "1234"
	expectedResponse := webResultGetFruit{
		Success: false,
		Data:    nil,
		Errors:  []string{"any error"},
	}
	errorToReturn := errors.New("any error")
	fruitEndyear := fruits.Endyear{
		GetFruitWithIDEndpoint: makeDummyGetFruitWithIDSuccessfullyEndpoint(t, nil, errorToReturn),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndyear, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/fruit/" + fruitID)
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
	newFruit := web.NewFruit{
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}
	newFruitJson, err := json.Marshal(newFruit)
	if err != nil {
		t.Errorf("unexpected error marshalling new fruit: %s", err)
		t.FailNow()
	}
	fruitEndyear := fruits.Endyear{
		CreateFruitEndpoint: makeDummyCreateFruitSuccessfullyEndpoint(t, 1234, nil),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitHandler := web.NewHTTPServer(fruitEndyear, logger)

	dummyServer := httptest.NewServer(fruitHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateFruit{
		Success: true,
		Data:    1234,
		Errors:  nil,
	}

	createRequest, err := http.NewRequest(http.MethodPut, dummyServer.URL+"/fruit", bytes.NewBuffer(newFruitJson))
	if err != nil {
		t.Errorf("unexpected error creating put request: %s", err)
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
	newFruit := web.NewFruit{
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}
	newFruitJson, err := json.Marshal(newFruit)
	if err != nil {
		t.Errorf("unexpected error marshalling new fruit: %s", err)
		t.FailNow()
	}
	fruitEndyear := fruits.Endyear{
		CreateFruitEndpoint: makeDummyCreateFruitSuccessfullyEndpoint(t, 0, errors.New("any error")),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitHandler := web.NewHTTPServer(fruitEndyear, logger)

	dummyServer := httptest.NewServer(fruitHandler)
	defer dummyServer.Close()

	expectedResponse := webResultCreateFruit{
		Success: false,
		Data:    0,
		Errors:  []string{"any error"},
	}

	createRequest, err := http.NewRequest(http.MethodPut, dummyServer.URL+"/fruit", bytes.NewBuffer(newFruitJson))
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
	expectedResponse := webResultGetStatus{
		Status:    "ok",
		Message:   "",
		Timestamp: 1234,
	}
	status := fruits.DatasetStateOK
	errorMessage := ""
	timeStamp := int64(1234)
	fruitEndyear := fruits.Endyear{
		GetStatusEndpoint: makeDummyGetStatusEndpoint(status, errorMessage, timeStamp),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	httpHandler := web.NewHTTPServer(fruitEndyear, logger)
	dummyServer := httptest.NewServer(httpHandler)
	defer dummyServer.Close()

	response, err := http.Get(dummyServer.URL + "/status")
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
		fruitID, ok := request.(int64)
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

func makeDummyCreateFruitSuccessfullyEndpoint(t *testing.T, newFruitID int64, err error) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		t.Helper()
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
