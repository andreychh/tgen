// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"github.com/andreychh/tgen/meta"
	"github.com/spf13/cobra"
)

// NewRootCommand returns the primary application command ("tgen").
func NewRootCommand() *cobra.Command {
	metadata := meta.NewMeta(meta.NewDetectedSource())
	cmd := &cobra.Command{
		Use:     "tgen",
		Short:   "Generate strongly-typed Telegram Bot API clients",
		Version: metadata.Release().Version(),
	}
	cmd.SetVersionTemplate(NewVersionMessage(metadata).String())
	cmd.AddCommand(NewGoCommand(metadata))
	cmd.AddCommand(NewPythonCommand(metadata))
	return cmd
}
