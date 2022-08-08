package monitoring_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/monitoring"
	"github.com/stretchr/testify/assert"
)

func TestStringReport(t *testing.T) {
	t.Parallel()
	expectedReport := "requests=2\nsuccess=1\nerror=1\navailability=50\nnum_fruits=1"
	fruitRepository := fruitRepoMock{
		size: 1,
	}
	metricRepository := metricRepoMock{}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	monitorData := monitoring.MonitorData{
		ReportFrequency:   100 * time.Millisecond,
		FruitRepository:   &fruitRepository,
		MetricsRepository: &metricRepository,
		Logger:            logger,
	}
	agent := monitoring.New(monitorData)
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go agent.Start(ctx)

	agent.CountRequest()
	agent.CountSuccess()
	agent.CountRequest()
	agent.CountError()

	metricRepository.wg.Add(1)
	metricRepository.wg.Wait()

	assert.Equal(t, expectedReport, metricRepository.logs[0])
}

func TestStringReportMultipleEvents(t *testing.T) {
	t.Parallel()
	expectedReport := "requests=4\nsuccess=2\nerror=2\navailability=50\nnum_fruits=2"
	fruitRepository := fruitRepoMock{
		size: 2,
	}
	metricRepository := metricRepoMock{}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	monitorData := monitoring.MonitorData{
		ReportFrequency:   100 * time.Millisecond,
		FruitRepository:   &fruitRepository,
		MetricsRepository: &metricRepository,
		Logger:            logger,
	}
	agent := monitoring.New(monitorData)
	ctx := context.TODO()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go agent.Start(ctx)

	var waitGroup sync.WaitGroup
	waitGroup.Add(3)
	go func() {
		defer waitGroup.Done()
		agent.CountRequest()
		agent.CountSuccess()
	}()
	go func() {
		defer waitGroup.Done()
		agent.CountRequest()
		agent.CountError()
		agent.CountRequest()
		agent.CountSuccess()
	}()
	go func() {
		defer waitGroup.Done()
		agent.CountRequest()
		agent.CountError()
	}()
	waitGroup.Wait()

	metricRepository.wg.Add(1)
	metricRepository.wg.Wait()

	assert.Equal(t, expectedReport, metricRepository.logs[0])
}

type fruitRepoMock struct {
	size int
}

func (r *fruitRepoMock) Count() int {
	return r.size
}

type metricRepoMock struct {
	wg   sync.WaitGroup
	logs []string
}

func (m *metricRepoMock) Push(report string) error {
	defer m.wg.Done()
	m.logs = append(m.logs, report)
	return nil
}
