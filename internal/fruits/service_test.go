package fruits_test

import (
	"context"
	"errors"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var errAnyError = errors.New("something went wrong accessing db")

func TestFindFruitSuccessfully(t *testing.T) {
	t.Parallel()

	fruitID := "1234"
	expectedFruit := fruits.Fruit{
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
		LocalName:      "Kerin Test",
		WikiPage:       "@kerinokeefe",
	}
	fruitRepository := new(fruitRepoMock)
	fruitRepository.repo = make(map[string]repository.Fruit)
	existingFruitID := "1234"
	existingFruit := repository.Fruit{
		ID:             repository.FruitID(existingFruitID),
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin Test",
		WikiPage:       "@kerinokeefe",
	}
	fruitRepository.repo[existingFruitID] = existingFruit
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	fruitFound, err := fruitService.GetFruitWithID(ctx, fruitID)

	assert.NoError(t, err)
	assert.Equal(t, &expectedFruit, fruitFound)
}

func TestFindFruitNotFound(t *testing.T) {
	t.Parallel()

	fruitID := "1234"
	fruitRepository := fruitRepoMock{
		repo: make(map[string]repository.Fruit),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	fruitFound, err := fruitService.GetFruitWithID(ctx, fruitID)

	assert.NoError(t, err)
	assert.Nil(t, fruitFound)
}

func TestFindFruitWithError(t *testing.T) {
	t.Parallel()

	fruitID := "1234"
	fruitRepository := fruitRepoMock{
		searchResult:  *new(repository.FindFruitsResult),
		repo:          make(map[string]repository.Fruit),
		err:           errAnyError,
		dataSetStatus: *new(repository.FruitDatasetStatus),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	fruitFound, err := fruitService.GetFruitWithID(ctx, fruitID)

	assert.Error(t, err)
	assert.Nil(t, fruitFound)
	assert.Equal(t, errAnyError, err)
}

func TestSearchFruitsSuccessfully(t *testing.T) {
	t.Parallel()

	givenFilter := fruits.SearchFruitFilter{
		Start: 1,
		Count: 10,
	}
	expectedResult := fruits.SearchFruitsResult{
		Fruits: []fruits.FruitItem{
			{
				ID:   "1234",
				Name: "Quinta dos Avidagos 2011 Avidagos Red (Douro)",
			},
			{
				ID:   "1240",
				Name: "Stemmari 2013 Dalila White (Terre Siciliane)",
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	searchResultFixture := repository.FindFruitsResult{
		Fruits: []repository.Fruit{
			{
				ID:   "1234",
				Name: "Quinta dos Avidagos 2011 Avidagos Red (Douro)",
			},
			{
				ID:   "1240",
				Name: "Stemmari 2013 Dalila White (Terre Siciliane)",
			},
		},
		Total: 2,
		Start: 1,
		Count: 10,
	}
	fruitRepository := fruitRepoMock{
		err:           nil,
		repo:          make(map[string]repository.Fruit),
		searchResult:  searchResultFixture,
		dataSetStatus: *new(repository.FruitDatasetStatus),
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	fruitsFound, err := fruitService.SearchFruits(ctx, givenFilter)

	assert.NoError(t, err)
	assert.Equal(t, &expectedResult, fruitsFound)
}

func TestDatasetOk(t *testing.T) {
	t.Parallel()

	expectedStatus := fruits.DatasetStateOK
	fruitRepository := fruitRepoMock{
		err:          nil,
		repo:         nil,
		searchResult: *new(repository.FindFruitsResult),
		dataSetStatus: repository.FruitDatasetStatus{
			Ok: true,
		},
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	got := fruitService.DatasetStatus(ctx)

	assert.Equal(t, expectedStatus, got.Status)
	assert.Empty(t, got.Message)
	assert.Greater(t, got.Timestamp, int64(0))
}

func TestDatasetWithError(t *testing.T) {
	t.Parallel()

	expectedStatus := fruits.DatasetStateError
	expectedMessage := "dataset source was not found"
	fruitRepository := fruitRepoMock{
		err:          nil,
		repo:         nil,
		searchResult: *new(repository.FindFruitsResult),
		dataSetStatus: repository.FruitDatasetStatus{
			Ok:      false,
			Message: "dataset source was not found",
		},
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitService := fruits.NewService(&fruitRepository, &publisherMock{}, logger)
	ctx := context.TODO()

	got := fruitService.DatasetStatus(ctx)

	assert.Equal(t, expectedStatus, got.Status)
	assert.Equal(t, expectedMessage, got.Message)
	assert.Greater(t, got.Timestamp, int64(0))
}

type publisherMock struct{}

func (p *publisherMock) Publish(_ context.Context, event repository.NewFruitEvent) error {
	return nil
}

type fruitRepoMock struct {
	err           error
	repo          map[string]repository.Fruit
	searchResult  repository.FindFruitsResult
	dataSetStatus repository.FruitDatasetStatus
}

func (u *fruitRepoMock) FindByID(_ context.Context, fruitID repository.FruitID) (*repository.Fruit, error) {
	if u.err != nil {
		return nil, u.err
	}

	var result *repository.Fruit

	fruitFound, ok := u.repo[repository.FruitIDValue(fruitID)]
	if !ok {
		return result, nil
	}

	return &fruitFound, nil
}

func (u *fruitRepoMock) Save(ctx context.Context, fruit repository.NewFruit) (repository.FruitID, error) {
	if u.err != nil {
		return "", u.err
	}
	id := uuid.New().String()
	newFruit := transformNewFruitToFruit(repository.FruitID(id), fruit)
	u.repo[id] = newFruit
	return repository.FruitID(id), nil
}

func (u *fruitRepoMock) Update(ctx context.Context, fruit repository.Fruit) error {
	if u.err != nil {
		return u.err
	}
	u.repo[repository.FruitIDValue(fruit.ID)] = fruit
	return nil
}

func (u *fruitRepoMock) SearchWithFilters(ctx context.Context, filter repository.FruitFilter) (repository.FindFruitsResult, error) {
	var result repository.FindFruitsResult
	if u.err != nil {
		return result, u.err
	}
	return u.searchResult, nil
}

func (u *fruitRepoMock) DatasetStatus(ctx context.Context) (repository.FruitDatasetStatus, error) {
	var result repository.FruitDatasetStatus
	if u.err != nil {
		return result, u.err
	}
	return u.dataSetStatus, nil
}

func transformNewFruitToFruit(fruitID repository.FruitID, newFruit repository.NewFruit) repository.Fruit {
	return repository.Fruit{
		ID:             fruitID,
		Name:           newFruit.Name,
		Variety:        newFruit.Variety,
		Year:           newFruit.Year,
		Price:          newFruit.Price,
		Vault:          newFruit.Vault,
		Country:        newFruit.Country,
		Province:       newFruit.Province,
		Region:         newFruit.Region,
		Finca:          newFruit.Finca,
		Description:    newFruit.Description,
		Classification: newFruit.Classification,
		LocalName:      newFruit.LocalName,
		WikiPage:       newFruit.WikiPage,
	}
}
