package document

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/google/uuid"
)

const fruitsTable = "fruits"

var (
	errLoadingAWSConfig = errors.New("unable to load aws config")
	errCreatingDynamodb = errors.New("unable to connect to DynamoDB")
	errSavingFruit      = errors.New("unable to save fruit")
	errGettingFruit     = errors.New("unable to get fruit")
)

// Setup contains dynamodb settings.
type Setup struct {
	Logger   *loggers.Logger
	Region   string
	Endpoint string
}

// DynamoDB defines logic for dynamodb repository.
type DynamoDB struct {
	client *dynamodb.Client
	logger *loggers.Logger
}

func NewDynamoDBClient(ctx context.Context, setup Setup) (*DynamoDB, error) {
	newDynamodb := new(DynamoDB)
	newDynamodb.logger = setup.Logger

	awsconfig, err := newDynamodb.getConfig(ctx, setup.Region, setup.Endpoint)
	if err != nil {
		return nil, errCreatingDynamodb
	}

	newDynamodb.client = dynamodb.NewFromConfig(awsconfig)

	return newDynamodb, nil
}

func (d *DynamoDB) getConfig(ctx context.Context, region, endpoint string) (aws.Config, error) {
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if endpoint != "" {
			return aws.Endpoint{
				// PartitionID:       "aws",
				URL:           endpoint,
				SigningRegion: region,
				// HostnameImmutable: true,
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
		d.logger.Error("unable to load aws config",
			loggers.Fields{
				"error": err,
			},
		)

		return cfg, errLoadingAWSConfig
	}

	return cfg, nil
}

func (d *DynamoDB) FindByID(ctx context.Context, fruitID repository.FruitID) (*repository.Fruit, error) {
	selectedKeys := map[string]string{
		"id": string(fruitID),
	}

	key, err := attributevalue.MarshalMap(selectedKeys)
	if err != nil {
		d.logger.Error("unable to marshal fruit keys", loggers.Fields{"error": err})

		return nil, errGettingFruit
	}

	data, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(fruitsTable),
		Key:       key,
	})
	if err != nil {
		d.logger.Error("unable to get fruit", loggers.Fields{"error": err})

		return nil, errGettingFruit
	}

	var fruit *repository.Fruit

	if data.Item == nil {
		return fruit, nil
	}

	var item Fruit

	err = attributevalue.UnmarshalMap(data.Item, &item)
	if err != nil {
		d.logger.Error("unable to unmarshal fruit", loggers.Fields{"error": err})

		return nil, errGettingFruit
	}

	fruit = item.toRepositoryFruit()

	return fruit, nil
}

func (d *DynamoDB) Save(ctx context.Context, fruit repository.NewFruit) (repository.FruitID, error) {
	newid := uuid.New().String()

	newFruit := transformFruit(newid, fruit)

	data, err := attributevalue.MarshalMap(newFruit)
	if err != nil {
		d.logger.Error("unable to marshal new fruit", loggers.Fields{"error": err})

		return repository.FruitID(""), errSavingFruit
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(fruitsTable),
		Item:      data,
	})
	if err != nil {
		d.logger.Error("unable to store fruit", loggers.Fields{"error": err})

		return repository.FruitID(""), errSavingFruit
	}

	d.logger.Debug(
		"new fruit stored",
		loggers.Fields{
			"id":     newid,
			"output": newFruit,
		},
	)

	return repository.FruitID(newid), nil
}

func (d *DynamoDB) SearchWithFilters(ctx context.Context, filter repository.FruitFilter) (repository.FindFruitsResult, error) {
	return repository.FindFruitsResult{}, nil
}

func (d *DynamoDB) DatasetStatus(ctx context.Context) (repository.FruitDatasetStatus, error) {
	return repository.FruitDatasetStatus{Ok: true}, nil
}

func (d *DynamoDB) Count() int {
	return 1
}
