// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package implicit

// Spec represents the set of tgen-defined types not present in the Telegram Bot API specification.
type Spec struct {
	unions Unions
}

// NewSpec constructs the canonical Spec of implicit types.
func NewSpec() Spec {
	return Spec{
		unions: Unions{
			discriminated: []DiscriminatedUnion{
				NewDiscriminatedUnion(
					"ChatId",
					"ChatId represents a chat identifier, which is either an integer or a string.",
					[]DiscriminatedVariant{
						NewDiscriminatedVariant("Id", "Integer"),
						NewDiscriminatedVariant("Username", "String"),
					},
				),
				NewDiscriminatedUnion(
					"ReplyMarkup",
					"ReplyMarkup represents a reply markup attached to a message.",
					[]DiscriminatedVariant{
						NewDiscriminatedVariant("InlineKeyboardMarkup", "InlineKeyboardMarkup"),
						NewDiscriminatedVariant("ReplyKeyboardMarkup", "ReplyKeyboardMarkup"),
						NewDiscriminatedVariant("ReplyKeyboardRemove", "ReplyKeyboardRemove"),
						NewDiscriminatedVariant("ForceReply", "ForceReply"),
					},
				),
				NewDiscriminatedUnion(
					"InputMediaGroup",
					"InputMediaGroup represents a media element in a media group.",
					[]DiscriminatedVariant{
						NewDiscriminatedVariant("InputMediaAudio", "InputMediaAudio"),
						NewDiscriminatedVariant("InputMediaDocument", "InputMediaDocument"),
						NewDiscriminatedVariant("InputMediaPhoto", "InputMediaPhoto"),
						NewDiscriminatedVariant("InputMediaVideo", "InputMediaVideo"),
					},
				),
				NewDiscriminatedUnion(
					"MaybeMessage",
					"MaybeMessage represents a method return type that is either an edited Message or True for inline messages.",
					[]DiscriminatedVariant{
						NewDiscriminatedVariant("Message", "Message"),
						NewDiscriminatedVariant("Ok", "Boolean"),
					},
				),
			},
		},
	}
}

// Unions returns the implicit union types.
func (s Spec) Unions() Unions {
	return s.unions
}
