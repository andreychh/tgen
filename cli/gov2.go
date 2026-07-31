// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/meta"
	ir "github.com/andreychh/tgen/model/ir/v2"
	"github.com/andreychh/tgen/output"
	"github.com/andreychh/tgen/source"
	"github.com/andreychh/tgen/targets/golangv2"
	"github.com/spf13/cobra"
)

// NewGoV2Command returns the "gov2" subcommand, which generates Go code from
// the nanopass pipeline rather than from the legacy chain "go" runs.
//
// TODO #248: Fold this command into "go" and delete the legacy chain.
func NewGoV2Command(m meta.Meta) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gov2",
		Short: "Generate Go client code from the nanopass pipeline (experimental)",
		RunE: func(cmd *cobra.Command, args []string) error {
			return goV2Action(cmd, args, m)
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

func goV2Action(cmd *cobra.Command, _ []string, m meta.Meta) error {
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
	spec, err := NewPipeline(doc).Specification()
	if err != nil {
		return fmt.Errorf("running the pipeline over %q: %w", location, err)
	}
	artifacts, err := golangv2.NewPass(
		golangv2.NewSpecification(ir.NewSpecification(spec), "api"),
	).Artifacts()
	if err != nil {
		return err
	}
	out := cmd.Flag("out").Value.String()
	err = output.NewFileset(artifacts).Emit(out)
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
