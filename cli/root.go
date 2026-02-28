// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import "github.com/spf13/cobra"

// NewRootCommand returns the primary application command ("tgen").
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tgen",
		Short: "Telegram Bot API code generator",
	}
	cmd.AddCommand(NewGoCommand())
	return cmd
}
