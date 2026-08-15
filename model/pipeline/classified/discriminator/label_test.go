// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package discriminator_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline/classified/discriminator"
	"github.com/andreychh/tgen/model/prose"
)

func TestLabel_Value(t *testing.T) {
	cases := []struct {
		name        string
		description prose.Phrase
		want        model.DiscriminatorValue
		wantOK      bool
	}{
		{
			name: "returns the value a must-be clause sets in italic",
			description: prose.NewPhrase(
				plain("Error source, must be "),
				italic("selfie"),
			),
			want:   model.DiscriminatorValue("selfie"),
			wantOK: true,
		},
		{
			name: "returns the value an always clause quotes",
			description: prose.NewPhrase(
				plain("Type of the area, always “suggested_reaction”"),
			),
			want:   model.DiscriminatorValue("suggested_reaction"),
			wantOK: true,
		},
		{
			name: "returns a quoted value written outside the ASCII range",
			description: prose.NewPhrase(
				plain("Тип области, always “отметка_погоды”"),
			),
			want:   model.DiscriminatorValue("отметка_погоды"),
			wantOK: true,
		},
		{
			name: "returns the italic value unchanged, keeping the space the run carries",
			description: prose.NewPhrase(
				plain("Error source, must be "),
				italic(" file "),
			),
			want:   model.DiscriminatorValue(" file "),
			wantOK: true,
		},
		{
			name: "returns the value of a must-be clause that runs on past it",
			description: prose.NewPhrase(
				plain("Error source, must be "),
				italic("front_side"),
				plain(", the side of the document"),
			),
			want:   model.DiscriminatorValue("front_side"),
			wantOK: true,
		},
		{
			name: "returns no value when the description names no fixed value",
			description: prose.NewPhrase(
				plain("Unique identifier of the message"),
			),
			wantOK: false,
		},
		{
			name: "returns no value when the quoted value does not close the description",
			description: prose.NewPhrase(
				plain("Source of the boost, always “premium” for a subscription"),
			),
			wantOK: false,
		},
		{
			name: "returns no value when the must-be words do not close the run before the italic",
			description: prose.NewPhrase(
				plain("The source must be one of the values below: "),
				italic("selfie"),
			),
			wantOK: false,
		},
		{
			name: "returns no value when the must-be clause carries no run after it",
			description: prose.NewPhrase(
				plain("Error source, must be "),
			),
			wantOK: false,
		},
		{
			name: "returns no value when the run after the must-be clause carries no emphasis",
			description: prose.NewPhrase(
				plain("Error source, must be "),
				plain("selfie"),
			),
			wantOK: false,
		},
		{
			name: "returns no value when the always clause sets its value in italic instead of quoting it",
			description: prose.NewPhrase(
				plain("Type of the area, always "),
				italic("link"),
			),
			wantOK: false,
		},
		{
			name: "returns no value when a link stands between the description and the must-be clause",
			description: prose.NewPhrase(
				plain("Source of the "),
				anchor("ChatBoost"),
				plain(", must be "),
				italic("premium"),
			),
			wantOK: false,
		},
		{
			name:        "returns no value when the description carries no runs",
			description: prose.NewPhrase(),
			wantOK:      false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := discriminator.NewLabel(tc.description).Value()
			if !tc.wantOK {
				assert.False(
					t,
					ok,
					"Label.Value must report no discriminator unless the description opens with one of the two forms",
				)
				return
			}
			require.True(t, ok)
			assert.Equal(
				t,
				tc.want,
				got,
				"Label.Value must decode the fixed value the description carries, whichever form writes it",
			)
		})
	}
}
