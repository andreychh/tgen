// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/explicit"
	"github.com/andreychh/tgen/pkg/iters"
)

// Specification represents the Telegram Bot API specification with field and
// method type overlays applied.
type Specification struct {
	inner   explicit.Specification
	overlay Overlay
}

// NewSpecification constructs a Specification from a parsed specification.
func NewSpecification(s explicit.Specification) Specification {
	return Specification{
		inner: s,
		overlay: NewSequential(
			ChatID{},
			ReplyMarkup{},
			InputMediaGroup{},
			InputFile{},
		),
	}
}

func (s Specification) Objects() iter.Seq[explicit.Object] {
	return func(yield func(explicit.Object) bool) {
		for o := range s.inner.Objects() {
			name, _ := o.Name().AsString()
			if name == "InputFile" {
				continue
			}
			if !yield(NewObject(o, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[explicit.Method] {
	return func(yield func(explicit.Method) bool) {
		for m := range s.inner.Methods() {
			if !yield(NewMethod(m, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedObjects() iter.Seq[explicit.DiscriminatedObject] {
	return func(yield func(explicit.DiscriminatedObject) bool) {
		for v := range s.inner.DiscriminatedObjects() {
			name, _ := v.Name().AsString()
			// InaccessibleMessage uses an integer discriminator (date == 0),
			// which the current string-only discriminator system cannot represent.
			// It is hardcoded in the template until typed discriminators are supported.
			if name == "InaccessibleMessage" {
				continue
			}
			if !yield(NewDiscriminatedObject(v, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedUnions() iter.Seq[explicit.DiscriminatedUnion] {
	return iters.NewMappedSeq(
		s.inner.DiscriminatedUnions(),
		func(d explicit.DiscriminatedUnion) explicit.DiscriminatedUnion {
			return NewDiscriminatedUnion(d, s.overlay)
		},
	)
}

func (s Specification) Release() explicit.Release {
	return s.inner.Release()
}
