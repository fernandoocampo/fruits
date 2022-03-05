package application

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitbucket.com/fernandoocampo/fruits/internal/adapter/loggers"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/memorydb"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/metrics"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/monitoring"
	"gitbucket.com/fernandoocampo/fruits/internal/adapter/web"
	"gitbucket.com/fernandoocampo/fruits/internal/configurations"
	"gitbucket.com/fernandoocampo/fruits/internal/fruits"
)

const applicationName = "fruits-service"

// Event contains an application event.
type Event struct {
	Message string
	Error   error
}

// Instance application instance
type Instance struct {
	configuration configurations.Application
	logger        *loggers.Logger
}

// NewInstance creates a new application instance
func NewInstance() *Instance {
	newInstance := Instance{
		logger: loggers.NewLoggerWithStdout(applicationName, loggers.Debug),
	}
	return &newInstance
}

// Run runs fruits application
func (i *Instance) Run() error {
	i.logger.Info("starting application", loggers.Fields{"pkg": "application"})

	confError := i.loadConfiguration()
	if confError != nil {
		return confError
	}
	i.logger.SetLoggerLevel(loggers.Level(i.configuration.LogLevel))
	i.logger.Debug("application configuration", loggers.Fields{"parameters": i.configuration})

	repoFruit := i.createFruitRepository()
	serviceFruit := fruits.NewService(repoFruit, i.logger)
	monitorWorker := i.createMonitoringWorker(repoFruit)
	defer monitorWorker.Shutdown()
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		monitorWorker.Start(ctx)
	}()
	middlewareFruit := fruits.NewFruitMiddleware(serviceFruit, monitorWorker)
	endyear := fruits.NewEndyear(middlewareFruit, i.logger)

	eventStream := make(chan Event)
	i.listenToOSSignal(eventStream)
	i.startWebServer(endyear, eventStream)

	eventMessage := <-eventStream
	i.logger.Info(
		"ending server",
		loggers.Fields{
			"event": eventMessage.Message,
		},
	)

	if eventMessage.Error != nil {
		i.logger.Error(
			"ending server with error",
			loggers.Fields{
				"error": eventMessage.Error,
			},
		)
		return eventMessage.Error
	}
	return nil
}

// Stop stop application, take advantage of this to clean resources
func (i *Instance) Stop() {
	i.logger.Info("stopping the application", loggers.Fields{"pkg": "application"})
}

func (i *Instance) listenToOSSignal(eventStream chan<- Event) {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		osSignal := fmt.Sprintf("%d", <-c)
		event := Event{
			Message: osSignal,
		}
		eventStream <- event
	}()
}

func (i *Instance) createMonitoringWorker(repoFruit monitoring.FruitRepository) *monitoring.Monitor {
	stderrorLogger := loggers.NewBasicLogger(os.Stderr)
	metrics := metrics.New(stderrorLogger)
	monitorData := monitoring.MonitorData{
		ReportFrequency:   time.Duration(i.configuration.MetricsIntervalMillis) * time.Millisecond,
		FruitRepository:   repoFruit,
		MetricsRepository: metrics,
		Logger:            i.logger,
	}
	monitorWorker := monitoring.New(monitorData)

	return monitorWorker
}

// startWebServer starts the web server.
func (i *Instance) startWebServer(endyear fruits.Endyear, eventStream chan<- Event) {
	go func() {
		i.logger.Info("starting http server", loggers.Fields{"http": i.configuration.ApplicationPort})
		handler := web.NewHTTPServer(endyear, i.logger)
		err := http.ListenAndServe(i.configuration.ApplicationPort, handler)
		if err != nil {
			eventStream <- Event{
				Message: "web server was ended with error",
				Error:   err,
			}
			return
		}
		eventStream <- Event{
			Message: "web server was ended",
		}
	}()
}

func (i *Instance) loadConfiguration() error {
	applicationSetUp, err := configurations.Load()
	if err != nil {
		i.logger.Error(
			"application setup could not be loaded",
			loggers.Fields{
				"error": err,
			},
		)
		return errors.New("application setup could not be loaded")
	}
	i.configuration = applicationSetUp
	return nil
}

func (i *Instance) createFruitRepository() *memorydb.FruitMemoryRepository {
	i.logger.Info("initializing database", loggers.Fields{})
	newRepository := memorydb.NewFruitRepository(i.logger)
	if i.configuration.LoadDataset {
		i.logger.Info("loading fruit dataset", loggers.Fields{})
		ctx := context.Background()
		err := newRepository.LoadDatasetWithFile(ctx, i.configuration.FilePath)
		if err != nil {
			i.logger.Error(
				"application failed to load fruit dataset",
				loggers.Fields{
					"error": err,
				},
			)
		}
	}
	return newRepository
}
