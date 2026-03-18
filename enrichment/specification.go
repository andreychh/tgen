// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

// Package enrichment supplements the parsed Telegram Bot API specification
// with tgen's editorial additions: canonical names for inline union types and
// rule-based field and method type overrides.
package enrichment

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/parsing"
)

// Specification represents an enriched Telegram Bot API specification.
type Specification struct {
	inner parsing.Specification
}

// NewSpecification creates a Specification from a parsed specification.
func NewSpecification(s parsing.Specification) Specification {
	return Specification{inner: s}
}

func (s Specification) Objects() iter.Seq[Object] {
	return func(yield func(Object) bool) {
		for o := range s.inner.Objects() {
			if !yield(NewObject(o)) {
				break
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[Method] {
	return func(yield func(Method) bool) {
		for m := range s.inner.Methods() {
			if !yield(NewMethod(m)) {
				break
			}
		}
	}
}

func (s Specification) Unions() iter.Seq[parsing.Union] {
	return s.inner.Unions()
}

func (s Specification) Release() parsing.Release {
	return s.inner.Release()
}

func (s Specification) ImplicitUnions() iter.Seq[ImplicitUnion] {
	return slices.Values([]ImplicitUnion{
		NewImplicitUnion(
			"ChatId",
			"ChatId represents a chat identifier, which is either an integer or a string.",
			[]string{"Integer", "String"},
		),
		NewImplicitUnion(
			"ReplyMarkup",
			"ReplyMarkup represents a reply markup attached to a message.",
			[]string{
				"InlineKeyboardMarkup",
				"ReplyKeyboardMarkup",
				"ReplyKeyboardRemove",
				"ForceReply",
			},
		),
		NewImplicitUnion(
			"InputMediaGroup",
			"InputMediaGroup represents a media element in a media group.",
			[]string{"InputMediaAudio", "InputMediaDocument", "InputMediaPhoto", "InputMediaVideo"},
		),
		NewImplicitUnion(
			"MaybeMessage",
			"MaybeMessage represents a method return type that is either an edited Message or True for inline messages.",
			[]string{"Message", "True"},
		),
	})
}
