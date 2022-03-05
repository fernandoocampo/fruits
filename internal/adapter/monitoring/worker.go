package monitoring

import (
	"context"
	"fmt"
	"time"

	"gitbucket.com/fernandoocampo/fruits/internal/adapter/loggers"
)

const (
	reportFormat = "requests=%d\nsuccess=%d\nerror=%d\navailability=%d\nnum_fruits=%d"
)
const (
	requests          = "requests"
	successfulRequest = "success"
	failedRequest     = "error"
	availability      = "availability"
	fruits            = "num_fruits"
)

// FruitRepository defines behavior to count the fruits.
type FruitRepository interface {
	Count() int
}

// MetricsRepository defines behavior to push report to the metrics server.
type MetricsRepository interface {
	Push(report string) error
}

// MonitorData monitor data to initialize the monitor worker
type MonitorData struct {
	ReportFrequency   time.Duration
	FruitRepository   FruitRepository
	MetricsRepository MetricsRepository
	Logger            *loggers.Logger
}

type Monitor struct {
	report            map[string]int
	eventStream       chan string
	ticker            *time.Ticker
	fruitRepository   FruitRepository
	metricsRepository MetricsRepository
	logger            *loggers.Logger
}

func New(params MonitorData) *Monitor {
	newMonitor := Monitor{
		report:            newReport(),
		eventStream:       make(chan string),
		ticker:            time.NewTicker(params.ReportFrequency),
		fruitRepository:   params.FruitRepository,
		metricsRepository: params.MetricsRepository,
		logger:            params.Logger,
	}
	return &newMonitor
}

// Start start to push reports to the metrics repository.
func (m *Monitor) Start(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() != nil {
				m.logger.Info("receiving signal to finish context", loggers.Fields{"reason": ctx.Err().Error()})
			}
			return
		case metricName, ok := <-m.eventStream:
			if !ok {
				return
			}
			m.report[metricName]++
		case <-m.ticker.C:
			m.Flush()
		}
	}
}

// CountRequest count a request
func (m *Monitor) CountRequest() {
	m.count(requests)
}

// CountSuccess count a successful request
func (m *Monitor) CountSuccess() {
	m.count(successfulRequest)
}

// CountError count an unsuccessful request
func (m *Monitor) CountError() {
	m.count(failedRequest)
}

func (m *Monitor) Flush() {
	m.metricsRepository.Push(m.generateReport())
	m.report = newReport()
}

func (m *Monitor) count(metricName string) {
	go func() {
		m.eventStream <- metricName
	}()
}

func (m *Monitor) generateReport() string {
	if m.report[requests] > 0 {
		m.report[availability] = int(100 * (float32(m.report[successfulRequest]) / float32(m.report[requests])))
	}
	m.report[fruits] = m.fruitRepository.Count()
	return fmt.Sprintf(
		reportFormat,
		m.report[requests],
		m.report[successfulRequest],
		m.report[failedRequest],
		m.report[availability],
		m.report[fruits],
	)
}

// Shutdown turn off the monitor
func (m *Monitor) Shutdown() {
	m.ticker.Stop()
	close(m.eventStream)
}

func newReport() map[string]int {
	return map[string]int{
		requests:          0,
		successfulRequest: 0,
		failedRequest:     0,
		availability:      0,
		fruits:            0,
	}
}
