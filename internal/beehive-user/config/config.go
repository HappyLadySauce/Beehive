package config

import "github.com/HappyLadySauce/Beehive/internal/beehive-user/options"

type Config struct {
	*options.Options
}

func CreateConfigFromOptions(opts *options.Options) (*Config, error) {
	return &Config{opts}, nil
}