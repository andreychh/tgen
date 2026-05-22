# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
from io import BytesIO

from api.methods import SendPhotoMethod
from api.types import FileID, ID, Upload

from helpers import CapturingConnection


def test_send_photo_method_call_serialises_upload_as_file_part_at_photo_key() -> None:
    conn = CapturingConnection()
    SendPhotoMethod(
        chat_id=ID(99),
        photo=Upload(name="photo.jpg", reader=BytesIO(b"data")),
    ).call(conn)
    assert conn.files["photo"] == "photo.jpg"


def test_send_photo_method_call_serialises_file_id_as_text_part_at_photo_key() -> None:
    conn = CapturingConnection()
    SendPhotoMethod(chat_id=ID(99), photo=FileID("AgACAgIAAxk")).call(conn)
    assert conn.data["photo"] == "AgACAgIAAxk"
