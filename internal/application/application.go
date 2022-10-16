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

	"github.com/fernandoocampo/fruits/internal/adapter/document"
	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/metrics"
	"github.com/fernandoocampo/fruits/internal/adapter/monitoring"
	"github.com/fernandoocampo/fruits/internal/adapter/topic"
	"github.com/fernandoocampo/fruits/internal/adapter/web"
	"github.com/fernandoocampo/fruits/internal/configurations"
	"github.com/fernandoocampo/fruits/internal/fruits"
)

const applicationName = "fruits-service"

// Event contains an application event.
type Event struct {
	Message string
	Error   error
}

// Instance application instance.
type Instance struct {
	configuration configurations.Application
	logger        *loggers.Logger
}

var (
	errCreatingTopic      = errors.New("unable to create topic client")
	errCreatingRepository = errors.New("unable to create repository client")
	errLoadingApplication = errors.New("application setup could not be loaded")
)

// NewInstance creates a new application instance.
func NewInstance() *Instance {
	newInstance := Instance{
		logger: loggers.NewLoggerWithStdout(applicationName, loggers.Debug),
	}

	return &newInstance
}

// Run runs fruits application.
func (i *Instance) Run() error {
	i.logger.Info("starting application", loggers.Fields{"pkg": "application"})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	confError := i.loadConfiguration()
	if confError != nil {
		return confError
	}

	i.logger.SetLoggerLevel(loggers.Level(i.configuration.LogLevel))
	i.logger.Debug("application configuration", loggers.Fields{"parameters": i.configuration})

	i.logger.Info("metadata",
		loggers.Fields{
			"version": i.configuration.Version,
			"commit":  i.configuration.CommitHash,
		},
	)

	repoFruit, err := i.createFruitRepository(ctx)
	if err != nil {
		return errLoadingApplication
	}

	repoTopic, err := i.createFruitTopic(ctx)
	if err != nil {
		return errLoadingApplication
	}

	serviceFruit := fruits.NewService(repoFruit, repoTopic, i.logger)

	monitorWorker := i.createMonitoringWorker(ctx, repoFruit)
	defer monitorWorker.Shutdown()

	middlewareFruit := fruits.NewFruitMiddleware(serviceFruit, monitorWorker)
	endpoints := fruits.NewEndpoints(middlewareFruit, i.logger)

	eventStream := make(chan Event)
	i.listenToOSSignal(eventStream)
	i.startWebServer(endpoints, eventStream)

	eventMessage := <-eventStream

	i.logger.Info("ending server",
		loggers.Fields{
			"event": eventMessage.Message,
		})

	if eventMessage.Error != nil {
		i.logger.Error("ending server with error",
			loggers.Fields{
				"error": eventMessage.Error,
			})

		return eventMessage.Error
	}

	return nil
}

// Stop stop application, take advantage of this to clean resources.
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
			Error:   nil,
		}
		eventStream <- event
	}()
}

func (i *Instance) createMonitoringWorker(ctx context.Context, repoFruit monitoring.FruitRepository) *monitoring.Monitor {
	stderrorLogger := loggers.NewBasicLogger(os.Stderr)
	metrics := metrics.New(stderrorLogger)
	monitorData := monitoring.MonitorData{
		ReportFrequency:   time.Duration(i.configuration.MetricsIntervalMillis) * time.Millisecond,
		FruitRepository:   repoFruit,
		MetricsRepository: metrics,
		Logger:            i.logger,
	}
	monitorWorker := monitoring.New(monitorData)
	monitorWorker.Start(ctx)

	return monitorWorker
}

// startWebServer starts the web server.
func (i *Instance) startWebServer(endpoints fruits.Endpoints, eventStream chan<- Event) {
	go func() {
		i.logger.Info("starting http server", loggers.Fields{"http": i.configuration.ApplicationPort})
		handler := web.NewHTTPServer(endpoints, i.logger)

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
			Error:   nil,
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

		return errLoadingApplication
	}

	i.configuration = applicationSetUp

	return nil
}

func (i *Instance) createFruitRepository(ctx context.Context) (*document.DynamoDB, error) {
	i.logger.Info("initializing database", loggers.Fields{})

	dbSetup := document.Setup{
		Logger:   i.logger,
		Region:   i.configuration.CloudRegion,
		Endpoint: i.configuration.CloudEndpointURL,
	}

	newRepository, err := document.NewDynamoDBClient(ctx, dbSetup)
	if err != nil {
		i.logger.Error("unable to create dynamodb client", loggers.Fields{"error": err})

		return nil, errCreatingRepository
	}

	return newRepository, nil
}

func (i *Instance) createFruitTopic(ctx context.Context) (*topic.SNS, error) {
	i.logger.Info("initializing topic client", loggers.Fields{})

	dbSetup := topic.Setup{
		Logger:   i.logger,
		Region:   i.configuration.CloudRegion,
		Endpoint: i.configuration.CloudEndpointURL,
	}

	newTopic, err := topic.NewSNSClient(ctx, dbSetup)
	if err != nil {
		i.logger.Error("unable to create sns client", loggers.Fields{"error": err})

		return nil, errCreatingTopic
	}

	return newTopic, nil
}
