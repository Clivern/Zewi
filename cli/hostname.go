// Copyright 2025 Clivern. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"

	"github.com/clivern/zewi/core"

	"github.com/spf13/cobra"
)

var hostnameCmd = &cobra.Command{
	Use:   "hostname",
	Short: "Start a server that returns hostname on the home page",
	Run: func(_ *cobra.Command, _ []string) {
		// Load configuration
		if err := core.Load(config); err != nil {
			panic(err.Error())
		}

		// Setup logging
		if err := core.SetupLogging(); err != nil {
			panic(err.Error())
		}

		// Setup and configure the HTTP server
		r := core.SetupHostname(Static)

		// Run the server
		if err := core.RunHostname(r); err != nil {
			panic(fmt.Sprintf("Server error: %s", err.Error()))
		}
	},
}

func init() {
	hostnameCmd.Flags().StringVarP(
		&config,
		"config",
		"c",
		"config.prod.yml",
		"Absolute path to config file (required)",
	)
	hostnameCmd.MarkFlagRequired("config")
	rootCmd.AddCommand(hostnameCmd)
}
