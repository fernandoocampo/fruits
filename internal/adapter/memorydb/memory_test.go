package memorydb_test

import (
	"context"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/memorydb"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

func TestCreateFruitInMemoryDB(t *testing.T) {
	t.Parallel()

	newFruitID := int64(1)
	newFruit := repository.Fruit{
		ID:             repository.FruitID(newFruitID),
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
	expectedFruit := repository.Fruit{
		ID:             repository.FruitID(newFruitID),
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
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	newDB := memorydb.NewRepository(logger)
	ctx := context.TODO()

	err := newDB.Save(ctx, newFruitID, newFruit)
	savedFruit, readErr := newDB.FindByID(ctx, newFruitID)

	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.True(t, newFruitID > 0)
	assert.Equal(t, expectedFruit, savedFruit)
}

func TestCreateFruitWithRepository(t *testing.T) {
	t.Parallel()

	fruitID := int64(1234)
	newFruit := repository.NewFruit{
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
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	newDB := memorydb.NewRepository(logger)
	ctx := context.TODO()

	err := newDB.Save(ctx, fruitID, newFruit)
	savedFruit, readErr := newDB.FindByID(ctx, fruitID)

	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, newFruit, savedFruit)
}

func TestCreateFruitInMemoryDBWithLimit(t *testing.T) {
	t.Parallel()

	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	newDB := memorydb.NewRepository(logger)
	ctx := context.TODO()
	newFruit := repository.NewFruit{
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

	for i := 1; i <= 100; i++ {
		err := newDB.Save(ctx, int64(i), newFruit)
		if !assert.NoError(t, err) {
			t.FailNow()
		}
	}
	assert.Equal(t, 100, newDB.Count())
}

func TestFindAllButEmpty(t *testing.T) {
	t.Parallel()

	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	newDB := memorydb.NewRepository(logger)
	ctx := context.TODO()

	result, err := newDB.FindAll(ctx, 10, 10)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	assert.Empty(t, result)
}
