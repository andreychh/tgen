// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package enrichment

import (
	"iter"

	"github.com/andreychh/tgen/parsing"
)

// Object represents an enriched Telegram Bot API object definition.
type Object struct {
	inner parsing.Object
	rules []FieldRule
}

// NewObject constructs an Object from a parsed object with all field rules applied.
func NewObject(o parsing.Object) Object {
	return Object{
		inner: o,
		rules: []FieldRule{ChatIdRule{}, ReplyMarkupRule{}, InputMediaGroupRule{}, InputFileRule{}},
	}
}

func (o Object) Ref() parsing.DefinitionRef {
	return o.inner.Ref()
}

//nolint:ireturn // ObjectName is the intentional public contract of this method
func (o Object) Name() parsing.ObjectName {
	return o.inner.Name()
}

//nolint:ireturn // DefinitionDescription is the intentional public contract of this method
func (o Object) Description() parsing.DefinitionDescription {
	return o.inner.Description()
}

func (o Object) Fields() iter.Seq[parsing.Field] {
	return func(yield func(parsing.Field) bool) {
		for f := range o.inner.Fields() {
			if !yield(o.applyRules(f)) {
				break
			}
		}
	}
}

func (o Object) applyRules(f parsing.Field) parsing.Field {
	for _, r := range o.rules {
		f = r.Apply(f)
	}
	return f
}
