// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/spec/gq"
	pkggq "github.com/andreychh/tgen/pkg/gq"
)

func trWith(desc string) pkggq.Selection {
	doc, _ := goquery.NewDocumentFromReader(strings.NewReader(fmt.Sprintf(
		`<table><tbody><tr><td>type</td><td>String</td><td>%s</td></tr></tbody></table>`,
		desc,
	)))
	return pkggq.NewNormSelection(doc.Find("tr"))
}

func TestFieldRow_Kind(t *testing.T) {
	cases := []struct {
		name string
		desc string
		want gq.FieldKind
	}{
		{
			name: "returns FieldKindDiscriminator for Always with curly-quoted value",
			desc: "Always “animation”",
			want: gq.FieldKindDiscriminator,
		},
		{
			name: "returns FieldKindDiscriminator for always with lowercase a and curly-quoted value",
			desc: "always “audio”",
			want: gq.FieldKindDiscriminator,
		},
		{
			name: "returns FieldKindDiscriminator for Always followed by a numeric dot sequence",
			desc: "Always 42.",
			want: gq.FieldKindDiscriminator,
		},
		{
			name: "returns FieldKindDiscriminator for must be with a lowercase identifier",
			desc: "must be animation",
			want: gq.FieldKindDiscriminator,
		},
		{
			name: "returns FieldKindDiscriminator for must be with a snake_case identifier",
			desc: "must be audio_message",
			want: gq.FieldKindDiscriminator,
		},
		{
			name: "returns FieldKindFree for a plain description",
			desc: "Type of the service message.",
			want: gq.FieldKindFree,
		},
		{
			name: "returns FieldKindFree for Must be with uppercase M",
			desc: "Must be animation",
			want: gq.FieldKindFree,
		},
		{
			name: "returns FieldKindFree for must be followed by an uppercase identifier",
			desc: "must be Animation",
			want: gq.FieldKindFree,
		},
		{
			name: "returns FieldKindFree when Always is not followed by a quoted value or number",
			desc: "always animation",
			want: gq.FieldKindFree,
		},
		{
			name: "returns FieldKindFree when the description column is empty",
			desc: "",
			want: gq.FieldKindFree,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(
				t,
				tc.want,
				gq.NewFieldRow(trWith(tc.desc)).Kind(),
				"FieldRow must classify a row as discriminator only when its description matches the Always-quoted or must-be pattern",
			)
		})
	}
}
