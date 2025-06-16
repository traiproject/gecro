// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"go.trai.ch/gecro/config"
	"go.trai.ch/gecro/internal/generator"

	"github.com/spf13/cobra"
)

// newCmd represents the new command for packages
var newCmd = &cobra.Command{
	Use:     "new [packagename]",
	Aliases: []string{"n"},
	Short:   "Creates a new shared package (library)",
	Long: `Creates a new shared package (library) in the /pkg directory.
    Example: gecro pkg new common-utils`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Note: We are reusing the ServiceName field from the config struct
		// for the package name to avoid changing the config logic.
		config := config.Cfg
		config.Name = args[0]

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}
		config.DryRun = dryRun

		gen, err := generator.NewGenerator()
		if err != nil {
			return fmt.Errorf("failed creating new generator: %w", err)
		}

		// Call the new generator function for packages
		if err := gen.GeneratePkg(*config); err != nil {
			return fmt.Errorf("failed generating new package: %w", err)
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
	PkgCmd.AddCommand(newCmd)

	newCmd.Flags().BoolP("dry-run", "d", false, "Dry run without writing files")
}
