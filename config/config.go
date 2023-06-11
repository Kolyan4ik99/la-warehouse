package config

import (
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const configPath = "./config/config.yml"

type Config struct {
	Postgres `yaml:"Postgres"`
}

type Postgres struct {
	Host     string `yaml:"Host"`
	Port     string `yaml:"Port"`
	User     string `yaml:"User"`
	Password string `yaml:"Password"`
	DBName   string `yaml:"DBName"`
	SSLMode  string `yaml:"SSLMode"`
	Settings struct {
		MaxOpenConns    int           `yaml:"MaxOpenConns"`
		ConnMaxLifeTime time.Duration `yaml:"ConnMaxLifeTime"`
		MaxIdleConns    int           `yaml:"MaxIdleConns"`
		MaxIdleLifeTime time.Duration `yaml:"MaxIdleLifeTime"`
	}
}

func LoadConfig() (*Config, error) {
	buf, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(buf, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
