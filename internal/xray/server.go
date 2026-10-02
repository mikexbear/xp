package xray

import (
	"encoding/json"
	"fmt"

	libXray "github.com/xtls/libxray"
)

type Server struct {
	XrayConfig *Config
}

func NewServer(xrayConfig *Config) *Server {
	return &Server{XrayConfig: xrayConfig}
}

func (s *Server) Start() error {
	request, err := json.Marshal(map[string]any{
		"apiVersion": 3,
		"method":     "runXray",
		"payload": map[string]string{
			"xrayJson": s.XrayConfig.Raw,
		},
	})

	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}

	rawResponse := libXray.Invoke(string(request))
	if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("failed to start server: %s", response.Error)
	}

	return nil
}

func (s *Server) Stop() error {
	request, err := json.Marshal(map[string]any{
		"apiVersion": 3,
		"method":     "stopXray",
		"payload":    map[string]any{},
	})

	if err != nil {
		return err
	}

	var response struct {
		Success bool   `json:"success"`
		Error   string `json:"error"`
	}

	rawResponse := libXray.Invoke(string(request))
	if err := json.Unmarshal([]byte(rawResponse), &response); err != nil {
		return err
	}

	if !response.Success {
		return fmt.Errorf("failed to stop server: %s", response.Error)
	}

	return nil
}
