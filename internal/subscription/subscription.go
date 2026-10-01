package subscription

import (
	"sync"
	"xp/internal/config"
	"xp/internal/xray"
)

type Subscription struct {
	xrayConfigs []*xray.Config
	mutex       sync.Mutex
	url         config.Url
}

func NewSubscription(url config.Url) (*Subscription, error) {
	xrayConfigs, err := buildXrayConfigs(url)
	if err != nil {
		return nil, err
	}

	return &Subscription{xrayConfigs: xrayConfigs, url: url}, nil
}

func NewSubscriptions(urls []config.Url) ([]*Subscription, error) {
	result := make([]*Subscription, len(urls))

	for i, url := range urls {
		sub, err := NewSubscription(url)
		if err != nil {
			return nil, err
		}

		result[i] = sub
	}

	return result, nil
}

func (s *Subscription) GetXrayConfigs() []*xray.Config {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	result := make([]*xray.Config, len(s.xrayConfigs))
	copy(result, s.xrayConfigs)

	return result
}

func (s *Subscription) UpdateXrayConfigs() error {
	xrayConfigs, err := buildXrayConfigs(s.url)
	if err != nil {
		return err
	}

	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.xrayConfigs = xrayConfigs

	return nil
}

func buildXrayConfigs(url config.Url) ([]*xray.Config, error) {
	raw, err := fetchRawXrayConfigs(url)
	if err != nil {
		return nil, err
	}

	xrayConfigs, err := xray.NewConfigs(raw)
	if err != nil {
		return nil, err
	}

	return xrayConfigs, nil
}
