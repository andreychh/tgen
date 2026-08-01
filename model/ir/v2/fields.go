// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/flattened"
	"github.com/andreychh/tgen/model/pipeline/separated"
	"github.com/andreychh/tgen/model/typeform"
)

// Fields is the ordered view of the fields one definition owns, each with its
// type bound to the definition it names. It gives the required fields before
// the optional ones, and within either group the order their source gave them:
// a caller fills in every required field, so those are what a reader of the
// generated type meets first.
type Fields struct {
	db    separated.Specification
	owner model.Reference
}

// NewFields constructs a Fields over the fields owner holds in the database.
func NewFields(db separated.Specification, owner model.Reference) Fields {
	return Fields{db: db, owner: owner}
}

// Value returns the fields of the owner, the required ones before the optional.
// It fails when a field's type names no definition.
func (f Fields) Value() ([]Field, error) {
	records := f.records()
	out := make([]Field, 0, len(records))
	for _, record := range records {
		field, err := f.field(record)
		if err != nil {
			return nil, err
		}
		out = append(out, field)
	}
	return out, nil
}

// Files returns the fields of the owner that reach a file, the required ones
// before the optional, each with the way it reaches one. A field typed as a
// built-in reaches nothing, since only a definition is ever marked. It fails
// when a field's type names no definition.
func (f Fields) Files() ([]FileField, error) {
	out := make([]FileField, 0)
	for _, record := range f.records() {
		named, ok := record.Type.Atom().(typeform.Named)
		if !ok {
			continue
		}
		mark, reaches := f.db.Files.Lookup(named.Ref())
		if !reaches {
			continue
		}
		field, err := f.field(record)
		if err != nil {
			return nil, err
		}
		out = append(out, FileField{Field: field, Kind: mark.Kind})
	}
	return out, nil
}

// field returns the record read as a field, with its type bound to the
// definition it names. It fails when that type names no definition.
func (f Fields) field(record flattened.Field) (Field, error) {
	typ, err := NewResolution(f.db, record.Type).Value()
	if err != nil {
		return Field{}, fmt.Errorf("binding the type of %s.%s: %w", f.owner, record.Key, err)
	}
	return Field{
		Key:         record.Key,
		Type:        typ,
		Optionality: record.Optionality,
		Description: record.Description,
	}, nil
}

// records returns the raw field records of the owner, grouped by optionality
// and ordered by position within either group.
func (f Fields) records() []flattened.Field {
	out := make([]flattened.Field, 0)
	for key, record := range f.db.Fields.All() {
		if key.Owner != f.owner {
			continue
		}
		out = append(out, record)
	}
	slices.SortFunc(out, func(a, b flattened.Field) int {
		return cmp.Or(
			cmp.Compare(f.rank(a.Optionality), f.rank(b.Optionality)),
			cmp.Compare(a.Position, b.Position),
		)
	})
	return out
}

// rank returns the group an optionality sorts into, ordering a required field
// ahead of an optional one.
func (f Fields) rank(opt model.Optionality) int {
	if opt {
		return 1
	}
	return 0
}
