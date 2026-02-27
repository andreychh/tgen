// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cmd

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/andreychh/tgen/rendering"
	"github.com/andreychh/tgen/rendering/golang"
	"github.com/andreychh/tgen/source"
	"github.com/spf13/cobra"
)

// NewGo returns the "go" subcommand.
//
// TODO #43: Change `spec` default value to Telegram Bot API URL.
// TODO #43: Add an option to specify the Go package name.
func NewGo() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go",
		Short: "Generate Go structures and methods",
		RunE:  goAction,
	}
	cmd.Flags().StringP(
		"spec",
		"s",
		".notes/api/api.html",
		"Path to the local api.html file",
	)
	cmd.Flags().StringP(
		"out",
		"o",
		"api",
		"Output directory for generated files",
	)
	return cmd
}

func goAction(cmd *cobra.Command, _ []string) error {
	location := filepath.Clean(cmd.Flag("spec").Value.String())
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	reader, err := source.NewLocationSource(location).Open(ctx)
	if err != nil {
		return fmt.Errorf("opening source %q: %w", location, err)
	}
	defer func() { _ = reader.Close() }()
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return fmt.Errorf("parsing HTML from %q: %w", location, err)
	}
	sel := dom.NewHTMLSelection(doc.Selection)
	spec := parsing.NewRawSpecification(sel)
	tmpl := golang.PrepareTemplate()
	fileset := rendering.NewFileset(rendering.Artifacts{
		"unions.go":  golang.NewUnionsView(tmpl, spec),
		"objects.go": golang.NewObjectsView(tmpl, spec),
	})
	out := filepath.Clean(cmd.Flag("out").Value.String())
	err = fileset.Emit(out)
	if err != nil {
		return fmt.Errorf("generating files in directory %q: %w", out, err)
	}
	return nil
}
