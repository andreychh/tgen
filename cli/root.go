// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import "github.com/spf13/cobra"

// NewRootCommand returns the primary application command ("tgen").
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tgen",
		Short: "Generate strongly-typed Telegram Bot API clients",
		Long: `tgen turns the Telegram Bot API HTML documentation into ready-to-use API bindings.

Instead of relying on manually updated boilerplate, tgen parses the specification to generate strongly-typed client code.`,
	}
	cmd.AddCommand(NewGoCommand())
	return cmd
}
