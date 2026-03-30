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

// ReturnType extracts the return TypeExpression of a method from its
// description paragraphs.
type ReturnType struct {
	h4 gq.Selection
}

// NewReturnType creates a ReturnType from an h4 selection.
func NewReturnType(h4 gq.Selection) ReturnType {
	return ReturnType{h4: h4}
}

// AsExpression parses the method description and returns the return type
// expression. Returns an error if no description paragraphs are found or the
// return type cannot be extracted.
func (t ReturnType) AsExpression() (types.TypeExpression, error) {
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
	return extractReturnType(strings.Join(parts, " "))
}

func extractReturnType(text string) (types.TypeExpression, error) {
	if m := returnConditional.FindStringSubmatch(text); m != nil {
		return types.NewUnionType(
			[]types.TypeExpression{types.NewNamedType(m[1]), types.NewNamedType(m[2])},
		), nil
	}
	if m := returnArray.FindStringSubmatch(text); m != nil {
		return types.NewArrayType(types.NewNamedType(m[1])), nil
	}
	if m := returnAsType.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnInFormOf.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnArticleObject.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnDirect.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnTheNamed.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnThePre.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	if m := returnFallback.FindStringSubmatch(text); m != nil {
		return types.NewNamedType(m[1]), nil
	}
	return nil, fmt.Errorf("cannot extract return type from: %q", text)
}
