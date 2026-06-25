// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package prose

// Tag is one of the closed set of HTML node kinds that occur in entity-section
// prose. The text node, which carries no tag, is [TagText]. A node outside the
// set — a tag that never appears in scope, such as <i>, <b>, or <pre> — is
// [TagUnknown], which the decoder rejects.
type Tag int

const (
	// TagUnknown marks a node whose tag falls outside entity-section prose.
	TagUnknown Tag = iota
	TagText
	TagEm
	TagStrong
	TagCode
	TagA
	TagBr
	TagImg
	TagP
	TagUl
	TagLi
	TagBlockquote
	TagDiv
)
