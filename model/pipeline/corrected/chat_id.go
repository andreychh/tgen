// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package corrected

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/pipeline/typed"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/prose"
	"github.com/andreychh/tgen/model/typeexpr"
)

const (
	chatIDRef   = model.Reference("chatid")
	idRef       = model.Reference("id")
	usernameRef = model.Reference("username")
)

// ChatID is a [Rule] that introduces the ChatID union — a chat addressed by
// numeric ID or by username — and redirects every Integer-or-String field to
// it.
type ChatID struct{}

// Apply implements [Rule]. It fails when chatid, id, or username already
// names something else in spec.
func (r ChatID) Apply(spec Specification) (Specification, error) {
	for _, ref := range []model.Reference{chatIDRef, idRef, usernameRef} {
		if alreadyExists(spec, ref) {
			return Specification{}, fmt.Errorf("%s already names something else in spec", ref)
		}
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id aliases: %w", err)
	}
	unions, err := pipeline.NewMergedTable(spec.Unions, r.union()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id union: %w", err)
	}
	variants, err := pipeline.NewMergedTable(spec.Variants, r.variants()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, chatIDMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting chat id fields: %w", err)
	}
	return Specification{
		Objects:        spec.Objects,
		Methods:        spec.Methods,
		Fields:         fields,
		Discriminators: spec.Discriminators,
		Unions:         unions,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// aliases returns the table of primitive aliases ChatID's variants need: a
// numeric id and a username, neither of which the documentation names.
func (r ChatID) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](2)
	out.Insert(idRef, Alias{
		Ref:  idRef,
		Name: "Id",
		Type: typeexpr.NewPrimitive(primitive.Integer),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ID represents a numeric Telegram chat or user identifier.",
			prose.StylePlain,
		))),
	})
	out.Insert(usernameRef, Alias{
		Ref:  usernameRef,
		Name: "Username",
		Type: typeexpr.NewPrimitive(primitive.String),
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"Username represents a Telegram username.",
			prose.StylePlain,
		))),
	})
	return out
}

// union returns the table holding the single ChatID union definition.
func (r ChatID) union() parsed.Unions {
	out := pipeline.NewMapTableWithCapacity[model.Reference, parsed.Union](1)
	out.Insert(chatIDRef, parsed.Union{
		Ref:  chatIDRef,
		Name: "ChatId",
		Description: prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ChatId represents a chat identifier, either a numeric ID or a username.",
			prose.StylePlain,
		))),
	})
	return out
}

// variants returns the table of ChatID's two variants, id and username.
func (r ChatID) variants() parsed.Variants {
	out := pipeline.NewMapTableWithCapacity[model.VariantKey, parsed.Variant](2)
	out.Insert(
		model.VariantKey{Owner: chatIDRef, Ref: idRef},
		parsed.Variant{Ref: idRef},
	)
	out.Insert(
		model.VariantKey{Owner: chatIDRef, Ref: usernameRef},
		parsed.Variant{Ref: usernameRef},
	)
	return out
}

// chatIDMapping is a [pipeline.Mapping] that redirects a field typed as
// Integer or String to the ChatID union; every other field rides through
// unchanged.
type chatIDMapping struct{}

// Apply implements [pipeline.Mapping]. It never fails.
func (chatIDMapping) Apply(field typed.Field) (typed.Field, error) {
	if !field.Type.Equals(typeexpr.NewUnion(
		typeexpr.NewPrimitive(primitive.Integer),
		typeexpr.NewPrimitive(primitive.String),
	)) {
		return field, nil
	}
	return typed.Field{
		Key:         field.Key,
		Position:    field.Position,
		Type:        typeexpr.NewNamed(chatIDRef),
		Optionality: field.Optionality,
		Description: field.Description,
	}, nil
}
