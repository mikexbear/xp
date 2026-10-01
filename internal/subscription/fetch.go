package subscription

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"xp/internal/config"
)

func fetchRawXrayConfigs(subscriptionUrl config.Url) ([]string, error) {
	url := subscriptionUrl.String()

	response, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", response.Status)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var rawConfigs []json.RawMessage
	if err := json.Unmarshal(body, &rawConfigs); err != nil {
		return nil, err
	}

	result := make([]string, len(rawConfigs))
	for i, rawConfig := range rawConfigs {
		result[i] = string(rawConfig)
	}

	return result, nil
}
