package fruits_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/stretchr/testify/assert"
)

func TestGetFruitSuccessfully(t *testing.T) {
	t.Parallel()

	fruitID := int64(1234)
	expectedResponse := fruits.GetFruitWithIDResult{
		Fruit: &fruits.Fruit{
			ID: 1234,
		},
		Err: "",
	}
	fruitRepository := fruitRepoMock{
		repo: make(map[int64]repository.Fruit),
	}
	existingFruitID := int64(1234)
	existingFruit := repository.Fruit{
		ID: repository.FruitID(existingFruitID),
	}
	fruitRepository.repo[existingFruitID] = existingFruit
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, logger)
	getFruitEndpoint := fruits.MakeGetFruitWithIDEndpoint(fruitService, logger)
	ctx := context.TODO()

	fruitFound, err := getFruitEndpoint(ctx, fruitID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, fruitFound)
}

func TestGetFruitNotFound(t *testing.T) {
	t.Parallel()

	fruitID := int64(1234)
	expectedResponse := fruits.GetFruitWithIDResult{
		Fruit: nil,
		Err:   "record not found",
	}
	fruitRepository := fruitRepoMock{
		repo: make(map[int64]repository.Fruit),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, logger)
	getFruitEndpoint := fruits.MakeGetFruitWithIDEndpoint(fruitService, logger)
	ctx := context.TODO()

	fruitFound, err := getFruitEndpoint(ctx, fruitID)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, fruitFound)
}

func TestCreateFruitSuccessfully(t *testing.T) {
	t.Parallel()

	fruitRepository := fruitRepoMock{
		repo: make(map[int64]repository.Fruit),
	}
	newFruit := fruits.NewFruit{
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
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, logger)
	createFruitEndpoint := fruits.MakeCreateFruitEndpoint(fruitService, logger)
	ctx := context.TODO()

	result, err := createFruitEndpoint(ctx, &newFruit)

	assert.NoError(t, err)
	assert.NotEmpty(t, result)
}

func TestSearchFruitsEndpointSuccessfully(t *testing.T) {
	t.Parallel()

	givenFilter := fruits.SearchFruitFilter{
		Start: 1,
		Count: 10,
	}
	expectedSearchResult := fruits.SearchFruitsResult{
		Fruits: []fruits.FruitItem{
			{
				ID: 1234,
			},
			{
				ID: 1240,
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	expectedResult := fruits.SearchFruitsDataResult{
		SearchResult: &expectedSearchResult,
	}
	searchResultFixture := repository.FindFruitsResult{
		Fruits: []repository.Fruit{
			{
				ID: 1234,
			},
			{
				ID: 1240,
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	fruitRepository := fruitRepoMock{
		repo:         make(map[int64]repository.Fruit),
		searchResult: searchResultFixture,
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, logger)
	searchFruitEndpoint := fruits.MakeSearchFruitsEndpoint(fruitService, logger)
	ctx := context.TODO()

	fruitsFound, err := searchFruitEndpoint(ctx, givenFilter)

	assert.NoError(t, err)
	assert.Equal(t, expectedResult, fruitsFound)
}
