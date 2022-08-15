package memorydb

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
)

const (
	fieldsNumber = 14
	bits32       = 32
	bits64       = 64
	base10       = 10
)

// fruit dataset columns index
const (
	idColumn = iota
	countryColumn
	descriptionColumn
	classificationColumn
	yearColumn
	priceColumn
	provinceColumn
	regionColumn
	fincaColumn
	localNameColumn
	wikiPageColumn
	nameColumn
	varietyColumn
	vaultColumn
)

var (
	errIDMustBeAnInteger    = fmt.Errorf("id must be an integer")
	errYearMustBeAnInteger  = fmt.Errorf("year must be an integer")
	errPriceMustBeAnInteger = errors.New("price must be an integer")
)

// ColumnsError defines colum error
type ColumnsError struct {
	FruitValues int
	Record      string
}

// FruitMemoryRepository is the repository handler for fruits in a memory db.
type FruitMemoryRepository struct {
	storage       *Repository
	datasetStatus repository.FruitDatasetStatus
	logger        *loggers.Logger
}

// NewFruitRepository creates a new fruit repository in a dry run repository
func NewFruitRepository(logger *loggers.Logger) *FruitMemoryRepository {
	newRepo := FruitMemoryRepository{
		storage: NewRepository(logger),
		datasetStatus: repository.FruitDatasetStatus{
			Ok: true,
		},
		logger: logger,
	}
	return &newRepo
}

func NewColumnsError(columns int, record string) *ColumnsError {
	return &ColumnsError{
		FruitValues: columns,
		Record:      record,
	}
}

func (c *ColumnsError) Error() string {
	return fmt.Sprintf("invalid fields number, expected: %d, but got %d, record %q", fieldsNumber, c.FruitValues, c.Record)
}

// Save save the given fruit in the postgresql database.
func (u *FruitMemoryRepository) Save(ctx context.Context, newfruit repository.NewFruit) (repository.FruitID, error) {
	u.logger.Debug(
		"save fruit",
		loggers.Fields{
			"method": "repository.FruitMemoryRepository.Save",
			"data":   newfruit,
		},
	)
	newFruitID := u.storage.NewID()
	fruitToPersist := newfruit.ToFruit(repository.FruitID(newFruitID))
	u.logger.Debug(
		"fruit data to persist",
		loggers.Fields{
			"method": "repository.FruitMemoryRepository.Save",
			"data":   fruitToPersist,
		},
	)
	err := u.storage.Save(ctx, newFruitID, fruitToPersist)
	if err != nil {
		u.logger.Error(
			"failed storing fruit",
			loggers.Fields{
				"method": "repository.FruitMemoryRepository.Save",
				"error":  err,
			},
		)
		return 0, errors.New("given fruit could not be stored")
	}
	return repository.FruitID(newFruitID), nil
}

// Update update the given fruit in the postgresql database.
func (u *FruitMemoryRepository) Update(ctx context.Context, fruit repository.Fruit) error {
	u.logger.Debug(
		"updating fruit",
		loggers.Fields{
			"method": "repository.FruitMemoryRepository.Update",
			"data":   fruit,
		},
	)
	err := u.storage.Update(ctx, repository.FruitIDValue(fruit.ID), fruit)
	if err != nil {
		u.logger.Error(
			"failed updating fruit",
			loggers.Fields{
				"method": "repository.FruitMemoryRepository.Update",
				"error":  err,
			},
		)
		return errors.New("given fruit could not be updated")
	}
	return nil
}

// FindByID look for an fruit with the given id
func (u *FruitMemoryRepository) FindByID(ctx context.Context, fruitID repository.FruitID) (*repository.Fruit, error) {
	u.logger.Debug(
		"reading fruit",
		loggers.Fields{
			"method":   "repository.FruitMemoryRepository.FindByID",
			"fruit_id": fruitID,
		},
	)
	result, err := u.storage.FindByID(ctx, repository.FruitIDValue(fruitID))
	if err != nil {
		u.logger.Error(
			"something went wrong trying to read a fruit",
			loggers.Fields{
				"method": "repository.FruitMemoryRepository.FindByID",
				"error":  err,
			},
		)
		return nil, errors.New("something went wrong trying to get the given fruit id")
	}
	var got *repository.Fruit
	if result == nil {
		return got, nil
	}
	fruit, ok := result.(repository.Fruit)
	if !ok {
		u.logger.Error(
			"failed reading fruit",
			loggers.Fields{
				"method": "repository.FruitMemoryRepository.FindByID",
				"error":  "unexpected object",
				"object": result,
			},
		)
		return nil, errors.New("something went wrong trying to get the given fruit id")
	}
	return &fruit, nil
}

// SearchWithFilters memory search
func (u *FruitMemoryRepository) SearchWithFilters(ctx context.Context, filter repository.FruitFilter) (repository.FindFruitsResult, error) {
	result, err := u.storage.FindAll(ctx, filter.Start, filter.Count)
	if err != nil {
		return repository.FindFruitsResult{}, err
	}
	findResult := repository.FindFruitsResult{
		Total: u.storage.Count(),
		Start: filter.Start,
		Count: filter.Count,
	}
	var fruits []repository.Fruit
	for _, v := range result {
		fruit, ok := v.(repository.Fruit)
		if !ok {
			continue
		}
		fruits = append(fruits, fruit)
	}
	findResult.Fruits = fruits
	return findResult, nil
}

// LoadDatasetWithFile load given file data into internal repository
func (u *FruitMemoryRepository) LoadDatasetWithFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		u.logger.Error(
			"loading fruit dataset, but something's wrong with the file",
			loggers.Fields{
				"method":   "memorydb.FruitMemoryRepository.LoadDatasetWithFile",
				"error":    err,
				"filepath": filePath,
			},
		)
		u.datasetStatus.Ok = false
		u.datasetStatus.Message = fmt.Sprintf("could not open file: %s", filePath)
		return errors.New("could not open file")
	}
	defer file.Close()
	err = u.LoadFruitDataset(ctx, bufio.NewScanner(file))
	if err != nil {
		u.datasetStatus.Ok = false
		u.datasetStatus.Message = err.Error()
		return errors.New("could not process file records")
	}
	return nil
}

// LoadFruitDataset load given data into internal repository
func (u *FruitMemoryRepository) LoadFruitDataset(ctx context.Context, scanner *bufio.Scanner) error {
	if scanner == nil {
		u.logger.Info(
			"loading fruit dataset, but the scanner is nil, nothing to do",
			loggers.Fields{
				"method": "memorydb.FruitMemoryRepository.LoadFruitDataset",
			},
		)
		return nil
	}

	recordLine := 1
	scanner.Split(bufio.ScanLines)
	scanner.Scan() // first line has only column names
	for scanner.Scan() {
		fruit, err := u.buildFruit(scanner.Text())
		if err != nil {
			u.logger.Error(
				"loading fruit dataset, but record is not valid",
				loggers.Fields{
					"method": "memorydb.FruitMemoryRepository.LoadFruitDataset",
					"error":  err,
					"line":   recordLine,
				},
			)
			return fmt.Errorf("loading fruit dataset, but record %d is not valid", recordLine)
		}
		err = u.storage.Save(ctx, repository.FruitIDValue(fruit.ID), *fruit)
		if err != nil {
			u.logger.Error(
				"loading fruit dataset, but record cannot be stored into the repository",
				loggers.Fields{
					"method": "memorydb.FruitMemoryRepository.LoadFruitDataset",
					"error":  err,
					"data":   fruit,
				},
			)
			u.datasetStatus.Ok = false
			u.datasetStatus.Message = fmt.Sprintf("loading fruit dataset, but record %d could not be stored", recordLine)
			return errors.New("dataset has invalid records")
		}
		u.storage.UpdateID(repository.FruitIDValue(fruit.ID))
		recordLine++
	}
	return nil
}

// DatasetStatus return the status of the fruit dataset.
func (u *FruitMemoryRepository) DatasetStatus(_ context.Context) (repository.FruitDatasetStatus, error) {
	return u.datasetStatus, nil
}

// Count counts the number of fruits in the memory repo
func (u *FruitMemoryRepository) Count() int {
	return u.storage.Count()
}

// buildFruit build Fruit data from given record data
func (u *FruitMemoryRepository) buildFruit(record string) (*repository.Fruit, error) {
	fruitValues := SplitAtCommas(record)
	if len(fruitValues) != fieldsNumber {
		return nil, NewColumnsError(len(fruitValues), record)
	}
	id, err := strconv.ParseInt(fruitValues[idColumn], base10, bits64)
	if err != nil {
		return nil, errIDMustBeAnInteger
	}
	year, err := strconv.Atoi(fruitValues[yearColumn])
	if err != nil {
		return nil, errYearMustBeAnInteger
	}
	var price float64
	if fruitValues[priceColumn] != "" {
		price, err = strconv.ParseFloat(fruitValues[priceColumn], bits32)
		if err != nil {
			return nil, errPriceMustBeAnInteger
		}
	}
	fruit := repository.Fruit{
		ID:             repository.FruitID(id),
		Name:           fruitValues[nameColumn],
		Variety:        fruitValues[varietyColumn],
		Year:           year,
		Price:          repository.FruitPrice(float32(price)),
		Vault:          fruitValues[vaultColumn],
		Country:        fruitValues[countryColumn],
		Province:       fruitValues[provinceColumn],
		Region:         fruitValues[regionColumn],
		Finca:          fruitValues[fincaColumn],
		Description:    strings.ReplaceAll(fruitValues[descriptionColumn], "\"", ""),
		Classification: fruitValues[classificationColumn],
		LocalName:      fruitValues[localNameColumn],
		WikiPage:       fruitValues[wikiPageColumn],
	}
	return &fruit, nil
}

// SplitAtCommas split s at commas, ignoring commas in strings.
func SplitAtCommas(value string) []string {
	res := []string{}
	var beg int
	var inString bool

	for idx := 0; idx < len(value); idx++ {
		if value[idx] == ',' && !inString {
			res = append(res, value[beg:idx])
			beg = idx + 1
		} else if value[idx] == '"' {
			if !inString {
				inString = true
			} else if idx > 0 && value[idx-1] != '\\' {
				inString = false
			}
		}
	}
	return append(res, value[beg:])
}
