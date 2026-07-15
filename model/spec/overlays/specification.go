// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package overlays

import (
	"iter"

	"github.com/andreychh/tgen/model/spec"
	"github.com/andreychh/tgen/pkg/iters"
)

// Specification represents the Telegram Bot API specification with field and
// method type overlays applied.
type Specification struct {
	inner   spec.Specification
	overlay Overlay
}

// NewSpecification constructs a Specification from a parsed specification.
func NewSpecification(s spec.Specification) Specification {
	return Specification{
		inner: s,
		overlay: NewSequential(
			ChatID{},
			ReplyMarkup{},
			InputMediaGroup{},
			InputRichMedia{},
			InputFile{},
		),
	}
}

func (s Specification) Objects() iter.Seq[spec.Object] {
	return func(yield func(spec.Object) bool) {
		for o := range s.inner.Objects() {
			name, _ := o.Name()
			if name == "InputFile" {
				continue
			}
			if !yield(NewObject(o, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) Methods() iter.Seq[spec.Method] {
	return func(yield func(spec.Method) bool) {
		for m := range s.inner.Methods() {
			if !yield(NewMethod(m, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedObjects() iter.Seq[spec.DiscriminatedObject] {
	return func(yield func(spec.DiscriminatedObject) bool) {
		for obj := range s.inner.DiscriminatedObjects() {
			name, _ := obj.Name()
			// InaccessibleMessage uses an integer discriminator (date == 0),
			// which the current string-only discriminator system cannot represent.
			// It is hardcoded in the template until typed discriminators are supported.
			if name == "InaccessibleMessage" {
				continue
			}
			if !yield(NewDiscriminatedObject(obj, s.overlay)) {
				break
			}
		}
	}
}

func (s Specification) DiscriminatedUnions() iter.Seq[spec.DiscriminatedUnion] {
	return iters.NewMappedSeq(
		s.inner.DiscriminatedUnions(),
		func(d spec.DiscriminatedUnion) spec.DiscriminatedUnion {
			return NewDiscriminatedUnion(d, s.overlay)
		},
	)
}

func (s Specification) Release() spec.Release {
	return s.inner.Release()
}
