// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

// config/config.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Versions struct {
	Go           string `mapstructure:"go"`
	Kratos       string `mapstructure:"kratos"`
	Wire         string `mapstructure:"wire"`
	Grpc         string `mapstructure:"grpc"`
	Protobuf     string `mapstructure:"protobuf"`
	Automaxprocs string `mapstructure:"automaxprocs"`
	Genproto     string `mapstructure:"genproto"`
}

type Config struct {
	Name           string
	MonorepoPrefix string   `mapstructure:"monorepo-prefix"`
	OutputDir      string   `mapstructure:"output-dir"`
	Versions       Versions `mapstructure:"versions"`
	DryRun         bool
}

var Cfg *Config

func init() {
	// Set default values
	viper.SetDefault("monorepo-prefix", "github.com/org/repo")
	viper.SetDefault("output-dir", ".")
	viper.SetDefault("versions.go", "1.24.3")
	viper.SetDefault("versions.kratos", "v2.8.4")
	viper.SetDefault("versions.wire", "v0.6.0")
	viper.SetDefault("versions.grpc", "v1.61.1")
	viper.SetDefault("versions.protobuf", "v1.33.0")
	viper.SetDefault("versions.automaxprocs", "v1.6.0")
	viper.SetDefault("versions.genproto", "v0.0.0-20240102182953-50ed04b92917")
}

// Load reads gecro.yaml and unmarshals into Cfg
func Load() error {
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot determine current directory: %w", err)
	}

	cfgPath := filepath.Join(cwd, "gecro.yaml")
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		return fmt.Errorf(
			"gecro.yaml not found at %q — please run this command from the monorepo root",
			cwd,
		)
	} else if err != nil {
		return fmt.Errorf("error checking for gecro.yaml at %q: %w", cfgPath, err)
	}

	viper.SetConfigFile(cfgPath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("failed to read config: %w", err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("unable to decode config into struct: %w", err)
	}
	Cfg = &cfg
	return nil
}
