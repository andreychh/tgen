// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package attached

import (
	"github.com/andreychh/tgen/model"
	"github.com/andreychh/tgen/model/pipeline"
)

// File is what one definition has to do with a file, classified by
// [model.FileKind].
type File struct {
	Ref  model.Reference
	Kind model.FileKind
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
	out.Insert(t.seed, File{Ref: t.seed, Kind: model.FileKindFile})
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
			out.Insert(ref, File{Ref: ref, Kind: model.FileKindCarrier})
			grew = true
		}
	}
	return grew
}
