// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.trai.ch/gecro/config"
	"go.trai.ch/gecro/internal/generator"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Creates a default gecro.yaml config file",
	RunE: func(cmd *cobra.Command, args []string) error {
		gen, err := generator.NewGenerator()
		if err != nil {
			return fmt.Errorf("failed creating new generator: %w", err)
		}

		force, err := cmd.Flags().GetBool("force")
		if err != nil {
			return err
		}

		cfg := config.Config{
			MonorepoPrefix: viper.GetString("monorepo-prefix"),
			OutputDir:      viper.GetString("output-dir"),
			Versions: config.Versions{
				Go:           viper.GetString("versions.go"),
				Kratos:       viper.GetString("versions.kratos"),
				Wire:         viper.GetString("versions.wire"),
				Grpc:         viper.GetString("versions.grpc"),
				Protobuf:     viper.GetString("versions.protobuf"),
				Automaxprocs: viper.GetString("versions.automaxprocs"),
				Genproto:     viper.GetString("versions.genproto"),
			},
		}

		if err := gen.GenerateConfig(force, cfg); err != nil {
			return fmt.Errorf("failed to initialize gecro project: %w", err)
		}
		fmt.Println("gecro.yaml created successfully.")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("force", "f", false, "Force overwrite of existing gecro.yaml")

	// Define flags in a map to add them in a loop.
	// Defaults are pulled from Viper, which are now reliably set.
	flags := map[string]string{
		"monorepo-prefix":       "Monorepo prefix",
		"output-dir":            "Output directory for generated files",
		"versions.go":           "Go version",
		"versions.kratos":       "Kratos version",
		"versions.wire":         "Wire version",
		"versions.grpc":         "gRPC version",
		"versions.protobuf":     "Protobuf version",
		"versions.automaxprocs": "Automaxprocs version",
		"versions.genproto":     "Genproto version",
	}

	for key, description := range flags {
		initCmd.Flags().String(key, viper.GetString(key), description)
		viper.BindPFlag(key, initCmd.Flags().Lookup(key))
	}
}
