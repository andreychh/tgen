// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package main is the entry point for the tgen CLI application.
package main

import (
	"context"
	"os"

	"github.com/andreychh/tgen/internal/cmd"
)

func main() {
	err := cmd.NewRoot().Run(context.Background(), os.Args)
	if err != nil {
		_, _ = os.Stderr.WriteString(err.Error() + "\n")
		os.Exit(1)
	}
}
