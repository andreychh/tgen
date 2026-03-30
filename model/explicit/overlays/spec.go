// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
)

// Spec represents the Telegram Bot API specification with field and method type overlays applied.
type Spec struct {
	inner   explicit.Specification
	overlay Overlay
}

// NewSpec constructs a Spec from a parsed specification.
func NewSpec(s explicit.Specification) Spec {
	return Spec{
		inner: s,
		overlay: NewSequential(
			ChatID{},
			ReplyMarkup{},
			InputMediaGroup{},
			InputFile{},
		),
	}
}

func (s Spec) Objects() iter.Seq[explicit.Object] {
	return func(yield func(explicit.Object) bool) {
		for o := range s.inner.Objects() {
			if !yield(NewObject(o, s.overlay)) {
				break
			}
		}
	}
}

func (s Spec) Methods() iter.Seq[explicit.Method] {
	return func(yield func(explicit.Method) bool) {
		for m := range s.inner.Methods() {
			if !yield(NewMethod(m, s.overlay)) {
				break
			}
		}
	}
}

func (s Spec) Unions() explicit.Unions {
	return s.inner.Unions()
}

func (s Spec) Release() explicit.Release {
	return s.inner.Release()
}
