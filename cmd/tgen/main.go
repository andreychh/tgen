// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// This is the main entry point for the tgen tool.
package main

import (
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/andreychh/tgen/rendering"
	"github.com/andreychh/tgen/rendering/golang"
)

func main() {
	file, err := os.Open(".notes/api/api.html")
	if err != nil {
		panic(fmt.Sprintf("failed to open HTML file: %v", err))
	}
	defer func() { _ = file.Close() }()
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		panic(fmt.Sprintf("failed to parse HTML file: %v", err))
	}
	sel := dom.NewHTMLSelection(doc.Selection)
	spec := parsing.NewRawSpecification(sel)
	tmpl := golang.PrepareTemplate()
	fileset := rendering.NewFileset(rendering.Artifacts{
		"views.go": golang.NewUnionsView(tmpl, spec),
	})
	err = fileset.Emit("api")
	if err != nil {
		panic(fmt.Sprintf("failed to emit files: %v", err))
	}
}
