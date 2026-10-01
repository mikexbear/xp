package xray

import (
	"encoding/json"

	"github.com/xtls/xray-core/infra/conf"
)

type Meta struct {
	ServerDescription string `yaml:"serverDescription"`
}

type Config struct {
	conf.Config
	Remarks string `json:"remarks"`
	Meta    *Meta  `json:"meta"`
	Raw     string `json:"-"`
}

func NewConfig(raw string) (*Config, error) {
	var result Config
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, err
	}

	result.Raw = raw

	return &result, nil
}

func NewConfigs(raw []string) ([]*Config, error) {
	result := make([]*Config, len(raw))

	for i, r := range raw {
		cfg, err := NewConfig(r)
		if err != nil {
			return nil, err
		}

		result[i] = cfg
	}

	return result, nil
}
