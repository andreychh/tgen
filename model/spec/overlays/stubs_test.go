// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays_test

import (
	"errors"
	"iter"
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/model/types"
)

type stubDesc struct {
	links    []string
	linksErr error
}

func (d stubDesc) Value() (string, error) { return "", nil }
func (d stubDesc) Links() ([]string, error) {
	return d.links, d.linksErr
}

var errLinks = errors.New("links unavailable")

type stubField struct {
	key      model.Key
	expr     types.Expression
	optional model.Optionality
	desc     model.Description
	typeErr  error
	optErr   error
}

func (f stubField) Key() (model.Key, error)                 { return f.key, nil }
func (f stubField) Type() (types.Expression, error)         { return f.expr, f.typeErr }
func (f stubField) Optionality() (model.Optionality, error) { return f.optional, f.optErr }
func (f stubField) Description() model.Description          { return f.desc }

var (
	errType = errors.New("type unavailable")
	errOpt  = errors.New("optionality unavailable")
)

type stubMethod struct {
	result    spec.Result
	resultErr error
	fields    []spec.Field
}

func (m stubMethod) Reference() (model.Reference, error) { return "", nil }
func (m stubMethod) Name() (model.Name, error)           { return "", nil }
func (m stubMethod) Description() model.Description      { return stubDesc{} }
func (m stubMethod) Result() (spec.Result, error)        { return m.result, m.resultErr }
func (m stubMethod) Fields() iter.Seq[spec.Field]        { return slices.Values(m.fields) }
