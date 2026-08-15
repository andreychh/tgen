// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types

import (
	"regexp"

	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/prose/grammar"
	"github.com/andreychh/tgen/model/typeexpr"
)

var (
	returnsSignal   = regexp.MustCompile(`(?i)\breturns\b`)
	returnedSignal  = regexp.MustCompile(`(?i)is returned`)
	arrayOfSignal   = regexp.MustCompile(`(?i)array of`)
	otherwiseSignal = regexp.MustCompile(`(?i)otherwise`)
)

// Rule is a pattern that recognizes one structural form of a return clause,
// decoding it into the type expression the clause names.
type Rule = grammar.Rule[typeexpr.Expression]

// ReturnsRule is a [Rule] that recognizes "Returns <type>", where the type is
// either a link to a documented type or an italic primitive word.
type ReturnsRule struct{}

// NewReturnsRule constructs a ReturnsRule.
func NewReturnsRule() ReturnsRule {
	return ReturnsRule{}
}

// Match implements [Rule].
func (ReturnsRule) Match(inlines []prose.Inline) (typeexpr.Expression, bool) {
	if len(inlines) < 2 {
		return nil, false
	}
	if !grammar.NewMarker(returnsSignal).Matches(inlines[0]) {
		return nil, false
	}
	if ref, ok := NewNamed().Matches(inlines[1]); ok {
		return typeexpr.NewNamed(ref), true
	}
	if kind, ok := NewPrimitive().Matches(inlines[1]); ok {
		return typeexpr.NewPrimitive(kind), true
	}
	return nil, false
}

// ReturnedRule is a [Rule] that recognizes "<type> is returned", where the type
// is either a link to a documented type or an italic primitive word.
type ReturnedRule struct{}

// NewReturnedRule constructs a ReturnedRule.
func NewReturnedRule() ReturnedRule {
	return ReturnedRule{}
}

// Match implements [Rule].
func (ReturnedRule) Match(inlines []prose.Inline) (typeexpr.Expression, bool) {
	if len(inlines) < 2 {
		return nil, false
	}
	if !grammar.NewMarker(returnedSignal).Matches(inlines[1]) {
		return nil, false
	}
	if ref, ok := NewNamed().Matches(inlines[0]); ok {
		return typeexpr.NewNamed(ref), true
	}
	if kind, ok := NewPrimitive().Matches(inlines[0]); ok {
		return typeexpr.NewPrimitive(kind), true
	}
	return nil, false
}

// ArrayRule is a [Rule] that recognizes "Array of <named>", a sequence whose
// element is a documented type.
type ArrayRule struct{}

// NewArrayRule constructs an ArrayRule.
func NewArrayRule() ArrayRule {
	return ArrayRule{}
}

// Match implements [Rule].
func (ArrayRule) Match(inlines []prose.Inline) (typeexpr.Expression, bool) {
	if len(inlines) < 2 {
		return nil, false
	}
	if !grammar.NewMarker(arrayOfSignal).Matches(inlines[0]) {
		return nil, false
	}
	ref, ok := NewNamed().Matches(inlines[1])
	if !ok {
		return nil, false
	}
	return typeexpr.NewArray(typeexpr.NewNamed(ref)), true
}

// UnionRule is a [Rule] that recognizes "<named> is returned, otherwise
// <primitive> is returned", a conditional return between a documented type and
// a primitive.
type UnionRule struct{}

// NewUnionRule constructs a UnionRule.
func NewUnionRule() UnionRule {
	return UnionRule{}
}

// Match implements [Rule].
func (UnionRule) Match(inlines []prose.Inline) (typeexpr.Expression, bool) {
	if len(inlines) < 4 {
		return nil, false
	}
	ref, ok := NewNamed().Matches(inlines[0])
	if !ok {
		return nil, false
	}
	if !grammar.NewMarker(otherwiseSignal).Matches(inlines[1]) {
		return nil, false
	}
	kind, found := NewPrimitive().Matches(inlines[2])
	if !found {
		return nil, false
	}
	if !grammar.NewMarker(returnedSignal).Matches(inlines[3]) {
		return nil, false
	}
	return typeexpr.NewUnion(typeexpr.NewNamed(ref), typeexpr.NewPrimitive(kind)), true
}
