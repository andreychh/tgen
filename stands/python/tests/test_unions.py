# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import json

from api.methods import AnswerInlineQueryMethod, EditMessageTextMethod, SendMessageMethod
from api.types import (
    ID,
    InaccessibleMessage,
    InlineQueryResultArticle,
    InlineQueryResultCachedPhoto,
    InlineQueryResultPhoto,
    InputLocationMessageContent,
    InputTextMessageContent,
    Message,
)

from helpers import CapturingConnection, RespondingConnection


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


def test_maybe_inaccessible_message_date_zero_deserialises_as_inaccessible_message() -> None:
    data = json.dumps({
        "message_id": 1, "date": 1000, "chat": {"id": 1, "type": "private"},
        "pinned_message": {"message_id": 0, "date": 0, "chat": {"id": 1, "type": "private"}},
    }).encode()
    conn = RespondingConnection(data)
    result = SendMessageMethod(chat_id=ID(1), text="hi").call(conn)
    assert isinstance(result.pinned_message, InaccessibleMessage), \
        "date:0 must be deserialised as InaccessibleMessage, not Message"


def test_maybe_inaccessible_message_nonzero_date_deserialises_as_message() -> None:
    data = json.dumps({
        "message_id": 1, "date": 1000, "chat": {"id": 1, "type": "private"},
        "pinned_message": {"message_id": 5, "date": 1234, "chat": {"id": 1, "type": "private"}},
    }).encode()
    conn = RespondingConnection(data)
    result = SendMessageMethod(chat_id=ID(1), text="hi").call(conn)
    assert isinstance(result.pinned_message, Message), \
        "non-zero date must be deserialised as Message, not InaccessibleMessage"
    assert result.pinned_message.message_id == 5, \
        "Message must carry the message_id from the response"


def test_input_message_content_text_produces_json_with_message_text_field() -> None:
    conn = CapturingConnection()
    AnswerInlineQueryMethod(
        inline_query_id="q1",
        results=[InlineQueryResultArticle(
            id="r1",
            title="title",
            input_message_content=InputTextMessageContent(message_text="hello"),
        )],
    ).call(conn)
    items = json.loads(conn.data["results"])
    assert "message_text" in items[0]["input_message_content"], \
        "InputTextMessageContent must produce JSON with a message_text field"


def test_input_message_content_location_produces_json_with_latitude_field() -> None:
    conn = CapturingConnection()
    AnswerInlineQueryMethod(
        inline_query_id="q1",
        results=[InlineQueryResultArticle(
            id="r1",
            title="title",
            input_message_content=InputLocationMessageContent(latitude=55.7558, longitude=37.6173),
        )],
    ).call(conn)
    items = json.loads(conn.data["results"])
    assert "latitude" in items[0]["input_message_content"], \
        "InputLocationMessageContent must produce JSON with a latitude field"


def test_inline_query_result_photo_emits_type_photo_with_photo_url() -> None:
    conn = CapturingConnection()
    AnswerInlineQueryMethod(
        inline_query_id="q1",
        results=[InlineQueryResultPhoto(
            id="r1",
            photo_url="https://example.com/p.jpg",
            thumbnail_url="https://example.com/t.jpg",
        )],
    ).call(conn)
    items = json.loads(conn.data["results"])
    assert items[0]["type"] == "photo", \
        "InlineQueryResultPhoto must inject type:photo discriminator"
    assert "photo_url" in items[0], \
        "InlineQueryResultPhoto must include photo_url to allow server-side dispatch within the group"


def test_inline_query_result_cached_photo_emits_type_photo_with_photo_file_id() -> None:
    conn = CapturingConnection()
    AnswerInlineQueryMethod(
        inline_query_id="q1",
        results=[InlineQueryResultCachedPhoto(id="r1", photo_file_id="AgACfile123")],
    ).call(conn)
    items = json.loads(conn.data["results"])
    assert items[0]["type"] == "photo", \
        "InlineQueryResultCachedPhoto must inject type:photo discriminator"
    assert "photo_file_id" in items[0], \
        "InlineQueryResultCachedPhoto must include photo_file_id to allow server-side dispatch within the group"
