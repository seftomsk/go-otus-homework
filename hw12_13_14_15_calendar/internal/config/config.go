package config

import (
	"errors"
	"fmt"
	"strings"
)

var ErrNotSetValue = errors.New("not set value")

var ErrInvalidValue = errors.New("invalid value")

var supportedStorages = map[string]struct{}{
	"memory":   {},
	"postgres": {},
}

var supportedLogLevels = map[string]struct{}{
	"error": {},
	"warn":  {},
	"info":  {},
	"debug": {},
}

var supportedEnvironments = map[string]struct{}{
	"dev":  {},
	"prod": {},
}

func New() Config {
	return Config{}
}

type Server struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

type Database struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Logger struct {
	LogLevel string `json:"loglevel"`
	FileName string `json:"filename"`
	StdErr   bool   `json:"stderr"`
}

type Config struct {
	Version     string   `json:"version"`
	Environment string   `json:"environment"`
	Storage     string   `json:"storage"`
	Server      Server   `json:"server"`
	Database    Database `json:"database"`
	Logger      Logger   `json:"logger"`
}

func (c Config) Host() string {
	return c.Server.Host
}

func (c Config) Port() string {
	return c.Server.Port
}

func validateSetValues(c Config) error {
	if c.Version == "" {
		return fmt.Errorf("validateSetValues - version: %w", ErrNotSetValue)
	}
	if c.Environment == "" {
		return fmt.Errorf("validateSetValues - environment: %w", ErrNotSetValue)
	}
	if c.Storage == "" {
		return fmt.Errorf("validateSetValues - storage: %w", ErrNotSetValue)
	}
	if c.Server.Host == "" {
		return fmt.Errorf("validateSetValues - host: %w", ErrNotSetValue)
	}
	if c.Server.Port == "" {
		return fmt.Errorf("validateSetValues - port: %w", ErrNotSetValue)
	}
	if c.Logger.FileName == "" {
		return fmt.Errorf("validateSetValues - logger - filename: %w", ErrNotSetValue)
	}
	if c.Storage != "postgres" {
		return nil
	}
	if c.Database.Host == "" {
		return fmt.Errorf("validateSetValues - database - host: %w", ErrNotSetValue)
	}
	if c.Database.Port == "" {
		return fmt.Errorf("validateSetValues - database - port: %w", ErrNotSetValue)
	}
	if c.Database.Name == "" {
		return fmt.Errorf("validateSetValues - database - name: %w", ErrNotSetValue)
	}
	if c.Database.Username == "" {
		return fmt.Errorf("validateSetValues - database - username: %w", ErrNotSetValue)
	}
	if c.Database.Password == "" {
		return fmt.Errorf("validateSetValues - database - password: %w", ErrNotSetValue)
	}

	return nil
}

func validateSupportedValues(c Config) error {
	if _, ok := supportedEnvironments[strings.ToLower(c.Environment)]; !ok {
		return fmt.Errorf("validateSupportedValues - environment: %w", ErrInvalidValue)
	}
	if _, ok := supportedStorages[strings.ToLower(c.Storage)]; !ok {
		return fmt.Errorf("validateSupportedValues - storage: %w", ErrInvalidValue)
	}
	if _, ok := supportedLogLevels[strings.ToLower(c.Logger.LogLevel)]; !ok {
		return fmt.Errorf("validateSupportedValues - logger - log_level: %w", ErrInvalidValue)
	}

	return nil
}

func (c Config) Validate() error {
	err := validateSetValues(c)
	if err != nil {
		return fmt.Errorf("validate - validateSetValies: %w", err)
	}

	err = validateSupportedValues(c)
	if err != nil {
		return fmt.Errorf("validate - validateSupportedValues: %w", err)
	}

	return nil
}
