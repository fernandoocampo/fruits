package topic

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
)

const fruitsTopic = "arn:aws:sns:us-east-1:000000000000:fruits"

var (
	errLoadingAWSConfig = errors.New("unable to load aws config")
	errCreatingDynamodb = errors.New("unable to connect to DynamoDB")
	errPublishingFruit  = errors.New("unable to publish new fruit")
)

// Setup contains dynamodb settings.
type Setup struct {
	Logger   *loggers.Logger
	Region   string
	Endpoint string
}

// SNS defines logic for sns.
type SNS struct {
	client *sns.Client
	logger *loggers.Logger
}

func NewSNSClient(ctx context.Context, setup Setup) (*SNS, error) {
	newsns := new(SNS)
	newsns.logger = setup.Logger

	awsconfig, err := newsns.getConfig(ctx, setup.Region, setup.Endpoint)
	if err != nil {
		return nil, errCreatingDynamodb
	}

	newsns.client = sns.NewFromConfig(awsconfig)

	return newsns, nil
}

func (s *SNS) getConfig(ctx context.Context, region, endpoint string) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				URL:           endpoint,
				SigningRegion: region,
			}, nil
		}

		return aws.Endpoint{}, nil
	})

	cfg, err := config.LoadDefaultConfig(
		ctx, config.WithRegion(region),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("d", "d", "")),
	)
	if err != nil {
		s.logger.Error("unable to load aws config",
			loggers.Fields{
				"error": err,
			},
		)

		return cfg, errLoadingAWSConfig
	}

	return cfg, nil
}

func (s *SNS) Publish(ctx context.Context, fruit repository.NewFruitEvent) error {
	message, err := json.Marshal(fruit)
	if err != nil {
		s.logger.Error("unable to marshal fruit message", loggers.Fields{"error": err})

		return errPublishingFruit
	}

	input := &sns.PublishInput{
		Message:  aws.String(string(message)),
		TopicArn: aws.String(fruitsTopic),
	}

	result, err := s.client.Publish(ctx, input)
	if err != nil {
		s.logger.Error("unable to publish fruit message", loggers.Fields{"error": err})

		return errPublishingFruit
	}

	if result != nil {
		s.logger.Info("publishing new fruit", loggers.Fields{"result": result.MessageId})
	}

	return nil
}
