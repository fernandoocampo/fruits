package configurations

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Application contains data related to application configuration parameters.
type Application struct {
	LogLevel              int    `env:"LOG_LEVEL" envDefault:"2"` // 1 debug
	ApplicationPort       string `env:"APPLICATION_PORT" envDefault:":8080"`
	MetricsIntervalMillis int    `env:"METRICS_INTERVAL_MILLIS" envDefault:"60000"`
	FilePath              string `env:"FILE_PATH" envDefault:"/opt/fruits/fruitmag-data.csv"`
	LoadDataset           bool   `env:"LOAD_DATASET" envDefault:"true"`
}

// Load load application configuration.
func Load() (Application, error) {
	cfg := new(Application)
	if err := env.Parse(cfg); err != nil {
		return *cfg, fmt.Errorf("something went wrong loading app configuration: %w", err)
	}

	return *cfg, nil
}
