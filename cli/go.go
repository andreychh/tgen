// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/andreychh/tgen/rendering"
	"github.com/andreychh/tgen/rendering/golang"
	"github.com/andreychh/tgen/source"
	"github.com/spf13/cobra"
)

// NewGoCommand returns the "go" subcommand.
//
// TODO #43: Add an option to specify the Go package name.
func NewGoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go",
		Short: "Generate Go client code",
		RunE:  goAction,
	}
	cmd.Flags().StringP(
		"spec",
		"s",
		"https://core.telegram.org/bots/api",
		"URL or local path to the Telegram Bot API HTML specification",
	)
	cmd.Flags().StringP(
		"out",
		"o",
		"./api",
		"Output directory for the generated Go files",
	)
	return cmd
}

func goAction(cmd *cobra.Command, _ []string) error {
	location := cmd.Flag("spec").Value.String()
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
	out := cmd.Flag("out").Value.String()
	err = fileset.Emit(out)
	if err != nil {
		return fmt.Errorf("generating files in directory %q: %w", out, err)
	}
	return nil
}
