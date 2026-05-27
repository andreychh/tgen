// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir_test

import (
	"iter"
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/types"
)

type stubDesc struct{}

func (d stubDesc) Value() (string, error)   { return "", nil }
func (d stubDesc) Links() ([]string, error) { return nil, nil }

type stubField struct {
	expr    types.Expression
	typeErr error
}

func (f stubField) Key() (model.Key, error)                 { return "", nil }
func (f stubField) Type() (types.Expression, error)         { return f.expr, f.typeErr }
func (f stubField) Optionality() (model.Optionality, error) { return false, nil }
func (f stubField) Description() model.Description          { return stubDesc{} }

type stubObject struct {
	fields []spec.Field
}

func (o stubObject) Reference() (model.Reference, error) { return "", nil }
func (o stubObject) Name() (model.Name, error)           { return "", nil }
func (o stubObject) Description() model.Description      { return stubDesc{} }
func (o stubObject) Fields() iter.Seq[spec.Field]        { return slices.Values(o.fields) }

type stubDiscriminator struct{}

func (d stubDiscriminator) Key() (model.Key, error)                  { return "", nil }
func (d stubDiscriminator) Value() (model.DiscriminatorValue, error) { return "", nil }

type stubFields struct {
	free []spec.Field
}

func (f stubFields) Free() iter.Seq[spec.Field]        { return slices.Values(f.free) }
func (f stubFields) Discriminator() spec.Discriminator { return stubDiscriminator{} }

type stubDiscriminatedObject struct {
	fields []spec.Field
}

func (d stubDiscriminatedObject) Reference() (model.Reference, error) { return "", nil }
func (d stubDiscriminatedObject) Name() (model.Name, error)           { return "", nil }
func (d stubDiscriminatedObject) Description() model.Description      { return stubDesc{} }

func (d stubDiscriminatedObject) Fields() spec.Fields { return stubFields{free: d.fields} }

func toSpecFields(stubs []stubField) []spec.Field {
	out := make([]spec.Field, len(stubs))
	for i, s := range stubs {
		out[i] = s
	}
	return out
}
