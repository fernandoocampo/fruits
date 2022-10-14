package fruits_test

import (
	"strconv"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/fernandoocampo/fruits/internal/fruits"
	"github.com/stretchr/testify/assert"
)

func TestToFruitPortOut(t *testing.T) {
	t.Parallel()

	expectedRepoFruit := repository.NewFruit{
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OaKeefe",
		WikiPage:       "@kerinokeefe",
	}
	givenFruit := fruits.NewFruit{
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OaKeefe",
		WikiPage:       "@kerinokeefe",
	}

	got := givenFruit.ToFruitPortOut()

	assert.Equal(t, expectedRepoFruit, got)
}

func TestNewFruit(t *testing.T) {
	t.Parallel()

	expectedFruit := fruits.Fruit{
		ID:             strconv.Itoa(1234),
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OaKeefe",
		WikiPage:       "@kerinokeefe",
	}
	givenFruitID := "1234"
	givenNewFruit := fruits.NewFruit{
		Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
		Variety:        "White Blend",
		Vault:          "Nicosia",
		Year:           87,
		Country:        "Italy",
		Province:       "Sicily & Sardinia",
		Region:         "Etna",
		Description:    "brisk acidity",
		Classification: "Vulka Bianco",
		LocalName:      "Kerin OaKeefe",
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
				Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
				Variety:        "White Blend",
				Vault:          "Nicosia",
				Year:           87,
				Country:        "Italy",
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulka Bianco",
				LocalName:      "Kerin OaKeefe",
				WikiPage:       "@kerinokeefe",
			},
		},
		"valid_without_more_fields": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
				Vault:          "Nicosia",
				Country:        "Italy",
				Classification: "Vulka Bianco",
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
				Classification: "Vulka Bianco",
				LocalName:      "Kerin OaKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: fruits.MandatoryError{
				Fields: []string{"name"},
			},
		},
		"invalid_classification": {
			data: fruits.NewFruit{
				Name:        "Nicosia 2013 Vulka Bianco  (Etna)",
				Variety:     "White Blend",
				Vault:       "Nicosia",
				Year:        87,
				Country:     "Italy",
				Province:    "Sicily & Sardinia",
				Region:      "Etna",
				Description: "brisk acidity",
				LocalName:   "Kerin OaKeefe",
				WikiPage:    "@kerinokeefe",
			},
			want: fruits.MandatoryError{
				Fields: []string{"classification"},
			},
		},
		"invalid_country": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
				Variety:        "White Blend",
				Vault:          "Nicosia",
				Year:           87,
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulka Bianco",
				LocalName:      "Kerin OaKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: fruits.MandatoryError{
				Fields: []string{"country"},
			},
		},
		"invalid_vault": {
			data: fruits.NewFruit{
				Name:           "Nicosia 2013 Vulka Bianco  (Etna)",
				Variety:        "White Blend",
				Year:           87,
				Country:        "Italy",
				Province:       "Sicily & Sardinia",
				Region:         "Etna",
				Description:    "brisk acidity",
				Classification: "Vulka Bianco",
				LocalName:      "Kerin OaKeefe",
				WikiPage:       "@kerinokeefe",
			},
			want: fruits.MandatoryError{
				Fields: []string{"vault"},
			},
		},
		"all_invalid": {
			data: fruits.NewFruit{
				Variety:     "White Blend",
				Year:        87,
				Province:    "Sicily & Sardinia",
				Region:      "Etna",
				Description: "brisk acidity",
				LocalName:   "Kerin OaKeefe",
				WikiPage:    "@kerinokeefe",
			},
			want: fruits.MandatoryError{
				Fields: []string{"name", "classification", "country", "vault"},
			},
		},
	}

	for name, test := range cases {
		name, test := name, test
		t.Run(name, func(st *testing.T) {
			st.Parallel()

			got := test.data.Validate()
			assert.Equal(st, test.want, got)
		})
	}
}
