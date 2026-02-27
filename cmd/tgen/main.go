// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package main is the entry point for the tgen CLI application.
package main

import (
	"os"

	"github.com/andreychh/tgen/internal/cmd"
)

func main() {
	err := cmd.NewRoot().Execute()
	if err != nil {
		os.Exit(1)
	}
}
