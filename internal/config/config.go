package config

import (
	"encoding/json"
	"github.com/nexryai/MateSSH/internal/hostkey"
	"os"
)

const jsonFilePath = "mate_ssh.json"

type ServerConfig struct {
	BindAddr       string          `json:"bind"`
	Port           int             `json:"port"`
	AuthorizedKeys []string        `json:"authorized_keys"`
	HostKeys       hostkey.Keyring `json:"host_keys"`
}

func IsExist() bool {
	_, err := os.Stat(jsonFilePath)
	return err == nil
}

func CreateConfig(HostKeyring hostkey.Keyring, AuthorizedKey string) error {
	config := ServerConfig{
		BindAddr: "0.0.0.0",
		Port:     2222,
		AuthorizedKeys: []string{
			AuthorizedKey,
		},
		HostKeys: HostKeyring,
	}

	data, err := json.Marshal(config)
	if err != nil {
		return err
	}

	f, err := os.OpenFile(jsonFilePath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.Write(data)
	return err
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
