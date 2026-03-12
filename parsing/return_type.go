// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsing

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/andreychh/tgen/parsing/gq"
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

type ReturnType struct {
	selection gq.Selection
}

func NewReturnType(h4 gq.Selection) ReturnType {
	return ReturnType{selection: h4}
}

func (r ReturnType) Value() (string, error) {
	var parts []string
	for node := range r.selection.Until("h3, h4, hr").Filter("p").All() {
		text := node.Text()
		if text != "" {
			parts = append(parts, text)
		}
	}
	if len(parts) == 0 {
		return "", errors.New("no description paragraphs found")
	}
	return extractReturnType(strings.Join(parts, " "))
}

func extractReturnType(text string) (string, error) {
	if m := returnConditional.FindStringSubmatch(text); m != nil {
		return m[1] + " or " + m[2], nil
	}
	if m := returnArray.FindStringSubmatch(text); m != nil {
		return "Array of " + m[1], nil
	}
	if m := returnAsType.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnInFormOf.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnArticleObject.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnDirect.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnTheNamed.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnThePre.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	if m := returnFallback.FindStringSubmatch(text); m != nil {
		return m[1], nil
	}
	return "", fmt.Errorf("cannot extract return type from: %q", text)
}
