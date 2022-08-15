package fruits_test

import (
	"errors"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/stretchr/testify/assert"
)

func TestToFruitPortOut(t *testing.T) {
	t.Parallel()

	expectedRepoFruit := repository.NewFruit{
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}
	givenFruit := fruits.NewFruit{
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}

	got := givenFruit.ToFruitPortOut()

	assert.Equal(t, expectedRepoFruit, got)
}

func TestNewFruit(t *testing.T) {
	t.Parallel()

	expectedFruit := fruits.Fruit{
		ID:             int64(1234),
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}
	givenFruitID := int64(1234)
	givenNewFruit := fruits.NewFruit{
		Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulk√† Bianco",
		LocalName:      "Kerin OÄôKeefe",
		WikiPage:       "@kerinokeefe",
	}

	got := givenNewFruit.NewFruit(givenFruitID)

	assert.Equal(t, expectedFruit, got)
}

func TestNewValidation(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		data fruits.NewFruit
		want error
	}{
		"valid": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
				Variety:        "White Blend",
				Vault:          "Nicosia",
				Year:           87,
				Country:        "Italy",
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulk√† Bianco",
				LocalName:      "Kerin OÄôKeefe",
				WikiPage:       "@kerinokeefe",
			},
		},
		"valid_without_more_fields": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
				Vault:          "Nicosia",
				Country:        "Italy",
				Classification: "Vulk√† Bianco",
			},
		},
		"invalid_name": {
			data: fruits.NewFruit{
				Variety:        "White Blend",
				Vault:          "Nicosia",
				Year:           87,
				Country:        "Italy",
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulk√† Bianco",
				LocalName:      "Kerin OÄôKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: errors.New("these fields are mandatory: name."),
		},
		"invalid_classification": {
			data: fruits.NewFruit{
				Name:        "Nicosia 2013 Vulk√† Bianco  (Etna)",
				Variety:     "White Blend",
				Vault:       "Nicosia",
				Year:        87,
				Country:     "Italy",
				Province:    "Sicily & Sardinia",
				Region:      "Etna",
				Description: "brisk acidity",
				LocalName:   "Kerin OÄôKeefe",
				WikiPage:    "@kerinokeefe",
			},
			want: errors.New("these fields are mandatory: classification."),
		},
		"invalid_country": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
				Variety:        "White Blend",
				Vault:          "Nicosia",
				Year:           87,
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulk√† Bianco",
				LocalName:      "Kerin OÄôKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: errors.New("these fields are mandatory: country."),
		},
		"invalid_vault": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulk√† Bianco  (Etna)",
				Variety:        "White Blend",
				Year:           87,
				Country:        "Italy",
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulk√† Bianco",
				LocalName:      "Kerin OÄôKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: errors.New("these fields are mandatory: vault."),
		},
		"all_invalid": {
			data: fruits.NewFruit{
				Variety:     "White Blend",
				Year:        87,
				Province:    "Sicily & Sardinia",
				Region:      "Etna",
				Description: "brisk acidity",
				LocalName:   "Kerin OÄôKeefe",
				WikiPage:    "@kerinokeefe",
			},
			want: errors.New("these fields are mandatory: name, classification, country, vault."),
		},
	}

	for k, v := range cases {
		got := v.data.Validate()
		assert.Equal(t, v.want, got, k)
	}
}
