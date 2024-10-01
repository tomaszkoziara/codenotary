package config

import (
	"fmt"

	env "github.com/caarlos0/env/v11"
)

type config struct {
	Port  int `env:"API_PORT" envDefault:"8000"`
	Vault Vault
}

type Vault struct {
	Ledger     string `env:"VAULT_LEDGER" envDefault:"default"`
	Collection string `env:"VAULT_COLLECTION" envDefault:"default"`
	APIKey     string `env:"VAULT_API_KEY" envDefault:"default"`
}

func Load() (config, error) {
	var c config
	err := env.Parse(&c)
	if err != nil {
		return config{}, fmt.Errorf("failed to load config: %w", err)
	}
	return c, nil
}
