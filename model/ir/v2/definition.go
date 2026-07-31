// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package ir

// Definition represents one thing the specification names, joined with what it
// owns. A target walks the sum rather than a sequence per kind, because the
// page numbers its headings in a single run: an object, a union and a method
// stand beside each other there, and only a sum can hold that order. The
// concrete variants are [Object], [DiscriminatedObject], [Union],
// [DiscriminatedUnion], [Alias] and [Method].
//
//sumtype:decl
type Definition interface {
	isDefinition()
}
