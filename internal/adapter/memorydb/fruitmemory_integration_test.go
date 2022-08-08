package memorydb_test

import (
	"context"
	"flag"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/memorydb"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

func TestLoadFruitDatasetFromFile(t *testing.T) {
	integration := flag.Bool("integration", false, "run integration tests")
	filePath := flag.String("filepath", "../../../data/fruitmag-data.csv", "fruit file dataset")
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}
	t.Parallel()
	ctx := context.TODO()
	expectedDatasetStatus := repository.FruitDatasetStatus{
		Ok: true,
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitRepo := memorydb.NewFruitRepository(logger)

	err := fruitRepo.LoadDatasetWithFile(ctx, *filePath)

	assert.NoError(t, err)
	status, err := fruitRepo.DatasetStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedDatasetStatus, status)

	filter := repository.FruitFilter{
		Start: 1,
		Count: 10,
	}
	result, err := fruitRepo.SearchWithFilters(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, 129971, result.Total)
}
