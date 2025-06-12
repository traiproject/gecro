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

var cfg = config.Config{
	MonorepoPrefix: "github.com/org/repo",
	OutputDir:      ".",
	Versions: config.Versions{
		Go:           "1.21",
		Kratos:       "v2.3.0",
		Wire:         "v0.4.0",
		Grpc:         "v1.65.0",
		Protobuf:     "v1.34.1",
		Automaxprocs: "v1.5.1",
		Genproto:     "v0.0.0-20240528184218-531527333157",
	},
}

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

	// Add flags for all config options
	initCmd.Flags().String("monorepo-prefix", "github.com/org/repo", "Monorepo prefix")
	initCmd.Flags().String("output-dir", ".", "Output directory for generated files")
	initCmd.Flags().String("versions.go", "1.22", "Go version")
	initCmd.Flags().String("versions.kratos", "v2.7.0", "Kratos version")
	initCmd.Flags().String("versions.wire", "v0.6.0", "Wire version")
	initCmd.Flags().String("versions.grpc", "1.62.1", "gRPC version")
	initCmd.Flags().String("versions.protobuf", "v1.33.0", "Protobuf version")
	initCmd.Flags().String("versions.automaxprocs", "1.5.1", "Automaxprocs version")
	initCmd.Flags().String("versions.genproto", "v2.3.4", "Genproto version")

	// Bind flags with viper
	viper.BindPFlag("monorepo-prefix", initCmd.Flags().Lookup("monorepo-prefix"))
	viper.BindPFlag("output-dir", initCmd.Flags().Lookup("output-dir"))
	viper.BindPFlag("versions.go", initCmd.Flags().Lookup("versions.go"))
	viper.BindPFlag("versions.kratos", initCmd.Flags().Lookup("versions.kratos"))
	viper.BindPFlag("versions.wire", initCmd.Flags().Lookup("versions.wire"))
	viper.BindPFlag("versions.grpc", initCmd.Flags().Lookup("versions.grpc"))
	viper.BindPFlag("versions.protobuf", initCmd.Flags().Lookup("versions.protobuf"))
	viper.BindPFlag("versions.automaxprocs", initCmd.Flags().Lookup("versions.automaxprocs"))
	viper.BindPFlag("versions.genproto", initCmd.Flags().Lookup("versions.genproto"))
}
