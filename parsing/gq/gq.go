// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package gq

import "iter"

type Selection interface {
	Text() string
	Attr(name string) (string, bool)
	Find(selector string) Selection
	Filter(selector string) Selection
	FilterFunc(f func(Selection) bool) Selection
	Until(selector string) Selection
	IsEmpty() bool
	At(index int) Selection
	All() iter.Seq[Selection]
}
