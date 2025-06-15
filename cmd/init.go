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

	// Add flags for all config options
	initCmd.Flags().String("monorepo-prefix", viper.GetString("monorepo-prefix"), "Monorepo prefix")
	initCmd.Flags().String("output-dir", viper.GetString("output-dir"), "Output directory for generated files")
	initCmd.Flags().String("versions.go", viper.GetString("versions.go"), "Go version")
	initCmd.Flags().String("versions.kratos", viper.GetString("versions.kratos"), "Kratos version")
	initCmd.Flags().String("versions.wire", viper.GetString("versions.wire"), "Wire version")
	initCmd.Flags().String("versions.grpc", viper.GetString("versions.grpc"), "gRPC version")
	initCmd.Flags().String("versions.protobuf", viper.GetString("versions.protobuf"), "Protobuf version")
	initCmd.Flags().String("versions.automaxprocs", viper.GetString("versions.automaxprocs"), "Automaxprocs version")
	initCmd.Flags().String("versions.genproto", viper.GetString("versions.genproto"), "Genproto version")

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
