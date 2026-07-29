// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package attached

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
)

// Kind tells apart the two ways a definition has to do with a file: KindFile
// marks the type a file is sent as, which is one definition and no other;
// KindCarrier marks a definition holding one somewhere inside, however deep.
// A caller reads the kind to know whether a value hands over its file itself
// or asks what it holds to hand over theirs.
type Kind string

const (
	KindFile    Kind = "file"
	KindCarrier Kind = "carrier"
)

// File is what one definition has to do with a file.
type File struct {
	Ref  model.Reference
	Kind Kind
}

// FileTable is the spreading operator: it starts from the type a file is sent
// as and marks everything the rules can reach from there.
type FileTable struct {
	seed  model.Reference
	rules []Rule
}

// NewFileTable constructs a FileTable spreading from seed by the given rules.
func NewFileTable(seed model.Reference, rules ...Rule) FileTable {
	return FileTable{seed: seed, rules: rules}
}

// Table returns every definition that can carry a file, keyed by reference. The
// spread repeats until a round marks nothing new, so a definition holding its
// own kind — a rich block nesting rich blocks — settles instead of circling.
func (t FileTable) Table() Files {
	out := pipeline.NewMapTable[model.Reference, File]()
	out.Insert(t.seed, File{Ref: t.seed, Kind: KindFile})
	for t.spread(out) {
	}
	return out
}

// spread runs every rule once against out and reports whether the round marked
// a definition out did not already hold.
func (t FileTable) spread(out pipeline.MapTable[model.Reference, File]) bool {
	grew := false
	for _, rule := range t.rules {
		for _, ref := range rule.Apply(out) {
			if _, marked := out.Lookup(ref); marked {
				continue
			}
			out.Insert(ref, File{Ref: ref, Kind: KindCarrier})
			grew = true
		}
	}
	return grew
}
