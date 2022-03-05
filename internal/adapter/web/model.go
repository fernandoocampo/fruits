package web

import "github.com/fernandoocampo/fruits/internal/fruits"

// Result standard result for the service
type Result struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
	Errors  []string    `json:"errors"`
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

// FruitItemResult contains data related to a fruit found during a search.
type FruitItemResult struct {
	ID int64 `json:"id"`
	// Name or name of the fruit.
	Name string `json:"name"`
}

// NewFruit contains the expected data for a new fruit.
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

// CreateFruitResponse standard response for create Fruit
type CreateFruitResponse struct {
	ID  string `json:"id"`
	Err string `json:"err,omitempty"`
}

// GetFruitWithIDResponse standard response for get a Fruit with an ID.
type GetFruitWithIDResponse struct {
	Fruit *Fruit `json:"fruit"`
	Err   string `json:"err,omitempty"`
}

// SearchFruitsResponse standard response for searching fruits with filters.
type SearchFruitsResponse struct {
	Fruits *SearchFruitsResult `json:"result"`
	Err    string              `json:"err,omitempty"`
}

// SearchFruitFilter contains filters to search fruits
type SearchFruitFilter struct {
	// Start record position to query
	Start int
	// rows per page
	Count int
}

// SearchFruitsResult contains search fruits result data.
type SearchFruitsResult struct {
	Fruits []FruitItemResult `json:"fruits"`
	Total  int               `json:"total"`
	Start  int               `json:"start"`
	Count  int               `json:"count"`
}

// FruitDatasetStatusResponse contains fruit dataset status result data.
type FruitDatasetStatusResponse struct {
	Status    string `json:"status"`
	Message   string `json:"msg"`
	Timestamp int64  `json:"ts"`
}

// toFruit transforms new fruit to a fruit object.
func toFruit(fruit *fruits.Fruit) *Fruit {
	if fruit == nil {
		return nil
	}
	webFruit := Fruit{
		ID:             fruit.ID,
		Name:           fruit.Name,
		Variety:        fruit.Variety,
		Year:           fruit.Year,
		Price:          fruit.Price,
		Vault:          fruit.Vault,
		Country:        fruit.Country,
		Province:       fruit.Province,
		Region:         fruit.Region,
		Finca:          fruit.Finca,
		Description:    fruit.Description,
		Classification: fruit.Classification,
		LocalName:      fruit.LocalName,
		WikiPage:       fruit.WikiPage,
	}
	return &webFruit
}

// toFruitItemResult transforms fruit data to a fruit result.
func toFruitItemResult(fruit *fruits.FruitItem) *FruitItemResult {
	if fruit == nil {
		return nil
	}
	webFruit := FruitItemResult{
		ID:   fruit.ID,
		Name: fruit.Name,
	}
	return &webFruit
}

// toSearchFruitResult transforms new fruit to a fruit object.
func toSearchFruitResult(result *fruits.SearchFruitsResult) *SearchFruitsResult {
	if result == nil {
		return nil
	}
	fruitsFound := make([]FruitItemResult, 0)
	for _, v := range result.Fruits {
		fruitFound := toFruitItemResult(&v)
		fruitsFound = append(fruitsFound, *fruitFound)
	}
	webFruit := SearchFruitsResult{
		Fruits: fruitsFound,
		Total:  result.Total,
		Start:  result.Start,
		Count:  result.Count,
	}
	return &webFruit
}

// toFruit transforms new fruit to a fruit object.
func (n *NewFruit) toFruit() *fruits.NewFruit {
	if n == nil {
		return nil
	}
	fruitDomain := fruits.NewFruit{
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
	return &fruitDomain
}

func toCreateFruitResponse(fruitResult fruits.CreateFruitResult) Result {
	var message Result
	if fruitResult.Err == "" {
		message.Success = true
		message.Data = fruitResult.ID
	}
	if fruitResult.Err != "" {
		message.Errors = []string{fruitResult.Err}
	}
	return message
}

func toGetFruitWithIDResponse(fruitResult fruits.GetFruitWithIDResult) Result {
	var message Result
	newFruit := toFruit(fruitResult.Fruit)
	if fruitResult.Err == "" {
		message.Success = true
		message.Data = newFruit
	}
	if fruitResult.Err != "" {
		message.Errors = []string{fruitResult.Err}
	}
	return message
}

func toSearchFruitsResponse(fruitResult fruits.SearchFruitsDataResult) Result {
	var message Result

	if fruitResult.Err == "" {
		message.Success = true
		message.Data = toSearchFruitResult(fruitResult.SearchResult)
	}
	if fruitResult.Err != "" {
		message.Errors = []string{fruitResult.Err}
	}
	return message
}

func toFruitDatasetStatusResponse(status fruits.DatasetStatus) FruitDatasetStatusResponse {
	response := FruitDatasetStatusResponse{
		Status:    string(status.Status),
		Message:   status.Message,
		Timestamp: status.Timestamp,
	}
	return response
}

func (s SearchFruitFilter) toSearchFruitFilter() fruits.SearchFruitFilter {
	return fruits.SearchFruitFilter{
		Start: s.Start,
		Count: s.Count,
	}
}
