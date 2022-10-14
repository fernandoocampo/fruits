package repository

// FruitID is the fruit identification type.
type FruitID string

// Fruit contains fruit data.
type Fruit struct {
	ID             FruitID  `json:"id"`
	Name           string   `json:"name"`
	Variety        string   `json:"variety"`
	Year           int      `json:"year"`
	Price          *float32 `json:"price,omitempty"`
	Vault          string   `json:"vault"`
	Country        string   `json:"country"`
	Province       string   `json:"province"`
	Region         string   `json:"region,omitempty"`
	Finca          string   `json:"finca,omitempty"`
	Description    string   `json:"description"`
	Classification string   `json:"classification"`
	LocalName      string   `json:"local_name"`
	WikiPage       string   `json:"wiki_page,omitempty"`
}

// NewFruit contains data to create a new fruit.
type NewFruit struct {
	Name           string   `json:"name"`
	Variety        string   `json:"variety"`
	Year           int      `json:"year"`
	Price          *float32 `json:"price,omitempty"`
	Vault          string   `json:"vault"`
	Country        string   `json:"country"`
	Province       string   `json:"province"`
	Region         string   `json:"region,omitempty"`
	Finca          string   `json:"finca,omitempty"`
	Description    string   `json:"description"`
	Classification string   `json:"classification"`
	LocalName      string   `json:"local_name"`
	WikiPage       string   `json:"wiki_page,omitempty"`
}

// FindFruitsResult contains the list of fruits found plus some metadata.
type FindFruitsResult struct {
	Fruits []Fruit
	Total  int
	Start  int
	Count  int
}

// FruitFilter contains filters to search fruits.
type FruitFilter struct {
	// Start record to query
	Start int
	// rows to return
	Count int
}

// FruitDatasetStatus contains data for dataset status.
type FruitDatasetStatus struct {
	Ok      bool
	Message string
}

// FruitPrice returns a pointer to the int value passed in.
func FruitPrice(v float32) *float32 {
	if v == 0 {
		return nil
	}

	return &v
}

// FruitPriceValue returns the value of the int pointer
// passed in or 0 if the pointer is nil.
func FruitPriceValue(v *float32) float32 {
	if v == nil {
		return 0
	}

	return *v
}

// FruitIDValue returns the value of the fruit id value as a int64.
func FruitIDValue(v FruitID) string {
	return string(v)
}

// ToFruit transforms new fruit to a fruit.
func (u NewFruit) ToFruit(fruitID FruitID) Fruit {
	return Fruit{
		ID:             fruitID,
		Name:           u.Name,
		Variety:        u.Variety,
		Year:           u.Year,
		Price:          u.Price,
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
