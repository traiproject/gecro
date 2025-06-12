// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/spf13/cobra"
	"go.trai.ch/gecro/config"
)

var ServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc", "s"},
	Short:   "Manage services",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return config.Load()
	},
}
