package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type CredentialsConfig struct {
	File    string `json:"file"`    // location of AWS credentials file
	Profile string `json:"profile"` // AWS credentials profile used by STS client
}

type RolesConfig map[string]string // collection of role alias->ARN pairs

type Config struct {
	DefaultRole string            `json:"default_role"` // role served as instance profile on startup
	Credentials CredentialsConfig `json:"credentials"`
	Roles       RolesConfig       `json:"roles"`
}

func LoadConfig(file string) (*Config, error) {
	cf, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %s", err)
	}

	var c *Config

	decoder := json.NewDecoder(cf)
	if err := decoder.Decode(&c); err != nil {
		return nil, fmt.Errorf("failed to decode %s: %s", file, err)
	}

	return c, nil
}

func (c *Config) String() string {
	config, _ := json.MarshalIndent(c, "", "  ")
	return string(config[:])
}
