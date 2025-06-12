// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package pkg

import (
	"github.com/spf13/cobra"
	"go.trai.ch/gecro/config"
)

// PkgCmd represents the base command for managing shared packages.
var PkgCmd = &cobra.Command{
	Use:     "pkg",
	Aliases: []string{"p"},
	Short:   "Manage shared packages (libraries)",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Load()
	},
}
