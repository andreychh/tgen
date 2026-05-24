# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import os
from io import BytesIO
from random import Random
from typing import IO

import httpx
import pytest
from api.client import Call, Connection, FakeConnection, HTTPConnection, SeqCallQueue
from api.method import Method
from api.methods import GetChatMethod, SendPhotoMethod, SetMessageReactionMethod
from api.types import (
    AcceptedGiftTypes,
    Chat,
    ChatFullInfo,
    ChatID,
    ID,
    InlineKeyboardButton,
    InlineKeyboardMarkup,
    Message,
    ReactionTypeEmoji,
    Upload,
    Username,
)


def announce_hoshi_cup(conn: Connection, chat_id: ChatID, poster: IO[bytes], rng: Random) -> None:
    info = GetChatMethod(chat_id=chat_id).call(conn)
    msg = SendPhotoMethod(
        chat_id=chat_id,
        photo=Upload(name="poster.jpg", reader=poster),
        caption="Hoshi Cup IX is open for registration.\nAll ranks welcome.",
        reply_markup=InlineKeyboardMarkup(inline_keyboard=[
            [InlineKeyboardButton(text="Register", url="https://hoshicup.ru/register")]
        ]),
    ).call(conn)
    found = [r for r in (info.available_reactions or []) if isinstance(r, ReactionTypeEmoji)]
    if not found:
        return
    SetMessageReactionMethod(
        chat_id=ID(msg.chat.id),
        message_id=msg.message_id,
        reaction=[rng.choice(found)],
    ).call(conn)


def _require_env(key: str) -> str:
    v = os.environ.get(key, "")
    if not v:
        pytest.skip(f"{key} not set")
    return v


def test_announce_hoshi_cup_makes_three_calls_in_sequence() -> None:
    queue = SeqCallQueue(
        Call(Method.GetChat, ChatFullInfo(
            id=0,
            type="supergroup",
            accent_color_id=0,
            max_reaction_count=0,
            accepted_gift_types=AcceptedGiftTypes(
                unlimited_gifts=False,
                limited_gifts=False,
                unique_gifts=False,
                premium_subscription=False,
                gifts_from_channels=False,
            ),
            available_reactions=[
                ReactionTypeEmoji(emoji="🎉"),
                ReactionTypeEmoji(emoji="👏"),
            ],
        )),
        Call(Method.SendPhoto, Message(
            message_id=101,
            date=0,
            chat=Chat(id=-1001234567890, type="supergroup"),
        )),
        Call(Method.SetMessageReaction, True),
    )
    conn = FakeConnection(queue)
    announce_hoshi_cup(conn, Username("@hoshi_cup"), BytesIO(b"poster"), Random(42))
    assert len(queue.calls()) == 3, \
        "announce_hoshi_cup must exhaust all three queued calls in sequence"


def test_announce_hoshi_cup_live() -> None:
    token = _require_env("BOT_API_TOKEN")
    chat_id = int(_require_env("CHAT_ID"))
    path = _require_env("POSTER_PATH")
    with open(path, "rb") as poster:
        conn = HTTPConnection(httpx.Client(timeout=30), token)
        announce_hoshi_cup(conn, ID(chat_id), poster, Random(0))
