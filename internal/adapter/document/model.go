package document

import "github.com/fernandoocampo/fruits/internal/adapter/repository"

// Fruit contains data to create a new fruit.
type Fruit struct {
	ID             string   `json:"id" dynamodbav:"id"`
	Name           string   `json:"name" dynamodbav:"name"`
	Variety        string   `json:"variety" dynamodbav:"variety"`
	Year           int      `json:"year" dynamodbav:"year"`
	Price          *float32 `json:"price,omitempty" dynamodbav:"price"`
	Vault          string   `json:"vault" dynamodbav:"vault"`
	Country        string   `json:"country" dynamodbav:"country"`
	Province       string   `json:"province" dynamodbav:"province"`
	Region         string   `json:"region,omitempty" dynamodbav:"region"`
	Finca          string   `json:"finca,omitempty" dynamodbav:"finca"`
	Description    string   `json:"description" dynamodbav:"description"`
	Classification string   `json:"classification" dynamodbav:"classification"`
	LocalName      string   `json:"local_name" dynamodbav:"local_name"`
	WikiPage       string   `json:"wiki_page,omitempty" dynamodbav:"wiki_page"`
}

// transformFruit transforms new fruit to a repository fruit.
func (f *Fruit) toRepositoryFruit() *repository.Fruit {
	if f == nil {
		return nil
	}

	return &repository.Fruit{
		ID:             repository.FruitID(f.ID),
		Name:           f.Name,
		Variety:        f.Variety,
		Year:           f.Year,
		Price:          f.Price,
		Vault:          f.Vault,
		Country:        f.Country,
		Province:       f.Province,
		Region:         f.Region,
		Finca:          f.Finca,
		Description:    f.Description,
		Classification: f.Classification,
		LocalName:      f.LocalName,
		WikiPage:       f.WikiPage,
	}
}

// transformFruit transforms new fruit to a fruit.
func transformFruit(fruitID string, fruit repository.NewFruit) Fruit {
	return Fruit{
		ID:             fruitID,
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
}
