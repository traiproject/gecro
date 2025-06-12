// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.trai.ch/gecro/cmd/pkg"
	"go.trai.ch/gecro/cmd/service"
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
	rootCmd.AddCommand(service.ServiceCmd)
	rootCmd.AddCommand(pkg.PkgCmd)
}
