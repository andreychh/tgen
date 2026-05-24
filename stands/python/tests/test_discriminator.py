# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import json

import pytest

from api.methods import GetChatAdministratorsMethod, GetChatMethod, GetChatMemberMethod, SetMessageReactionMethod
from api.types import (
    ID,
    ChatMemberMember,
    ChatMemberOwner,
    ReactionType,
    ReactionTypeCustomEmoji,
    ReactionTypeEmoji,
    ReactionTypePaid,
)

from helpers import CapturingConnection, RespondingConnection

_USER_ALICE = {"id": 1, "is_bot": False, "first_name": "Alice"}
_USER_BOB = {"id": 2, "is_bot": False, "first_name": "Bob"}

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


@pytest.mark.parametrize("status_json,want_type", [
    ({"status": "creator", "user": _USER_ALICE, "is_anonymous": False}, ChatMemberOwner),
    ({"status": "member", "user": _USER_BOB}, ChatMemberMember),
])
def test_get_chat_member_method_call_dispatches_status_to_correct_variant(
    status_json: dict, want_type: type
) -> None:
    conn = RespondingConnection(json.dumps(status_json).encode())
    result = GetChatMemberMethod(chat_id=ID(99), user_id=status_json["user"]["id"]).call(conn)
    assert isinstance(result, want_type)


def test_get_chat_member_method_call_owner_carries_user_id_from_response() -> None:
    conn = RespondingConnection(
        json.dumps({"status": "creator", "user": _USER_ALICE, "is_anonymous": False}).encode()
    )
    result = GetChatMemberMethod(chat_id=ID(99), user_id=1).call(conn)
    assert isinstance(result, ChatMemberOwner) and result.user.id == 1


def test_get_chat_administrators_method_call_dispatches_first_element_as_owner() -> None:
    data = json.dumps([
        {"status": "creator", "user": _USER_ALICE, "is_anonymous": False},
        {"status": "member", "user": _USER_BOB},
    ]).encode()
    conn = RespondingConnection(data)
    result = GetChatAdministratorsMethod(chat_id=ID(99)).call(conn)
    assert isinstance(result[0], ChatMemberOwner)


def test_get_chat_administrators_method_call_dispatches_second_element_as_member() -> None:
    data = json.dumps([
        {"status": "creator", "user": _USER_ALICE, "is_anonymous": False},
        {"status": "member", "user": _USER_BOB},
    ]).encode()
    conn = RespondingConnection(data)
    result = GetChatAdministratorsMethod(chat_id=ID(99)).call(conn)
    assert isinstance(result[1], ChatMemberMember)
