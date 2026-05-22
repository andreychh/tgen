# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import pytest

from api.methods import SendMessageMethod
from api.types import ID, ChatID, Username

from helpers import CapturingConnection


@pytest.mark.parametrize("chat_id,want", [
    (ID(12345), "12345"),
    (Username("@octocat"), "@octocat"),
])
def test_send_message_method_call_serialises_chat_id_to_expected_string(
    chat_id: ChatID, want: str
) -> None:
    conn = CapturingConnection()
    SendMessageMethod(chat_id=chat_id, text="hi").call(conn)
    assert conn.data["chat_id"] == want
