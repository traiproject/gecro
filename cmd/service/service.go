// Copyright (c) TRAI
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/spf13/cobra"
)

var ServiceCmd = &cobra.Command{
	Use:     "service",
	Aliases: []string{"svc", "s"},
	Short:   "Manage services",
}
