// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package dom_test

import (
	"strings"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/andreychh/tgen/parsing/dom"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTMLSelection_All(t *testing.T) {
	tests := []struct {
		name string
		html string
		want []string
	}{
		{
			name: "yields all elements in document order",
			html: "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			want: []string{"One", "Two", "Three"},
		},
		{
			name: "yields the single matching element",
			html: "<ul><li>One</li></ul>",
			want: []string{"One"},
		},
		{
			name: "yields nothing from an empty selection",
			html: "<ul></ul>",
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			selection := dom.NewHTMLSelection(doc.Selection).Find("li")
			got := make([]string, len(tt.want))
			for i, item := range selection.All() {
				got[i] = item.Text()
			}
			assert.Equal(t, tt.want, got, "iterator does not yield the expected sequence")
		})
	}
}

func TestHTMLSelection_At(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		index    int
		wantText string
	}{
		{
			name:     "returns the first element",
			html:     "<ul><li>One</li></ul>",
			index:    0,
			wantText: "One",
		},
		{
			name:     "returns a middle element",
			html:     "<ul><li>One</li><li>Two</li></ul>",
			index:    1,
			wantText: "Two",
		},
		{
			name:     "returns the last element",
			html:     "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			index:    2,
			wantText: "Three",
		},
		{
			name:     "returns an empty selection for an out-of-bounds index",
			html:     "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			index:    5,
			wantText: "",
		},
		{
			name:     "returns an empty selection for a negative index",
			html:     "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			index:    -1,
			wantText: "",
		},
		{
			name:     "returns an empty selection from an empty list",
			html:     "<ul></ul>",
			index:    0,
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find("li").At(tt.index)
			assert.Equal(
				t,
				tt.wantText,
				got.Text(),
				"selection does not return the expected text at index %d",
				tt.index,
			)
		})
	}
}

func TestHTMLSelection_Attr(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		attrName  string
		wantValue string
		wantExist bool
	}{
		{
			name:      "returns value and true for existing attribute",
			html:      `<div id="my-id"></div>`,
			attrName:  "id",
			wantValue: "my-id",
			wantExist: true,
		},
		{
			name:      "normalizes whitespace in attribute value",
			html:      `<div class="  foo   bar  "></div>`,
			attrName:  "class",
			wantValue: "foo bar",
			wantExist: true,
		},
		{
			name:      "returns empty string and true for empty attribute",
			html:      `<div data-empty=""></div>`,
			attrName:  "data-empty",
			wantValue: "",
			wantExist: true,
		},
		{
			name:      "returns empty string and false for missing attribute",
			html:      `<div></div>`,
			attrName:  "href",
			wantValue: "",
			wantExist: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			gotValue, gotExist := dom.NewHTMLSelection(doc.Selection).Find("div").Attr(tt.attrName)
			assert.Equal(
				t,
				tt.wantExist,
				gotExist,
				"attribute existence flag does not match expectation",
			)
			assert.Equal(
				t,
				tt.wantValue,
				gotValue,
				"attribute value does not match expectation",
			)
		})
	}
}

func TestHTMLSelection_Filter(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		wantText string
	}{
		{
			name:     "filters elements by class selector",
			html:     `<ul><li class="keep">First</li><li class="drop">Second</li><li class="keep">Third</li></ul>`,
			selector: ".keep",
			wantText: "FirstThird",
		},
		{
			name:     "returns empty set for non-matching selector",
			html:     `<ul><li>Item</li></ul>`,
			selector: ".missing",
			wantText: "",
		},
		{
			name:     "filters elements by position pseudo-selector",
			html:     `<ul><li>First</li><li>Second</li></ul>`,
			selector: ":first-child",
			wantText: "First",
		},
		{
			name:     "filters elements by negation selector",
			html:     `<ul><li class="keep">First</li><li class="drop">Second</li></ul>`,
			selector: ":not(.drop)",
			wantText: "First",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find("li").Filter(tt.selector)
			assert.Equal(
				t,
				tt.wantText,
				got.Text(),
				"filtered selection text does not match expectation",
			)
		})
	}
}

func TestHTMLSelection_FilterFunc(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		predicate func(dom.Selection) bool
		wantText  string
	}{
		{
			name: "filters elements by specific text content",
			html: "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			predicate: func(s dom.Selection) bool {
				return s.Text() == "Two"
			},
			wantText: "Two",
		},
		{
			name: "filters elements by attribute value",
			html: `<ul><li class="odd">One</li><li class="even">Two</li><li class="odd">Three</li></ul>`,
			predicate: func(s dom.Selection) bool {
				val, ok := s.Attr("class")
				return ok && val == "odd"
			},
			wantText: "OneThree",
		},
		{
			name: "keeps all elements when predicate is true",
			html: "<ul><li>One</li><li>Two</li></ul>",
			predicate: func(s dom.Selection) bool {
				return true
			},
			wantText: "OneTwo",
		},
		{
			name: "removes all elements when predicate is false",
			html: "<ul><li>One</li><li>Two</li></ul>",
			predicate: func(s dom.Selection) bool {
				return false
			},
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find("li").FilterFunc(tt.predicate)
			assert.Equal(
				t,
				tt.wantText,
				got.Text(),
				"filtered selection text does not match expectation",
			)
		})
	}
}

func TestHTMLSelection_Find(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		wantText string
	}{
		{
			name:     "finds direct child elements",
			html:     `<div id="scope"><span class="unique">Unique Inner</span></div>`,
			selector: ".unique",
			wantText: "Unique Inner",
		},
		{
			name:     "finds nested descendants within scope",
			html:     `<div id="scope"><span class="target">Direct</span><div><span class="target">Deep</span></div></div><span class="target">Outside</span>`,
			selector: ".target",
			wantText: "DirectDeep",
		},
		{
			name:     "returns empty selection for missing descendants",
			html:     `<div id="scope"><span>Content</span></div>`,
			selector: ".non-existent",
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find("#scope").Find(tt.selector)
			assert.Equal(t, tt.wantText, got.Text(), "found elements do not match expectation")
		})
	}
}

func TestHTMLSelection_First(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		wantText string
	}{
		{
			name:     "returns first element from multiple items",
			html:     "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			wantText: "One",
		},
		{
			name:     "returns element itself when only one exists",
			html:     "<ul><li>One</li></ul>",
			wantText: "One",
		},
		{
			name:     "returns empty selection from empty set",
			html:     "<ul></ul>",
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find("li").First()
			assert.Equal(
				t,
				tt.wantText,
				got.Text(),
				"first element text does not match expectation",
			)
		})
	}
}

func TestHTMLSelection_IsEmpty(t *testing.T) {
	tests := []struct {
		name      string
		html      string
		selector  string
		wantEmpty bool
	}{
		{
			name:      "returns false for non-empty selection",
			html:      "<div>Content</div>",
			selector:  "div",
			wantEmpty: false,
		},
		{
			name:      "returns true for empty selection",
			html:      "<div>Content</div>",
			selector:  ".missing",
			wantEmpty: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find(tt.selector).IsEmpty()
			assert.Equal(t, tt.wantEmpty, got, "empty status does not match expectation")
		})
	}
}

func TestHTMLSelection_Length(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		wantLen  int
	}{
		{
			name:     "returns correct count for multiple elements",
			html:     "<ul><li>One</li><li>Two</li><li>Three</li></ul>",
			selector: "li",
			wantLen:  3,
		},
		{
			name:     "returns one for single matching element",
			html:     "<div>Content</div>",
			selector: "div",
			wantLen:  1,
		},
		{
			name:     "returns zero for empty selection",
			html:     "<div>Content</div>",
			selector: ".missing",
			wantLen:  0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find(tt.selector).Length()
			assert.Equal(t, tt.wantLen, got, "selection length does not match expectation")
		})
	}
}

func TestHTMLSelection_NextUntil(t *testing.T) {
	tests := []struct {
		name          string
		html          string
		startSelector string
		untilSelector string
		wantText      string
	}{
		{
			name:          "selects siblings until stopper",
			html:          `<dl><dt id="t1">Term</dt><dd>Def 1</dd><dd>Def 2</dd><dt>Next</dt></dl>`,
			startSelector: "#t1",
			untilSelector: "dt",
			wantText:      "Def 1Def 2",
		},
		{
			name:          "selects all remaining siblings if stopper missing",
			html:          `<dl><dt id="t1">Term</dt><dd>Def 1</dd><dd>Def 2</dd></dl>`,
			startSelector: "#t1",
			untilSelector: "h1",
			wantText:      "Def 1Def 2",
		},
		{
			name:          "returns empty if immediate sibling matches stopper",
			html:          `<dl><dt id="t1">Term</dt><dd>Def 1</dd></dl>`,
			startSelector: "#t1",
			untilSelector: "dd",
			wantText:      "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).
				Find(tt.startSelector).
				NextUntil(tt.untilSelector)
			assert.Equal(t, tt.wantText, got.Text(), "siblings text does not match expectation")
		})
	}
}

func TestHTMLSelection_Text(t *testing.T) {
	tests := []struct {
		name     string
		html     string
		selector string
		wantText string
	}{
		{
			name:     "returns text content of single element",
			html:     "<div>Hello World</div>",
			selector: "div",
			wantText: "Hello World",
		},
		{
			name:     "concatenates text of multiple elements",
			html:     "<span>One</span><span>Two</span>",
			selector: "span",
			wantText: "OneTwo",
		},
		{
			name:     "includes text of descendants",
			html:     "<div>Outer <span>Inner</span></div>",
			selector: "div",
			wantText: "Outer Inner",
		},
		{
			name:     "returns empty string for empty selection",
			html:     "<div></div>",
			selector: ".missing",
			wantText: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(tt.html))
			require.NoError(t, err, "HTML fixture does not parse correctly")
			got := dom.NewHTMLSelection(doc.Selection).Find(tt.selector).Text()
			assert.Equal(t, tt.wantText, got, "extracted text does not match expectation")
		})
	}
}
