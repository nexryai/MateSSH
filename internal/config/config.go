package config

import (
	"encoding/json"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"os"
)

const jsonFilePath = "~/mate_ssh.json"

type ServerConfig struct {
	BindAddr       string          `json:"bind"`
	Port           int             `json:"port"`
	AuthorizedKeys []string        `json:"authorized_keys"`
	HostKeys       hostkey.Keyring `json:"host_keys"`
}

func CreateConfig(HostKeyring hostkey.Keyring) error {
	config := ServerConfig{
		BindAddr: "0.0.0.0",
		Port:     2222,
		HostKeys: HostKeyring,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(jsonFilePath, data, 0600)
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

	err = os.WriteFile(jsonFilePath, data, 0600)
	if err != nil {
		return err
	}

	return nil
}

func LoadConfig() (*ServerConfig, error) {
	data, err := os.ReadFile(jsonFilePath)
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
