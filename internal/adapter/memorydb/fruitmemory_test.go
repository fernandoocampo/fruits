package memorydb_test

import (
	"bufio"
	"bytes"
	"context"
	"testing"

	"github.com/fernandoocampo/fruits/internal/adapter/loggers"
	"github.com/fernandoocampo/fruits/internal/adapter/memorydb"
	"github.com/fernandoocampo/fruits/internal/adapter/repository"
	"github.com/stretchr/testify/assert"
)

func TestSaveFruit(t *testing.T) {
	newFruitID := int64(1)
	newFruit := repository.NewFruit{
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
	expectedFruit := repository.Fruit{
		ID:             repository.FruitID(newFruitID),
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
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitDB := memorydb.NewFruitRepository(logger)
	ctx := context.TODO()

	fruitID, err := fruitDB.Save(ctx, newFruit)
	savedFruit, readErr := fruitDB.FindByID(ctx, fruitID)

	assert.NoError(t, err)
	assert.NoError(t, readErr)
	assert.Equal(t, &expectedFruit, savedFruit)
}

func TestLoadFruitDataset(t *testing.T) {
	givenDataset := datasetFixture()
	expectedDatasetStatus := repository.FruitDatasetStatus{
		Ok: true,
	}
	expectedDataset := []repository.Fruit{
		{
			ID:             repository.FruitID(129969),
			Name:           "Domaine Marcel Deiss 2012 Pinot Gris (Alsace)",
			Variety:        "Pinot Gris",
			Year:           90,
			Price:          repository.FruitPrice(32.0),
			Vault:          "Domaine Marcel Deiss",
			Country:        "France",
			Province:       "Alsace",
			Region:         "Alsace",
			Description:    "A dry style of Pinot Gris.",
			Classification: "",
			LocalName:      "Roger Voss",
			WikiPage:       "@vossroger",
		},
		{
			ID:             repository.FruitID(129970),
			Name:           "Domaine Schoffit 2012 Lieu-dit Harth Cuvée Caroline Gewurztraminer (Alsace)",
			Variety:        "Gewürztraminer",
			Year:           90,
			Price:          repository.FruitPrice(21.0),
			Vault:          "Domaine Schoffit",
			Country:        "France",
			Province:       "Alsace",
			Region:         "Alsace",
			Description:    "Big, rich and off-dry.",
			Classification: "Lieu-dit Harth Cuvée Caroline",
			LocalName:      "Roger Voss",
			WikiPage:       "@vossroger",
		},
	}
	logger := loggers.NewLoggerWithStdout("", loggers.Debug)
	fruitRepo := memorydb.NewFruitRepository(logger)
	filter := repository.FruitFilter{
		Start: 1,
		Count: 10,
	}
	ctx := context.TODO()
	datasetBuffer := bytes.NewBuffer([]byte(givenDataset))
	scanner := bufio.NewScanner(datasetBuffer)

	err := fruitRepo.LoadFruitDataset(ctx, scanner)

	assert.NoError(t, err)
	result, err := fruitRepo.SearchWithFilters(ctx, filter)
	assert.NoError(t, err)
	assert.Equal(t, len(expectedDataset), result.Total)
	assert.Equal(t, expectedDataset, result.Fruits)
	status, err := fruitRepo.DatasetStatus(ctx)
	assert.NoError(t, err)
	assert.Equal(t, expectedDatasetStatus, status)
}

func datasetFixture() string {
	return `,country,description,classification,year,price,province,region,finca,local_name,wiki_page,name,variety,vault
129969,France,"A dry style of Pinot Gris.",,90,32.0,Alsace,Alsace,,Roger Voss,@vossroger,Domaine Marcel Deiss 2012 Pinot Gris (Alsace),Pinot Gris,Domaine Marcel Deiss
129970,France,"Big, rich and off-dry.",Lieu-dit Harth Cuvée Caroline,90,21.0,Alsace,Alsace,,Roger Voss,@vossroger,Domaine Schoffit 2012 Lieu-dit Harth Cuvée Caroline Gewurztraminer (Alsace),Gewürztraminer,Domaine Schoffit
`
}
