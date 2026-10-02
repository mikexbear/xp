package config

import (
	"errors"
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	SubscriptionUrls                []Url `yaml:"subscriptionUrls"`
	ProbeUrls                       []Url `yaml:"probeUrls"`
	ProbeTimeoutSec                 int   `yaml:"probeTimeoutSec"`
	ProbeIntervalSec                int   `yaml:"probeIntervalSec"`
	SubscriptionFetchingIntervalSec int   `yaml:"subscriptionFetchingIntervalSec"`
	ServerPort                      int   `yaml:"serverPort"`
}

func Default() *Config {
	return &Config{
		SubscriptionUrls:                []Url{},
		ProbeUrls:                       []Url{},
		ProbeTimeoutSec:                 5,
		ProbeIntervalSec:                300,
		SubscriptionFetchingIntervalSec: 900,
		ServerPort:                      8080,
	}
}

func Load(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open configuration file (%s): %w", path, err)
	}

	defer func() { _ = f.Close() }()

	decoder := yaml.NewDecoder(f)
	decoder.KnownFields(true)

	cfg := Default()

	if err := decoder.Decode(cfg); err != nil && !errors.Is(err, io.EOF) {
		return nil, fmt.Errorf("failed to decode configuration (%s): %w", path, err)
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate configuration (%s): %w", path, err)
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if len(c.SubscriptionUrls) == 0 {
		return fmt.Errorf("no subscription URLs configured")
	}

	if len(c.ProbeUrls) == 0 {
		return fmt.Errorf("no probe URLs configured")
	}

	if c.ProbeTimeoutSec <= 0 {
		return fmt.Errorf("probeTimeoutSec must be a positive number (got %d)", c.ProbeTimeoutSec)
	}

	if c.ProbeIntervalSec <= 0 {
		return fmt.Errorf("probeIntervalSec must be a positive number (got %d)", c.ProbeIntervalSec)
	}

	if c.SubscriptionFetchingIntervalSec <= 0 {
		return fmt.Errorf(
			"subscriptionFetchingIntervalSec must be a positive number (got %d)",
			c.SubscriptionFetchingIntervalSec,
		)
	}

	if c.ServerPort <= 0 || c.ServerPort > 65535 {
		return fmt.Errorf("serverPort out of range (%d)", c.ServerPort)
	}

	return nil
}
