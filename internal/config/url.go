package config

import (
	"net/url"

	"gopkg.in/yaml.v3"
)

type Url struct {
	url.URL
}

func (u *Url) UnmarshalYAML(node *yaml.Node) error {
	var s string
	if err := node.Decode(&s); err != nil {
		return err
	}

	parsed, err := url.Parse(s)
	if err != nil {
		return err
	}

	u.URL = *parsed

	return nil
}
