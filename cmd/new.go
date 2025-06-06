// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package cmd

import (
	"errors"
	"fmt"
	"os"

	"go.trai.ch/msgen/internal/generator"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new [servicename]",
	Short: "Creates a new microservice",
	Long: `Creates a new microservice with a specified name
	Example: msgen new test`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		const moduleFile = "MODULE.bazel"

		if _, err := os.Stat(moduleFile); errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("This command must be run from the monorepo root (missing %s)", moduleFile)
		}

		serviceName := args[0]

		serviceParams := generator.ServiceParams{
			ServiceName:       serviceName,
			MonorepoPrefix:    "go.trai.ch/trai/",
			OutputDir:         ".",
			KratosVersion:     "2.8.4",
			WireVersion:       "0.6.0",
			GoVersion:         "1.24.3",
			DefaultConfigPath: "proto/config/common/",
		}

		generator, err := generator.NewGenerator()
		if err != nil {
			return fmt.Errorf("failed creating new generator: %w", err)
		}

		if err := generator.GenerateService(serviceParams); err != nil {
			return fmt.Errorf("failed generating new service: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("dry-run", "d", true, "Dry run")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
