// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package pythonv2

import (
	"fmt"

	"github.com/andreychh/tgen/model"
)

// Assignment represents what a field declaration stands equal to: the None an
// optional field defaults to, and the key it is read from where the attribute
// it declares is not spelled the way the key is.
//
// The alias is written per field rather than derived for every field at once.
// One key of the documented API — from — is a word Python reserves, so a rule
// covering every model would be a rule the whole API pays for and one field
// uses, and it would misread the first key that ever ends in an underscore of
// its own.
type Assignment struct {
	key model.Key
	opt model.Optionality
}

// NewAssignment creates an Assignment for a field read from key, optional or
// not.
func NewAssignment(k model.Key, opt model.Optionality) Assignment {
	return Assignment{key: k, opt: opt}
}

// Value returns what stands to the right of the field's annotation, empty for a
// required field whose attribute is spelled the way its key is, which is what
// leaves the annotation standing alone.
func (a Assignment) Value() string {
	alias := a.alias()
	if alias == "" {
		if a.opt {
			return "None"
		}
		return ""
	}
	if a.opt {
		return fmt.Sprintf("Field(default=None, alias=%q)", alias)
	}
	return fmt.Sprintf("Field(alias=%q)", alias)
}

// alias returns the key the field is read from, empty where the attribute it
// declares already spells that key.
func (a Assignment) alias() string {
	if NewAttribute(a.key).Value() == string(a.key) {
		return ""
	}
	return string(a.key)
}
