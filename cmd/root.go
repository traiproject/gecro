// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.trai.ch/gecro/config"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gecro",
	Short: "A generative CLI for Go-Kratos microservices with Bazel",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(func() {
		if err := config.Load(); err != nil {
			fmt.Fprintf(os.Stderr, "Config error: %v", err)
			os.Exit(1)
		}
	})
}
