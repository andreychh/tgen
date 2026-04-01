// SPDX-FileCopyrightText: 2026 Andrey Chernykh
// SPDX-License-Identifier: MIT

package literals

import "github.com/andreychh/tgen/model/implicit"

// Specification represents the set of tgen-defined types not present in the
// Telegram Bot API specification.
type Specification struct {
	unions Unions
}

// NewSpecification constructs the canonical Specification of implicit types.
func NewSpecification() Specification {
	return Specification{
		unions: Unions{
			discriminated: []implicit.DiscriminatedUnion{
				NewDiscriminatedUnion(
					"InputMediaGroup",
					"InputMediaGroup represents a media element in a media group.",
					"type",
					[]implicit.DiscriminatedVariant{
						NewDiscriminatedVariant("InputMediaAudio", "audio"),
						NewDiscriminatedVariant("InputMediaDocument", "document"),
						NewDiscriminatedVariant("InputMediaPhoto", "photo"),
						NewDiscriminatedVariant("InputMediaVideo", "video"),
					},
				),
			},
			structured: []implicit.StructuredUnion{
				NewStructuredUnion(
					"ReplyMarkup",
					"ReplyMarkup represents a reply markup attached to a message.",
					[]implicit.StructuredVariant{
						NewStructuredVariant("InlineKeyboardMarkup", "InlineKeyboardMarkup"),
						NewStructuredVariant("ReplyKeyboardMarkup", "ReplyKeyboardMarkup"),
						NewStructuredVariant("ReplyKeyboardRemove", "ReplyKeyboardRemove"),
						NewStructuredVariant("ForceReply", "ForceReply"),
					},
				),
				NewStructuredUnion(
					"ChatId",
					"ChatId represents a chat identifier, which is either an integer or a string.",
					[]implicit.StructuredVariant{
						NewStructuredVariant("Id", "Integer"),
						NewStructuredVariant("Username", "String"),
					},
				),
				NewStructuredUnion(
					"MaybeMessage",
					"MaybeMessage represents a method return type that is either an edited Message or True for inline messages.",
					[]implicit.StructuredVariant{
						NewStructuredVariant("Message", "Message"),
						NewStructuredVariant("Ok", "Boolean"),
					},
				),
			},
		},
	}
}

// Unions returns the implicit union types.
func (s Specification) Unions() implicit.Unions {
	return s.unions
}
