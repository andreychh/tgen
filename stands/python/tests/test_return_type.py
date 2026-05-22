# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import json

import pytest

from api.methods import (
    EditMessageTextMethod,
    GetChatAdministratorsMethod,
    GetChatMemberMethod,
)
from api.types import ID, ChatMemberMember, ChatMemberOwner, Message

from helpers import RespondingConnection

_USER_ALICE = {"id": 1, "is_bot": False, "first_name": "Alice"}
_USER_BOB = {"id": 2, "is_bot": False, "first_name": "Bob"}


def test_edit_message_text_method_call_deserialises_object_response_as_message() -> None:
    conn = RespondingConnection(
        b'{"message_id": 42, "date": 0, "chat": {"id": 1, "type": "private"}}'
    )
    result = EditMessageTextMethod(text="hi").call(conn)
    assert isinstance(result, Message), "object response must be deserialised as Message, not True"
    assert result.message_id == 42, "Message must carry the message_id from the response"


def test_edit_message_text_method_call_deserialises_boolean_response_as_literal_true() -> None:
    conn = RespondingConnection(b"true")
    result = EditMessageTextMethod(text="hi").call(conn)
    assert result is True, "boolean true response must be deserialised as True, not Message"


@pytest.mark.parametrize("status_json,user,want_type,want_user_id", [
    (
        {"status": "creator", "user": _USER_ALICE, "is_anonymous": False},
        _USER_ALICE,
        ChatMemberOwner,
        1,
    ),
    (
        {"status": "member", "user": _USER_BOB},
        _USER_BOB,
        ChatMemberMember,
        None,
    ),
])
def test_get_chat_member_method_call_dispatches_status_to_correct_variant(
    status_json: dict, user: dict, want_type: type, want_user_id: int | None
) -> None:
    conn = RespondingConnection(json.dumps(status_json).encode())
    result = GetChatMemberMethod(chat_id=ID(99), user_id=user["id"]).call(conn)
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
