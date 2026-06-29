// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package unified

import (
	"errors"
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/parsed"
	"github.com/andreychh/tgen/model/prose"
)

// Field is a field of an object or a parameter of a method, reduced to a single
// shape: its key, the verbatim prose of its type, whether it is optional, and
// its description with the optional marker removed.
type Field struct {
	Key         model.Key
	Type        prose.Phrase
	Optionality model.Optionality
	Description prose.Phrase
}

// FieldMapping maps a parsed object field into a unified field. An object field
// carries its optionality as an italic "Optional" prefix in its description,
// which the mapping lifts out and strips along with the ". " separator that
// follows it.
type FieldMapping struct{}

// NewFieldMapping constructs a FieldMapping.
func NewFieldMapping() FieldMapping {
	return FieldMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the description carries
// the "Optional" marker without its following ". " separator.
func (m FieldMapping) Apply(field parsed.Field) (Field, error) {
	description, optional, err := m.resolveDescription(field.Description)
	if err != nil {
		return Field{}, fmt.Errorf("resolving description: %w", err)
	}
	return Field{
		Key:         field.Key,
		Type:        field.Type,
		Optionality: model.Optionality(optional),
		Description: description,
	}, nil
}

// resolveDescription splits a field's description into its optionality and
// normalized prose. A description that does not open with the italic "Optional"
// marker belongs to a required field and rides through unchanged. One that does
// is optional: the marker and the ". " separator that follows it are removed.
// It fails when the marker is present but the ". " separator is not.
func (m FieldMapping) resolveDescription(
	description prose.Phrase,
) (prose.Phrase, bool, error) {
	inlines := description.Inlines()
	marker := prose.NewText("Optional", prose.StyleItalic)
	if len(inlines) == 0 || !inlines[0].Equals(marker) {
		return description, false, nil
	}
	if len(inlines) < 2 {
		return prose.Phrase{}, false, errors.New(
			`optional field description lacks a ". " separator`,
		)
	}
	text, ok := inlines[1].(prose.Text)
	if !ok || text.Style() != prose.StylePlain || !text.HasPrefix(". ") {
		return prose.Phrase{}, false, errors.New(
			`optional field description lacks a ". " separator`,
		)
	}
	rest := append([]prose.Inline{text.TrimPrefix(". ")}, inlines[2:]...)
	return prose.NewPhrase(rest...), true, nil
}

// ParamMapping maps a parsed method parameter into a unified field. A parameter
// carries its optionality in a dedicated Required column rather than a
// description prefix, so its description needs no normalization.
type ParamMapping struct{}

// NewParamMapping constructs a ParamMapping.
func NewParamMapping() ParamMapping {
	return ParamMapping{}
}

// Apply implements [pipeline.Mapping]. It fails when the Required column is
// neither "Yes" nor "Optional".
func (m ParamMapping) Apply(param parsed.Param) (Field, error) {
	optional, err := m.resolveRequired(param.Required)
	if err != nil {
		return Field{}, fmt.Errorf("resolving required column: %w", err)
	}
	return Field{
		Key:         param.Key,
		Type:        param.Type,
		Optionality: model.Optionality(optional),
		Description: param.Description,
	}, nil
}

// resolveRequired reads a parameter's Required column, yielding true when it
// holds "Optional" and false when it holds "Yes". It fails on any other
// content.
func (m ParamMapping) resolveRequired(required prose.Phrase) (bool, error) {
	inlines := required.Inlines()
	optional := prose.NewText("Optional", prose.StylePlain)
	yes := prose.NewText("Yes", prose.StylePlain)
	if len(inlines) != 1 {
		return false, errors.New("unexpected required column")
	}
	if inlines[0].Equals(optional) {
		return true, nil
	}
	if inlines[0].Equals(yes) {
		return false, nil
	}
	return false, errors.New("unexpected required column")
}
