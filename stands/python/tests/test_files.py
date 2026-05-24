# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import json
from io import BytesIO

from api.methods import PostStoryMethod, SendMediaGroupMethod, SendPhotoMethod
from api.types import (
    ID,
    FileID,
    InputMediaPhoto,
    InputStoryContentPhoto,
    Upload,
)

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


def test_post_story_method_call_links_content_photo_attach_reference_to_uploaded_file() -> None:
    conn = CapturingConnection()
    PostStoryMethod(
        business_connection_id="bc1",
        active_period=86400,
        content=InputStoryContentPhoto(photo=Upload(name="story.jpg", reader=BytesIO(b"data"))),
    ).call(conn)
    body = json.loads(conn.data["content"])
    assert body["photo"].startswith("attach://"), \
        "content JSON photo field must be an attach reference"
    key = body["photo"].removeprefix("attach://")
    assert conn.files.get(key) == "story.jpg", \
        "file must be registered under the key referenced in the JSON"


def test_post_story_method_call_uses_file_id_verbatim_in_content_json() -> None:
    conn = CapturingConnection()
    PostStoryMethod(
        business_connection_id="bc1",
        active_period=86400,
        content=InputStoryContentPhoto(photo=FileID("AgACAgIAAxk")),
    ).call(conn)
    body = json.loads(conn.data["content"])
    assert body["photo"] == "AgACAgIAAxk", \
        "FileID must appear verbatim in content JSON without attach:// wrapper"


def test_send_media_group_method_call_links_each_photo_to_distinct_file_part_via_attach() -> None:
    conn = CapturingConnection()
    SendMediaGroupMethod(
        chat_id=ID(1),
        media=[
            InputMediaPhoto(media=Upload(name="photo1.jpg", reader=BytesIO(b"d1"))),
            InputMediaPhoto(media=Upload(name="photo2.jpg", reader=BytesIO(b"d2"))),
        ],
    ).call(conn)
    items = json.loads(conn.data["media"])
    assert items[0]["media"].startswith("attach://"), "item 0 must be an attach reference"
    assert items[1]["media"].startswith("attach://"), "item 1 must be an attach reference"
    key0 = items[0]["media"].removeprefix("attach://")
    key1 = items[1]["media"].removeprefix("attach://")
    assert key0 != key1, "each media item must reference a distinct attach key"
    assert conn.files.get(key0) == "photo1.jpg"
    assert conn.files.get(key1) == "photo2.jpg"


def test_send_media_group_method_call_uses_file_id_verbatim_in_media_json() -> None:
    conn = CapturingConnection()
    SendMediaGroupMethod(
        chat_id=ID(1),
        media=[
            InputMediaPhoto(media=FileID("AgACAgIAAxk1")),
            InputMediaPhoto(media=FileID("AgACAgIAAxk2")),
        ],
    ).call(conn)
    items = json.loads(conn.data["media"])
    assert items[0]["media"] == "AgACAgIAAxk1", \
        "FileID must appear verbatim in media JSON element without attach:// wrapper"
