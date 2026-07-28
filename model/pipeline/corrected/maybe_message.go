// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/resolved"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/result"
	typetree "github.com/andreychh/tgen/model/types/v2"
)

const (
	maybeMessageRef = model.Reference("maybemessage")
	messageRef      = model.Reference("message")
	trueRef         = model.Reference("true")
)

// MaybeMessage is a [Rule] that introduces the MaybeMessage union — a method
// return value that is either an edited Message or True for inline
// messages — and redirects every method returning their union to it.
type MaybeMessage struct{}

// Apply implements [Rule]. It fails when maybemessage or true already names
// something else in spec.
func (r MaybeMessage) Apply(spec Specification) (Specification, error) {
	for _, ref := range []model.Reference{maybeMessageRef, trueRef} {
		if alreadyExists(spec, ref) {
			return Specification{}, fmt.Errorf("%s already names something else in spec", ref)
		}
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message aliases: %w", err)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing maybe message variants: %w", err)
	}
	methods, err := pipeline.NewMappedTable(spec.Methods, maybeMessageMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting maybe message methods: %w", err)
	}
	return Specification{
		Objects:        spec.Objects,
		Methods:        methods,
		Fields:         spec.Fields,
		Discriminators: spec.Discriminators,
		Unions:         unions,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// aliases returns the table of primitive aliases MaybeMessage's variants
// need: True, which the documentation writes as a bare keyword, not a named
// type.
func (r MaybeMessage) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](1)
	out.Insert(trueRef, Alias{
		Ref:  trueRef,
		Name: "True",
		Type: typetree.NewPrimitive(typetree.True),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"True represents the boolean true value in Telegram API responses.",
			prose.StylePlain,
		))),
	})
	return out
}

// union returns the table holding the single MaybeMessage union definition.
func (r MaybeMessage) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(maybeMessageRef, parsed.Union{
		Ref:  maybeMessageRef,
		Name: "MaybeMessage",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"MaybeMessage represents a method return value that is either an "+
				"edited Message or True for inline messages.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of MaybeMessage's two variants, message and
// true.
func (r MaybeMessage) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](2)
	out.Insert(
		model.VariantKey{Owner: maybeMessageRef, Ref: messageRef},
		parsed.Variant{Ref: messageRef},
	)
	out.Insert(
		model.VariantKey{Owner: maybeMessageRef, Ref: trueRef},
		parsed.Variant{Ref: trueRef},
	)
	return out
}

// maybeMessageMapping is a [pipeline.Mapping] that redirects a method
// returning the Message-or-True union to MaybeMessage; every other method
// rides through unchanged.
type maybeMessageMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (maybeMessageMapping) Apply(method resolved.Method) (resolved.Method, error) {
	value, ok := method.Result.(result.Value)
	if !ok {
		return method, nil
	}
	if !value.Type().Equals(typetree.NewUnion(
		typetree.NewNamed(messageRef),
		typetree.NewPrimitive(typetree.True),
	)) {
		return method, nil
	}
	return resolved.Method{
		Ref:         method.Ref,
		Name:        method.Name,
		Description: method.Description,
		Result:      result.NewValue(typetree.NewNamed(maybeMessageRef)),
	}, nil
}
