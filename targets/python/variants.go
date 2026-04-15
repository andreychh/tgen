// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package python

import (
	"fmt"
	"iter"
	"strings"

	"github.com/andreychh/tgen/model/explicit"
)

type Variants struct {
	inner iter.Seq[explicit.DiscriminatedObject]
}

func NewVariants(s iter.Seq[explicit.DiscriminatedObject]) Variants {
	return Variants{inner: s}
}

func (v Variants) AsString() (string, error) {
	var names []string
	for o := range v.inner {
		name, err := NewClassName(o.Name()).AsString()
		if err != nil {
			return "", fmt.Errorf("variant name: %w", err)
		}
		names = append(names, name)
	}
	return strings.Join(names, " | "), nil
}
