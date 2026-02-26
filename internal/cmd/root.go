// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cmd

import "github.com/urfave/cli/v3"

// NewRoot returns the primary application command ("tgen").
func NewRoot() *cli.Command {
	return &cli.Command{
		Name:    "tgen",
		Usage:   "Telegram Bot API code generator",
		Version: "",
		// TODO #39: Add an explicit 'completion [shell]' command to make the built-in (but
		// 	hidden) shell completion generator visible to users in the 'tgen --help'
		// 	output.
		EnableShellCompletion: true,
		Commands: []*cli.Command{
			NewGo(),
		},
	}
}
