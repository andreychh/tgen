// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

import (
	"fmt"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/corrected"
	"github.com/andreychh/tgen/model/pipeline/separated"
)

// Reading is one definition read as the record a target renders: it joins the
// definition with what the tables hold for it, and gives back the record of the
// kind the definition is of.
type Reading struct {
	db         separated.Specification
	definition corrected.Definition
}

// NewReading constructs a Reading of one definition against the database
// holding what that definition owns.
func NewReading(db separated.Specification, definition corrected.Definition) Reading {
	return Reading{db: db, definition: definition}
}

// Value returns the definition read as the record of its kind. An object is
// read as the variant its discriminator makes it, since a union telling it apart
// changes what a target writes for it. It fails when the join fails, and when no
// record stands for the kind the definition is of.
func (r Reading) Value() (Definition, error) {
	switch r.definition.Kind {
	case model.DefinitionKindObject:
		if _, discriminated := r.db.Discriminators.Lookup(r.definition.Ref); discriminated {
			return r.discriminatedObject()
		}
		return r.object()
	case model.DefinitionKindMethod:
		return r.method()
	case model.DefinitionKindAlias:
		return r.alias()
	case model.DefinitionKindUnion:
		return r.union()
	default:
		return nil, fmt.Errorf(
			"%s names a definition of unknown kind %q",
			r.definition.Ref,
			r.definition.Kind,
		)
	}
}

// object returns the definition joined with the fields it owns.
func (r Reading) object() (Object, error) {
	fields, err := NewFields(r.db, r.definition.Ref).Value()
	if err != nil {
		return Object{}, fmt.Errorf("joining object %s: %w", r.definition.Ref, err)
	}
	files, err := NewFields(r.db, r.definition.Ref).Files()
	if err != nil {
		return Object{}, fmt.Errorf("joining object %s: %w", r.definition.Ref, err)
	}
	direction, err := r.direction()
	if err != nil {
		return Object{}, fmt.Errorf("joining object %s: %w", r.definition.Ref, err)
	}
	return Object{
		Ref:         r.definition.Ref,
		Name:        r.definition.Name,
		Description: r.definition.Description,
		Fields:      fields,
		Files:       files,
		Rewrites:    r.rewrites(),
		Direction:   direction,
		Introduced:  r.definition.Introduced,
	}, nil
}

// direction returns which way the definition travels. It fails when the
// specification names no direction for it, which only a method is left without,
// so failing here means a method was read as something it is not.
func (r Reading) direction() (model.Direction, error) {
	direction, found := r.db.Directions.Lookup(r.definition.Ref)
	if !found {
		return "", fmt.Errorf("%s travels nowhere", r.definition.Ref)
	}
	return direction, nil
}

// rewrites reports whether a union reaching a file admits the definition, which
// obliges it to rewrite itself into JSON however little it has to hand over.
func (r Reading) rewrites() bool {
	for key := range r.db.Variants.All() {
		if key.Ref != r.definition.Ref {
			continue
		}
		mark, reaches := r.db.Files.Lookup(key.Owner)
		if reaches && mark.Kind == model.FileKindCarrier {
			return true
		}
	}
	return false
}

// discriminatedObject returns the definition joined with the fields it owns and
// the fixed value it is told apart by. It fails when no union tells it apart.
func (r Reading) discriminatedObject() (DiscriminatedObject, error) {
	discriminator, discriminated := r.db.Discriminators.Lookup(r.definition.Ref)
	if !discriminated {
		return DiscriminatedObject{}, fmt.Errorf(
			"object %s is told apart by nothing",
			r.definition.Ref,
		)
	}
	fields, err := NewFields(r.db, r.definition.Ref).Value()
	if err != nil {
		return DiscriminatedObject{}, fmt.Errorf(
			"joining discriminated object %s: %w",
			r.definition.Ref,
			err,
		)
	}
	files, err := NewFields(r.db, r.definition.Ref).Files()
	if err != nil {
		return DiscriminatedObject{}, fmt.Errorf(
			"joining discriminated object %s: %w",
			r.definition.Ref,
			err,
		)
	}
	direction, err := r.direction()
	if err != nil {
		return DiscriminatedObject{}, fmt.Errorf(
			"joining discriminated object %s: %w",
			r.definition.Ref,
			err,
		)
	}
	return DiscriminatedObject{
		Ref:         r.definition.Ref,
		Name:        r.definition.Name,
		Description: r.definition.Description,
		Fields:      fields,
		Files:       files,
		Rewrites:    r.rewrites(),
		Direction:   direction,
		Introduced:  r.definition.Introduced,
		Discriminator: Discriminator{
			Key:   discriminator.Key,
			Value: discriminator.Value,
		},
	}, nil
}

// union returns the definition joined with the variants it stands for: told
// apart by a key when one key tells every variant apart with a value of its own,
// and told apart by nothing otherwise, since then only a target knows how. It
// fails when a variant names no definition.
func (r Reading) union() (Definition, error) {
	mark, reaches := r.db.Files.Lookup(r.definition.Ref)
	carrier := reaches && mark.Kind == model.FileKindCarrier
	direction, err := r.direction()
	if err != nil {
		return nil, fmt.Errorf("joining union %s: %w", r.definition.Ref, err)
	}
	key, discriminated, err := NewVariants(r.db, r.definition.Ref).Discriminated()
	if err != nil {
		return nil, fmt.Errorf("joining union %s: %w", r.definition.Ref, err)
	}
	if len(discriminated) > 0 {
		return DiscriminatedUnion{
			Ref:         r.definition.Ref,
			Name:        r.definition.Name,
			Description: r.definition.Description,
			Key:         key,
			Variants:    discriminated,
			Carrier:     carrier,
			Direction:   direction,
			Introduced:  r.definition.Introduced,
		}, nil
	}
	variants, err := NewVariants(r.db, r.definition.Ref).Value()
	if err != nil {
		return nil, fmt.Errorf("joining union %s: %w", r.definition.Ref, err)
	}
	return Union{
		Ref:         r.definition.Ref,
		Name:        r.definition.Name,
		Description: r.definition.Description,
		Variants:    variants,
		Carrier:     carrier,
		Direction:   direction,
		Introduced:  r.definition.Introduced,
	}, nil
}

// method returns the definition joined with the parameters it takes and the
// result it returns. It fails when the pipeline separated it into no result.
func (r Reading) method() (Method, error) {
	record, found := r.db.Methods.Lookup(r.definition.Ref)
	if !found {
		return Method{}, fmt.Errorf("method %s returns nothing", r.definition.Ref)
	}
	params, err := NewFields(r.db, r.definition.Ref).Value()
	if err != nil {
		return Method{}, fmt.Errorf("joining method %s: %w", r.definition.Ref, err)
	}
	files, err := NewFields(r.db, r.definition.Ref).Files()
	if err != nil {
		return Method{}, fmt.Errorf("joining method %s: %w", r.definition.Ref, err)
	}
	res, err := NewReturned(r.db, record.Result).Value()
	if err != nil {
		return Method{}, fmt.Errorf("joining method %s: %w", r.definition.Ref, err)
	}
	return Method{
		Ref:         r.definition.Ref,
		Name:        r.definition.Name,
		Description: r.definition.Description,
		Params:      params,
		Files:       files,
		Result:      res,
		Introduced:  r.definition.Introduced,
	}, nil
}

// alias returns the definition joined with the type it stands for. It fails when
// it stands for nothing.
func (r Reading) alias() (Alias, error) {
	record, found := r.db.Aliases.Lookup(r.definition.Ref)
	if !found {
		return Alias{}, fmt.Errorf("alias %s stands for nothing", r.definition.Ref)
	}
	typ, err := NewResolution(r.db, record.Type).Value()
	if err != nil {
		return Alias{}, fmt.Errorf("resolving alias %s: %w", r.definition.Ref, err)
	}
	direction, err := r.direction()
	if err != nil {
		return Alias{}, fmt.Errorf("joining alias %s: %w", r.definition.Ref, err)
	}
	return Alias{
		Ref:         r.definition.Ref,
		Name:        r.definition.Name,
		Type:        typ,
		Description: r.definition.Description,
		Direction:   direction,
	}, nil
}
