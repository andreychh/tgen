// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

import (
	"cmp"
	"iter"
	"slices"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// ExplicitStructuredUnion represents a structured union from the Telegram Bot API spec for Go code generation.
type ExplicitStructuredUnion struct {
	inner explicit.StructuredUnion
}

// NewExplicitStructuredUnion constructs an ExplicitStructuredUnion from a parsed structured union.
func NewExplicitStructuredUnion(u explicit.StructuredUnion) ExplicitStructuredUnion {
	return ExplicitStructuredUnion{inner: u}
}

func (u ExplicitStructuredUnion) Name() Name {
	return NewName(u.inner.Name())
}

func (u ExplicitStructuredUnion) Doc() GoDoc {
	return NewGoDoc(NewDefinitionDoc(u.inner.Reference(), u.inner.Description()))
}

func (u ExplicitStructuredUnion) Variants() iter.Seq[Object] {
	return iters.NewMappedSeq(
		slices.Values(slices.SortedFunc(u.inner.Variants(), u.sortKey)),
		func(o explicit.Object) Object {
			return NewExplicitObject(o)
		},
	)
}

func (u ExplicitStructuredUnion) sortKey(a, b explicit.Object) int {
	return cmp.Compare(u.specificity(b), u.specificity(a))
}

// specificity returns a best-effort count of required fields in o. Fields whose
// optionality cannot be resolved are not counted.
func (u ExplicitStructuredUnion) specificity(o explicit.Object) int {
	var count int
	for f := range o.Fields() {
		opt, _ := f.Optionality().AsBool()
		if !opt {
			count++
		}
	}
	return count
}
