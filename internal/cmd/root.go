// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cmd

import "github.com/spf13/cobra"

// NewRoot returns the primary application command ("tgen").
func NewRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tgen",
		Short: "Telegram Bot API code generator",
	}
	cmd.AddCommand(NewGo())
	return cmd
}
