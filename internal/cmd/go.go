// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/andreychh/tgen/rendering"
	"github.com/andreychh/tgen/rendering/golang"
	"github.com/urfave/cli/v3"
)

// NewGo returns the "go" subcommand.
func NewGo() *cli.Command {
	return &cli.Command{
		Name:  "go",
		Usage: "Generate Go structures and methods",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "spec",
				Aliases: []string{"s"},
				Usage:   "Path to the local api.html file",
				// TODO: use "https://core.telegram.org/bots/api" as the default value.
				Value: ".notes/api/api.html",
			},
			&cli.StringFlag{
				Name:    "out",
				Aliases: []string{"o"},
				Usage:   "Output directory for generated files",
				Value:   "api",
			},
			// TODO: add an option to specify the Go package name.
			// &cli.StringFlag{
			//     Name:    "package",
			//     Aliases: []string{"p"},
			//     Usage:   "Go package name for generated code",
			//     Value:   "api",
			// },
		},
		Action: goAction,
	}
}

func goAction(_ context.Context, cmd *cli.Command) error {
	path := filepath.Clean(cmd.String("spec"))
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open HTML file (%s): %w", path, err)
	}
	defer func() { _ = file.Close() }()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return fmt.Errorf("failed to parse HTML file: %w", err)
	}
	sel := dom.NewHTMLSelection(doc.Selection)
	spec := parsing.NewRawSpecification(sel)
	tmpl := golang.PrepareTemplate()
	fileset := rendering.NewFileset(rendering.Artifacts{
		"unions.go":  golang.NewUnionsView(tmpl, spec),
		"objects.go": golang.NewObjectsView(tmpl, spec),
	})
	out := cmd.String("out")
	err = fileset.Emit(out)
	if err != nil {
		return fmt.Errorf("failed to emit files to directory '%s': %w", out, err)
	}
	return nil
}
