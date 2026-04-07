// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/model/types"
	"github.com/andreychh/tgen/pkg/gq"
)

var (
	returnConditional = regexp.MustCompile(
		`the (?:\w+ )?([A-Z][a-zA-Z]+) is returned, otherwise ([A-Z][a-zA-Z]+) is returned`,
	)
	returnArray  = regexp.MustCompile(`(?i)array of ([A-Z][a-zA-Z]+)`)
	returnAsType = regexp.MustCompile(
		`[Rr]eturns?.{0,100}as (?:a )?([A-Z][a-zA-Z]+)(?: object)?`,
	)
	returnInFormOf      = regexp.MustCompile(`in form of a ([A-Z][a-zA-Z]+) object`)
	returnArticleObject = regexp.MustCompile(`\ba ([A-Z][a-zA-Z]+) object`)
	returnDirect        = regexp.MustCompile(`Returns ([A-Z][a-zA-Z]+) on success`)
	returnTheNamed      = regexp.MustCompile(`Returns the (?:\w+ )?([A-Z][a-zA-Z]+)`)
	returnThePre        = regexp.MustCompile(`the (?:\w+ )?([A-Z][a-zA-Z]+) is returned`)
	returnFallback      = regexp.MustCompile(`([A-Z][a-zA-Z]+) is returned`)
)

// ReturnType extracts the return Expression of a method from its
// description paragraphs.
type ReturnType struct {
	root, h4 gq.Selection
}

// NewReturnType creates a ReturnType from an h4 selection.
func NewReturnType(root, h4 gq.Selection) ReturnType {
	return ReturnType{root: root, h4: h4}
}

// AsExpression parses the method description and returns the return type
// expression. Returns an error if no description paragraphs are found or the
// return type cannot be extracted.
func (t ReturnType) AsExpression() (types.Expression, error) {
	var parts []string
	for node := range t.h4.Until("h3, h4, hr").Filter("p").All() {
		text := node.Text()
		if text != "" {
			parts = append(parts, text)
		}
	}
	if len(parts) == 0 {
		return nil, errors.New("no description paragraphs found")
	}
	return t.extractReturnType(strings.Join(parts, " "))
}

func (t ReturnType) extractReturnType(text string) (types.Expression, error) {
	if m := returnConditional.FindStringSubmatch(text); m != nil {
		first, err := t.named(m[1])
		if err != nil {
			return nil, err
		}
		second, err := t.named(m[2])
		if err != nil {
			return nil, err
		}
		return types.NewUnion(first, second), nil
	}
	if m := returnArray.FindStringSubmatch(text); m != nil {
		elem, err := t.named(m[1])
		if err != nil {
			return nil, err
		}
		return types.NewArray(elem), nil
	}
	if m := returnAsType.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnInFormOf.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnArticleObject.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnDirect.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnTheNamed.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnThePre.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	if m := returnFallback.FindStringSubmatch(text); m != nil {
		return t.named(m[1])
	}
	return nil, fmt.Errorf("cannot extract return type from: %q", text)
}

func (t ReturnType) named(name string) (types.Expression, error) {
	kind, found := NewCatalog(t.root).Lookup(name)
	if !found {
		return nil, fmt.Errorf("unknown type %q", name)
	}
	return types.NewNamed(name, kind), nil
}
