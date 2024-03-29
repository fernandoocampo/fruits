package fruits

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/fernandoocampo/fruits/internal/adapter/repository"
)

// FruitService defines behavior for fruit service business logic.
type FruitService interface {
	GetFruitWithID(ctx context.Context, fruitID string) (*Fruit, error)
	Create(ctx context.Context, newfruit NewFruit) (string, error)
	SearchFruits(ctx context.Context, givenFilter SearchFruitFilter) (*SearchFruitsResult, error)
	DatasetStatus(ctx context.Context) DatasetStatus
}

// MandatoryError define an error for mandatory fields.
type MandatoryError struct {
	Fields []string
}

// DatasetState define fruit dataset state.
type DatasetState string

// CreateFruitResult standard response for create Fruit.
type CreateFruitResult struct {
	ID  string
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

// SearchFruitFilter contains filters to search fruits.
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
	ID             string  `json:"id"`
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
	ID string `json:"id"`
	// Name fruit's name.
	Name string `json:"name"`
}

// DatasetStatus contains data about the fruit dataset result.
type DatasetStatus struct {
	Status    DatasetState
	Message   string
	Timestamp int64
}

const (
	DatasetStateOK    DatasetState = "ok"
	DatasetStateError DatasetState = "error"
)

func (m MandatoryError) Error() string {
	return fmt.Sprintf(
		"these fields are mandatory: %s.",
		strings.Join(m.Fields, ", "),
	)
}

// ToFruitPortOut transforms new fruit to a fruit port out.
func (n NewFruit) ToFruitPortOut() repository.NewFruit {
	return repository.NewFruit{
		Name:           n.Name,
		Variety:        n.Variety,
		Year:           n.Year,
		Price:          repository.FruitPrice(n.Price),
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

// Validate check is the given data to create a fruit is correct.
func (n NewFruit) Validate() error {
	var invalidList []string

	if n.Name == "" {
		invalidList = append(invalidList, "name")
	}

	if n.Classification == "" {
		invalidList = append(invalidList, "classification")
	}

	if n.Country == "" {
		invalidList = append(invalidList, "country")
	}

	if n.Vault == "" {
		invalidList = append(invalidList, "vault")
	}

	if len(invalidList) == 0 {
		return nil
	}

	return MandatoryError{
		Fields: invalidList,
	}
}

// NewFruit transforms new fruit to a fruit port out.
func (n NewFruit) NewFruit(fruitID string) Fruit {
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

// newGetFruitWithIDResult create a new GetFruitWithIDResult.
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

// newSearchFruitsResult create a new SearchFruitsResult.
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

// newCreateFruitResult create a new CreateFruitResponse.
func newCreateFruitResult(fruitID string, err error) CreateFruitResult {
	var errmessage string

	if err != nil {
		errmessage = err.Error()
	}

	return CreateFruitResult{
		ID:  fruitID,
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
	fruitCollection := make([]FruitItem, len(repoResult.Fruits))

	for index := range repoResult.Fruits {
		fruitFound := &repoResult.Fruits[index]
		fruitToAdd := transformFruitPortOuttoFruitItem(fruitFound)
		fruitCollection[index] = *fruitToAdd
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
