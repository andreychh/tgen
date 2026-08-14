// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/andreychh/tgen/model/pipeline/typed/types"
)

func TestCursor_Peek(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		peeks int
		want  []string
	}{
		{
			name:  "shows the same item however many times it is asked",
			items: []string{"Array", "of"},
			peeks: 3,
			want:  []string{"Array", "Array", "Array"},
		},
		{
			name:  "shows the zero value when constructed over no items",
			items: nil,
			peeks: 2,
			want:  []string{"", ""},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := types.NewCursor(tc.items)
			got := make([]string, 0, tc.peeks)
			for range tc.peeks {
				item, _ := cursor.Peek()
				got = append(got, item)
			}
			assert.Equal(
				t,
				tc.want,
				got,
				"Cursor.Peek must show the item the cursor stands on without advancing past it",
			)
		})
	}
}

func TestCursor_Take(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		takes int
		want  []string
	}{
		{
			name:  "hands out every item in the order it was given",
			items: []string{"Array", "of", "Süßigkeit"},
			takes: 3,
			want:  []string{"Array", "of", "Süßigkeit"},
		},
		{
			name:  "hands out nothing once the last item is consumed",
			items: []string{"или", "и"},
			takes: 4,
			want:  []string{"или", "и"},
		},
		{
			name:  "hands out nothing when constructed over no items",
			items: nil,
			takes: 2,
			want:  nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := types.NewCursor(tc.items)
			var got []string
			for range tc.takes {
				item, ok := cursor.Take()
				if !ok {
					break
				}
				got = append(got, item)
			}
			assert.Equal(
				t,
				tc.want,
				got,
				"Cursor.Take must hand out each item once, front to back, and report none after the last",
			)
		})
	}
}

func TestCursor_Skip(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		skips int
		want  []string
	}{
		{
			name:  "passes over the item the cursor stands on",
			items: []string{"Array", "of", "Süßigkeit"},
			skips: 1,
			want:  []string{"of", "Süßigkeit"},
		},
		{
			name:  "leaves an exhausted cursor exhausted",
			items: []string{"или"},
			skips: 5,
			want:  nil,
		},
		{
			name:  "leaves a cursor over no items exhausted",
			items: nil,
			skips: 3,
			want:  nil,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := types.NewCursor(tc.items)
			for range tc.skips {
				cursor.Skip()
			}
			var got []string
			for {
				item, ok := cursor.Take()
				if !ok {
					break
				}
				got = append(got, item)
			}
			assert.Equal(
				t,
				tc.want,
				got,
				"Cursor.Skip must pass over one item and leave nothing behind it reachable",
			)
		})
	}
}

func TestCursor_Done(t *testing.T) {
	cases := []struct {
		name  string
		items []string
		takes int
		want  bool
	}{
		{
			name:  "reports items left while one is still unconsumed",
			items: []string{"Array", "of"},
			takes: 1,
			want:  false,
		},
		{
			name:  "reports no items left once the last one is consumed",
			items: []string{"Array", "of"},
			takes: 2,
			want:  true,
		},
		{
			name:  "reports no items left when constructed over no items",
			items: nil,
			takes: 0,
			want:  true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cursor := types.NewCursor(tc.items)
			for range tc.takes {
				cursor.Take()
			}
			assert.Equal(
				t,
				tc.want,
				cursor.Done(),
				"Cursor.Done must report whether any item is still there to consume",
			)
		})
	}
}
