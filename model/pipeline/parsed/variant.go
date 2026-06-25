// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
)

// Variant is the decoded record of one union member: the reference of the type
// it points to. Its owning union is the key.
type Variant struct {
	Ref model.Reference
}

// VariantItem is one variant's <li> item of a union's list.
type VariantItem struct {
	li *goquery.Selection
}

// NewVariantItem constructs a VariantItem over a variant's <li> item.
func NewVariantItem(li *goquery.Selection) VariantItem {
	return VariantItem{li: li}
}

// Record returns the variant decoded from the item: the reference it points to.
// It fails when the item has no link or the reference is malformed.
func (i VariantItem) Record() (Variant, error) {
	href, found := i.li.Find("a").Attr("href")
	if !found {
		return Variant{}, errors.New("variant href not found")
	}
	ref := strings.TrimPrefix(href, "#")
	if !refPattern.MatchString(ref) {
		return Variant{}, fmt.Errorf("variant reference %q is malformed", ref)
	}
	return Variant{Ref: model.Reference(ref)}, nil
}

// VariantItems are the variant items of every union on a documentation page.
type VariantItems struct {
	doc *goquery.Document
}

// NewVariantItems constructs a VariantItems over a parsed documentation page.
func NewVariantItems(doc *goquery.Document) VariantItems {
	return VariantItems{doc: doc}
}

// Table returns the variants table, one record per variant item, keyed by owning
// union and variant reference. It fails when any reference or variant item is
// malformed.
func (i VariantItems) Table() (pipeline.MapTable[model.VariantKey, Variant], error) {
	out := pipeline.NewMapTable[model.VariantKey, Variant]()
	for _, h4 := range i.doc.Find("h4").EachIter() {
		if NewHeading(h4).Kind() != KindUnion {
			continue
		}
		owner, err := NewReference(h4).Value()
		if err != nil {
			return out, fmt.Errorf("parsing union reference: %w", err)
		}
		items := h4.NextUntil("h3, h4, hr").Filter("ul").First().Find("li")
		for _, li := range items.EachIter() {
			variant, err := NewVariantItem(li).Record()
			if err != nil {
				return out, fmt.Errorf("parsing variant: %w", err)
			}
			out.Insert(model.VariantKey{Owner: owner, Ref: variant.Ref}, variant)
		}
	}
	return out, nil
}
