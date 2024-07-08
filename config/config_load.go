package config

import (
	"math/big"

	"github.com/spf13/viper"
)

type Config struct {
	G           *big.Int
	H           *big.Int
	Q           *big.Int
	VerifierURL string
}

// LoadConfig loads the configuration settings from environment variables using Viper.
// It sets default values for the configuration options if the corresponding environment variable is not set.
// The function returns a pointer to a Config struct that contains the loaded configuration values.
func LoadConfig() *Config {
	viper.SetEnvPrefix("zkp")

	_ = viper.BindEnv("verifier_url")
	viper.SetDefault("verifier_url", "localhost:50051")

	_ = viper.BindEnv("g")
	viper.SetDefault("g", "2")

	_ = viper.BindEnv("h")
	viper.SetDefault("h", "5")

	_ = viper.BindEnv("q")
	viper.SetDefault("q", "100")

	return &Config{
		G:           big.NewInt(viper.GetInt64("g")),
		H:           big.NewInt(viper.GetInt64("h")),
		Q:           big.NewInt(viper.GetInt64("q")),
		VerifierURL: viper.GetString("verifier_url"),
	}
}
