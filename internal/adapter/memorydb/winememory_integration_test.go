package memorydb_test

import (
	"context"
	"flag"
	"testing"

	"gitbucket.com/fernandoocampo/fruits/internal/adapter/loggers"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/memorydb"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

var (
	integration = flag.Bool("integration", false, "run integration tests")
	filePath    = flag.String("filepath", "../../../data/fruitmag-data.csv", "fruit file dataset")
)

func TestLoadFruitDatasetFromFile(t *testing.T) {
	if !*integration {
		t.Skip("this is an integration test, to execute this test send integration flag to true")
	}
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
