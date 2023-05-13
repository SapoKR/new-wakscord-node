package config

import "github.com/caarlos0/env/v8"

// Initialize load environment variables and save as config struct
func Initialize() error {
	return env.Parse(&Default)
}
