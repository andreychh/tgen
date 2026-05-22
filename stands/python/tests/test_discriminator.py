# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import json

import pytest

from api.methods import GetChatMethod, SetMessageReactionMethod
from api.types import (
    ID,
    ReactionType,
    ReactionTypeCustomEmoji,
    ReactionTypeEmoji,
    ReactionTypePaid,
)

from helpers import CapturingConnection, RespondingConnection

_CHAT_FULL_INFO_BASE = {
    "id": 1,
    "type": "private",
    "accent_color_id": 0,
    "max_reaction_count": 0,
    "accepted_gift_types": {
        "unlimited_gifts": False,
        "limited_gifts": False,
        "unique_gifts": False,
        "premium_subscription": False,
        "gifts_from_channels": False,
    },
}


@pytest.mark.parametrize("reaction,want_type", [
    (ReactionTypeEmoji(emoji="👍"), "emoji"),
    (ReactionTypeCustomEmoji(custom_emoji_id="id123"), "custom_emoji"),
    (ReactionTypePaid(), "paid"),
])
def test_set_message_reaction_method_call_injects_type_discriminator(
    reaction: ReactionType, want_type: str
) -> None:
    conn = CapturingConnection()
    SetMessageReactionMethod(
        chat_id=ID(1), message_id=1, reaction=[reaction]
    ).call(conn)
    assert json.loads(conn.data["reaction"])[0]["type"] == want_type


@pytest.mark.parametrize("reaction_json,want_type", [
    ({"type": "emoji", "emoji": "👍"}, ReactionTypeEmoji),
    ({"type": "custom_emoji", "custom_emoji_id": "id123"}, ReactionTypeCustomEmoji),
    ({"type": "paid"}, ReactionTypePaid),
])
def test_get_chat_method_call_dispatches_available_reaction_to_correct_variant(
    reaction_json: dict, want_type: type
) -> None:
    data = json.dumps({**_CHAT_FULL_INFO_BASE, "available_reactions": [reaction_json]}).encode()
    conn = RespondingConnection(data)
    result = GetChatMethod(chat_id=ID(99)).call(conn)
    assert isinstance(result.available_reactions[0], want_type)
