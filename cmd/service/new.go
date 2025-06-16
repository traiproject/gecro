// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package service

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.trai.ch/gecro/config"
	"go.trai.ch/gecro/internal/generator"

	"github.com/spf13/cobra"
)

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:     "new [servicename]",
	Aliases: []string{"n"},
	Short:   "Creates a new microservice",
	Long: `Creates a new microservice with a specified name
	Example: gecro new test`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		config := config.Cfg
		config.Name = args[0]

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		config.DryRun = dryRun

		generator, err := generator.NewGenerator()
		if err != nil {
			return fmt.Errorf("failed creating new generator: %w", err)
		}

		if err := generator.GenerateService(*config); err != nil {
			return fmt.Errorf("failed generating new service: %w", err)
		}

		if err := runPostCreateHook(config.DryRun); err != nil {
			return err
		}

		return nil
	},
}

func runPostCreateHook(dryRun bool) error {
	if dryRun {
		return nil
	}

	command := "bazel"
	args := []string{"run", "//:gazelle"}

	fmt.Printf("Running post-create hook: %s %s\n", command, strings.Join(args, " "))
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func init() {
	ServiceCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("dry-run", "d", false, "Dry run")
}
