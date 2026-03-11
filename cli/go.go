// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/meta"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/gq"
	"github.com/andreychh/tgen/rendering"
	"github.com/andreychh/tgen/rendering/golang"
	"github.com/andreychh/tgen/source"
	"github.com/spf13/cobra"
)

// NewGoCommand returns the "go" subcommand.
//
// TODO #43: Add an option to specify the Go package name.
func NewGoCommand(m meta.Meta) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "go",
		Short: "Generate Go client code",
		RunE: func(cmd *cobra.Command, args []string) error {
			return goAction(cmd, args, m)
		},
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

func goAction(cmd *cobra.Command, _ []string, m meta.Meta) error {
	snapshot := meta.NewSnapshot(m)
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
	artifacts, err := golang.NewArtifacts(
		parsing.NewSpecification(
			gq.NewNormSelection(doc.Selection),
		),
	).Value()
	if err != nil {
		return err
	}
	out := cmd.Flag("out").Value.String()
	err = rendering.NewFileset(artifacts).Emit(out)
	if err != nil {
		return fmt.Errorf("generating files in directory %q: %w", out, err)
	}
	_, err = fmt.Fprintf(
		cmd.ErrOrStderr(),
		"done in %s\n",
		snapshot.Elapsed().Round(time.Millisecond),
	)
	return err
}
