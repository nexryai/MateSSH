package config

import (
	"encoding/json"
	"os"
)

type ServerConfig struct {
	BindAddr       string   `json:"bind"`
	Port           int      `json:"port"`
	AuthorizedKeys []string `json:"authorized_keys"`
	Fingerprint    string   `json:"fingerprint"`
}

func CreateConfig(fingerprint string) error {
	config := ServerConfig{
		BindAddr:    "0.0.0.0",
		Port:        2222,
		Fingerprint: fingerprint,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile("~/mate_ssh.json", data, 0600)
}

func AddAuthorizedKey(key string) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	config.AuthorizedKeys = append(config.AuthorizedKeys, key)

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	err = os.WriteFile("~/mate_ssh.json", data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*ServerConfig, error) {
	data, err := os.ReadFile("~/mate_ssh.json")
	if err != nil {
		return nil, err
	}

	var config ServerConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
