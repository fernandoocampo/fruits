package configurations

import (
	"fmt"

	"github.com/caarlos0/env"
)

// Application contains data related to application configuration parameters.
type Application struct {
	Version               string `env:"VERSION" envDefault:"local"`
	CommitHash            string `env:"COMMIT_HASH" envDefault:"local"`
	LogLevel              int    `env:"LOG_LEVEL" envDefault:"2"` // 1 debug
	ApplicationPort       string `env:"APPLICATION_PORT" envDefault:":8080"`
	MetricsIntervalMillis int    `env:"METRICS_INTERVAL_MILLIS" envDefault:"60000"`
	CloudRegion           string `env:"CLOUD_REGION" envDefault:"us-east-1"`
	CloudEndpointURL      string `env:"CLOUD_ENDPOINT_URL" envDefault:"aws"`
}

// Load load application configuration.
func Load() (Application, error) {
	cfg := new(Application)
	if err := env.Parse(cfg); err != nil {
		return *cfg, fmt.Errorf("something went wrong loading app configuration: %w", err)
	}

	return *cfg, nil
}
