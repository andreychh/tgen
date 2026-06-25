// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package parsed

import (
	"fmt"
	"regexp"

	"github.com/PuerkitoBio/goquery"

	"github.com/andreychh/tgen/model"
)

// keyPattern matches a field or parameter key: a snake_case identifier.
var keyPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Key is the first cell of a field or parameter row, holding its snake_case key.
type Key struct {
	td *goquery.Selection
}

// NewKey constructs a Key over a key table cell.
func NewKey(td *goquery.Selection) Key {
	return Key{td: td}
}

// Value returns the key. It fails when the cell is not a snake_case identifier.
func (k Key) Value() (model.Key, error) {
	key := k.td.Text()
	if !keyPattern.MatchString(key) {
		return "", fmt.Errorf("key %q is not a snake_case identifier", key)
	}
	return model.Key(key), nil
}
