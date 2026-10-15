package xray

import (
	"encoding/json"
	"fmt"

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

func (c *Config) GetSocksInboundHostAndPort() (string, int, error) {
	socksInbound := c.getSocksInbound()
	if socksInbound == nil {
		return "", 0, fmt.Errorf("no SOCKS inbound found in Xray configuration")
	}

	if socksInbound.PortList == nil || len(socksInbound.PortList.Range) == 0 {
		return "", 0, fmt.Errorf("SOCKS inbound has no port")
	}
	port := socksInbound.PortList.Range[0].From

	host := "127.0.0.1"
	if socksInbound.ListenOn != nil && socksInbound.ListenOn.String() != "" {
		host = socksInbound.ListenOn.String()
	}

	return host, int(port), nil
}

func (c *Config) getSocksInbound() *conf.InboundDetourConfig {
	for _, inbound := range c.InboundConfigs {
		if inbound.Protocol == "socks" {
			return &inbound
		}
	}

	return nil
}
