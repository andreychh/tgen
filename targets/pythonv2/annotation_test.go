// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/primitive"
	"github.com/andreychh/tgen/model/typebound"
	"github.com/andreychh/tgen/model/typeform"
	"github.com/andreychh/tgen/targets/pythonv2"
)

func bound(atom typebound.Atom, dim typeform.Dimensionality) typebound.Type {
	return typebound.NewType(atom, dim)
}

func TestAnnotation_Value(t *testing.T) {
	cases := []struct {
		name string
		typ  typebound.Type
		opt  model.Optionality
		want string
	}{
		{
			name: "writes a required integer as the built-in Python counts with",
			typ:  bound(typebound.NewPrimitive(primitive.Integer), 0),
			opt:  false,
			want: "int",
		},
		{
			name: "writes a required float as the built-in Python measures with",
			typ:  bound(typebound.NewPrimitive(primitive.Float), 0),
			opt:  false,
			want: "float",
		},
		{
			name: "writes a required string as the built-in Python spells text with",
			typ:  bound(typebound.NewPrimitive(primitive.String), 0),
			opt:  false,
			want: "str",
		},
		{
			name: "writes the True of the documentation as a plain boolean",
			typ:  bound(typebound.NewPrimitive(primitive.True), 0),
			opt:  false,
			want: "bool",
		},
		{
			name: "admits None where a primitive field is optional",
			typ:  bound(typebound.NewPrimitive(primitive.Boolean), 0),
			opt:  true,
			want: "bool | None",
		},
		{
			name: "writes an object atom as the class name it declares",
			typ:  bound(typebound.NewObject("MessageEntity"), 0),
			opt:  false,
			want: "MessageEntity",
		},
		{
			name: "restores the capitals of an acronym in an object name",
			typ:  bound(typebound.NewObject("MessageId"), 0),
			opt:  false,
			want: "MessageID",
		},
		{
			name: "writes a union atom as the class name it declares",
			typ:  bound(typebound.NewUnion("MaybeInaccessibleMessage"), 0),
			opt:  false,
			want: "MaybeInaccessibleMessage",
		},
		{
			name: "writes an alias atom by its own name and not by what it stands for",
			typ: bound(typebound.NewAlias(
				"RichTextPlain",
				bound(typebound.NewPrimitive(primitive.String), 0),
			), 0),
			opt:  false,
			want: "RichTextPlain",
		},
		{
			name: "encloses an atom of one dimension in a list",
			typ:  bound(typebound.NewObject("PhotoSize"), 1),
			opt:  false,
			want: "list[PhotoSize]",
		},
		{
			name: "encloses an atom of two dimensions in a list per dimension",
			typ:  bound(typebound.NewObject("InlineKeyboardButton"), 2),
			opt:  false,
			want: "list[list[InlineKeyboardButton]]",
		},
		{
			name: "admits None beside the outermost list and not inside it",
			typ:  bound(typebound.NewObject("MessageEntity"), 1),
			opt:  true,
			want: "list[MessageEntity] | None",
		},
		{
			name: "admits None beside a nested list where the field is optional",
			typ:  bound(typebound.NewPrimitive(primitive.String), 3),
			opt:  true,
			want: "list[list[list[str]]] | None",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := pythonv2.NewAnnotation(c.typ, c.opt).Value()
			assert.Equal(
				t,
				c.want,
				got,
				"Annotation must enclose the atom once per dimension and admit None only where the field is optional",
			)
		})
	}
}
