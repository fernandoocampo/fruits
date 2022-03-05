package fruits

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/fernandoocampo/fruits/internal/adapter/repository"
)

// FruitService defines behavior for fruit service business logic.
type FruitService interface {
	GetFruitWithID(ctx context.Context, fruitID int64) (*Fruit, error)
	Create(ctx context.Context, newfruit NewFruit) (int64, error)
	SearchFruits(ctx context.Context, givenFilter SearchFruitFilter) (*SearchFruitsResult, error)
	DatasetStatus(ctx context.Context) DatasetStatus
}

// DatasetState define fruit dataset state.
type DatasetState string

// CreateFruitResult standard response for create Fruit
type CreateFruitResult struct {
	ID  int64
	Err string
}

// GetFruitWithIDResult standard roespnse for get a Fruit with an ID.
type GetFruitWithIDResult struct {
	Fruit *Fruit
	Err   string
}

// SearchFruitsDataResult standard roespnse for get a Fruit with an ID.
type SearchFruitsDataResult struct {
	SearchResult *SearchFruitsResult
	Err          string
}

// SearchFruitFilter contains filters to search fruits
type SearchFruitFilter struct {
	// Page page to query
	Start int
	// rows per page
	Count int
}

// SearchFruitsResult contains search fruits result data.
type SearchFruitsResult struct {
	Fruits []FruitItem
	Total  int
	Start  int
	Count  int
}

// NewFruit contains fruit data.
type NewFruit struct {
	Name           string  `json:"name"`
	Variety        string  `json:"variety"`
	Vault          string  `json:"vault"`
	Year           int     `json:"year"`
	Price          float32 `json:"price,omitempty"`
	Country        string  `json:"country"`
	Province       string  `json:"province"`
	Region         string  `json:"region,omitempty"`
	Finca          string  `json:"finca,omitempty"`
	Description    string  `json:"description"`
	Classification string  `json:"classification"`
	LocalName      string  `json:"local_name"`
	WikiPage       string  `json:"wiki_page"`
}

// Fruit contains fruit data.
type Fruit struct {
	ID             int64   `json:"id"`
	Name           string  `json:"name"`
	Variety        string  `json:"variety"`
	Vault          string  `json:"vault"`
	Year           int     `json:"year"`
	Price          float32 `json:"price,omitempty"`
	Country        string  `json:"country"`
	Province       string  `json:"province"`
	Region         string  `json:"region,omitempty"`
	Finca          string  `json:"finca,omitempty"`
	Description    string  `json:"description"`
	Classification string  `json:"classification"`
	LocalName      string  `json:"local_name"`
	WikiPage       string  `json:"wiki_page"`
}

// FruitItem contains few fruit data, just to show reference data.
type FruitItem struct {
	ID int64 `json:"id"`
	// Name fruit's name.
	Name string `json:"name"`
}

// DatasetStatus contains data about the fruit dataset result
type DatasetStatus struct {
	Status    DatasetState
	Message   string
	Timestamp int64
}

const (
	DatasetStateOK    DatasetState = "ok"
	DatasetStateError DatasetState = "error"
)

// ToFruitPortOut transforms new fruit to a fruit port out.
func (u NewFruit) ToFruitPortOut() repository.NewFruit {
	return repository.NewFruit{
		Name:           u.Name,
		Variety:        u.Variety,
		Year:           u.Year,
		Price:          repository.FruitPrice(u.Price),
		Vault:          u.Vault,
		Country:        u.Country,
		Province:       u.Province,
		Region:         u.Region,
		Finca:          u.Finca,
		Description:    u.Description,
		Classification: u.Classification,
		LocalName:      u.LocalName,
		WikiPage:       u.WikiPage,
	}
}

// Validate check is the given data to create a fruit is correct.
func (n NewFruit) Validate() error {
	var errorMessage string
	if n.Name == "" {
		errorMessage = "name, "
	}
	if n.Classification == "" {
		errorMessage += "classification, "
	}
	if n.Country == "" {
		errorMessage += "country, "
	}
	if n.Vault == "" {
		errorMessage += "vault, "
	}
	if errorMessage == "" {
		return nil
	}
	result := fmt.Sprintf("these fields are mandatory: %s.", errorMessage[:len(errorMessage)-2])
	return errors.New(result)
}

// NewFruit transforms new fruit to a fruit port out.
func (n NewFruit) NewFruit(fruitID int64) Fruit {
	return Fruit{
		ID:             fruitID,
		Name:           n.Name,
		Variety:        n.Variety,
		Year:           n.Year,
		Price:          n.Price,
		Vault:          n.Vault,
		Country:        n.Country,
		Province:       n.Province,
		Region:         n.Region,
		Finca:          n.Finca,
		Description:    n.Description,
		Classification: n.Classification,
		LocalName:      n.LocalName,
		WikiPage:       n.WikiPage,
	}
}

// transformFruitPortOuttoFruit transforms the given fruit port out to service fruit.
func transformFruitPortOuttoFruit(fruitRepo *repository.Fruit) *Fruit {
	if fruitRepo == nil {
		return nil
	}
	newfruit := Fruit{
		ID:             repository.FruitIDValue(fruitRepo.ID),
		Name:           fruitRepo.Name,
		Variety:        fruitRepo.Variety,
		Year:           fruitRepo.Year,
		Price:          repository.FruitPriceValue(fruitRepo.Price),
		Vault:          fruitRepo.Vault,
		Country:        fruitRepo.Country,
		Province:       fruitRepo.Province,
		Region:         fruitRepo.Region,
		Finca:          fruitRepo.Finca,
		Description:    fruitRepo.Description,
		Classification: fruitRepo.Classification,
		LocalName:      fruitRepo.LocalName,
		WikiPage:       fruitRepo.WikiPage,
	}
	return &newfruit
}

// transformFruitPortOuttoFruitItem transforms the given fruit port out to service fruit item.
func transformFruitPortOuttoFruitItem(fruitRepo *repository.Fruit) *FruitItem {
	if fruitRepo == nil {
		return nil
	}
	newfruit := FruitItem{
		ID:   repository.FruitIDValue(fruitRepo.ID),
		Name: fruitRepo.Name,
	}
	return &newfruit
}

// newGetFruitWithIDResult create a new GetFruitWithIDResult
func newGetFruitWithIDResult(fruit *Fruit, err error) GetFruitWithIDResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	if fruit == nil && err == nil {
		errmessage = "record not found"
	}
	return GetFruitWithIDResult{
		Fruit: fruit,
		Err:   errmessage,
	}
}

// newSearchFruitsResult create a new SearchFruitsResult
func newSearchFruitsDataResult(result *SearchFruitsResult, err error) SearchFruitsDataResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return SearchFruitsDataResult{
		SearchResult: result,
		Err:          errmessage,
	}
}

// newCreateFruitResult create a new CreateFruitResponse
func newCreateFruitResult(id int64, err error) CreateFruitResult {
	var errmessage string
	if err != nil {
		errmessage = err.Error()
	}
	return CreateFruitResult{
		ID:  id,
		Err: errmessage,
	}
}

func (u Fruit) String() string {
	b, err := json.Marshal(u)
	if err != nil {
		return ""
	}
	return string(b)
}

func (n NewFruit) String() string {
	b, err := json.Marshal(n)
	if err != nil {
		return ""
	}
	return string(b)
}

func (s SearchFruitFilter) toRepositoryFilters() repository.FruitFilter {
	return repository.FruitFilter{
		Start: s.Start,
		Count: s.Count,
	}
}

func toSearchFruitsResult(repoResult repository.FindFruitsResult) SearchFruitsResult {
	var fruitCollection []FruitItem
	for _, v := range repoResult.Fruits {
		fruitFound := &v
		fruitToAdd := transformFruitPortOuttoFruitItem(fruitFound)
		fruitCollection = append(fruitCollection, *fruitToAdd)
	}
	return SearchFruitsResult{
		Fruits: fruitCollection,
		Total:  repoResult.Total,
		Start:  repoResult.Start,
		Count:  repoResult.Count,
	}
}

func toDatasetStatus(datasetStatus repository.FruitDatasetStatus) DatasetStatus {
	status := DatasetStateOK
	if !datasetStatus.Ok {
		status = DatasetStateError
	}
	return DatasetStatus{
		Message:   datasetStatus.Message,
		Status:    status,
		Timestamp: time.Now().Unix(),
	}
}
