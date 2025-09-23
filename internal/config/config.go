package config

import (
	"encoding/json"
	"os"
	"path"

	"github.com/PaleBlueDot1990/gator/internal/database"
)

const configFileName = ".gatorconfig.json"

type State struct {
	DbQueries *database.Queries
	Cfg *Config
}

type Config struct {
	DB_URL            string `json:"db_url"`
	CURRENT_USER_NAME string `json:"current_user_name"`
}

func Read() (*Config, error) {
	homeDir, err  := os.UserHomeDir()
	if err != nil {
		return nil, err 
	}

	configFilePath := path.Join(homeDir, configFileName)
	rawBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, err
	}

	cfg := Config {}
	err = json.Unmarshal(rawBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func (cfg* Config) SetUser(curr_user_name string) error {
	homeDir, err  := os.UserHomeDir()
	if err != nil {
		return err 
	}

	configFilePath := path.Join(homeDir, configFileName)
	cfg.CURRENT_USER_NAME = curr_user_name
	rawBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}

	return os.WriteFile(configFilePath, rawBytes, 0666)
}
