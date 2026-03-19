// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package golang

// VariantType represents the Go type for an ImplicitVariant field, resolving
// primitive Telegram API type names to their Go equivalents.
type VariantType struct {
	source RawValue
}

// NewVariantType constructs a VariantType from an ObjectName.
func NewVariantType(n RawValue) VariantType {
	return VariantType{source: n}
}

func (t VariantType) Value() (string, error) {
	name, err := t.source.Value()
	if err != nil {
		return "", err
	}
	if goName, ok := namedTypes[name]; ok {
		return goName, nil
	}
	return NewDefaultName(NewStaticName(name)).Value()
}
