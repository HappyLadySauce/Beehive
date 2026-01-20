package config

import "github.com/HappyLadySauce/Beehive/internal/beehive-gateway/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}
