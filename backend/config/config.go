package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"
)

type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	Upload   UploadConfig   `yaml:"upload"`
}

type ServerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type JWTConfig struct {
	SecretKey string        `yaml:"secret_key"`
	ExpiresIn time.Duration `yaml:"expires_in"`
}

type UploadConfig struct {
	MaxFileSize int64  `yaml:"max_file_size"`
	TempDir     string `yaml:"temp_dir"`
}

func Load() (*Config, error) {
	config := &Config{}

	if err := loadFromFile(config); err != nil {
		return nil, err
	}

	if err := validate(config); err != nil {
		return nil, err
	}

	return config, nil
}

func loadFromFile(config *Config) error {
	configPath := "./config/config.yml"
	file, err := os.Open(configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	return decoder.Decode(config)
}

func validate(config *Config) error {
	if config.Database.Password == "" {
		return fmt.Errorf("database password is required")
	}

	if config.JWT.SecretKey == "" || config.JWT.SecretKey == "your-secret-key" {
		return fmt.Errorf("JWT secret key must be set and not use default value")
	}

	if config.Upload.MaxFileSize <= 0 {
		return fmt.Errorf("max file size must be greater than 0")
	}

	return nil
}

func (c *Config) GetDatabaseDSN() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}
