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
	definitions, err := r.definitions(spec.Definitions)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id definitions: %w", err)
	}
	aliases, err := pipeline.NewMergedTable(spec.Aliases, r.aliases()).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id aliases: %w", err)
	}
	variants, err := r.variants(spec.Variants)
	if err != nil {
		return Specification{}, fmt.Errorf("introducing chat id variants: %w", err)
	}
	fields, err := pipeline.NewMappedTable(spec.Fields, chatIDMapping{}).Apply()
	if err != nil {
		return Specification{}, fmt.Errorf("redirecting chat id fields: %w", err)
	}
	return Specification{
		Definitions:    definitions,
		Methods:        spec.Methods,
		Fields:         fields,
		Discriminators: spec.Discriminators,
		Variants:       variants,
		Aliases:        aliases,
		Release:        spec.Release,
	}, nil
}

// definitions returns base holding what ChatID names too: the union itself and
// the two primitive aliases its variants stand for, neither of which the
// documentation names. It fails when any of the three references is taken.
func (r ChatID) definitions(base Definitions) (Definitions, error) {
	out := NewDefinitionTable(base)
	err := out.Insert(
		chatIDRef,
		"ChatId",
		model.DefinitionKindUnion,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ChatId represents a chat identifier, either a numeric ID or a username.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the chat id union: %w", err)
	}
	err = out.Insert(
		idRef,
		"Id",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"ID represents a numeric Telegram chat or user identifier.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the id alias: %w", err)
	}
	err = out.Insert(
		usernameRef,
		"Username",
		model.DefinitionKindAlias,
		prose.NewPassage(prose.NewParagraph(prose.NewText(
			"Username represents a Telegram username.",
			prose.StylePlain,
		))),
	)
	if err != nil {
		return nil, fmt.Errorf("naming the username alias: %w", err)
	}
	return out, nil
}

// aliases returns the type each alias ChatID introduces stands for.
func (r ChatID) aliases() Aliases {
	out := pipeline.NewMapTableWithCapacity[model.Reference, Alias](2)
	out.Insert(idRef, Alias{Ref: idRef, Type: typeexpr.NewPrimitive(primitive.Integer)})
	out.Insert(usernameRef, Alias{Ref: usernameRef, Type: typeexpr.NewPrimitive(primitive.String)})
	return out
}

// variants returns base listing ChatID's two variants, id and username. It
// fails when the union already lists either.
func (r ChatID) variants(base parsed.Variants) (parsed.Variants, error) {
	out := NewVariantTable(base, chatIDRef)
	err := out.Insert(
		idRef,
		usernameRef,
	)
	if err != nil {
		return nil, fmt.Errorf("listing chat id variants: %w", err)
	}
	return out, nil
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
