# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import email.parser
import json
from typing import TypeVar
from urllib.parse import parse_qsl

import httpx
from pydantic import TypeAdapter

from api.client import Payload
from api.method import Method

T = TypeVar("T")


class CapturingConnection:
    """Materialises the Payload into an HTTP request and parses its body.

    Handles both multipart/form-data (when files are present) and
    application/x-www-form-urlencoded (text-only payloads). Returns None for
    every call — capture tests never use the return value.
    """

    def __init__(self) -> None:
        self.data: dict[str, str] = {}
        self.files: dict[str, str] = {}

    def do(self, _method: Method, payload: Payload, _adapter: TypeAdapter[T]) -> T:  # type: ignore[return]
        req = payload.request("POST", "http://fake")
        req.read()
        ct = req.headers.get("content-type", "")
        if ct.startswith("multipart/form-data"):
            self._parse_multipart(req, ct)
        elif ct.startswith("application/x-www-form-urlencoded"):
            self._parse_urlencoded(req)
        return None  # type: ignore[return-value]

    def _parse_multipart(self, req: httpx.Request, ct: str) -> None:
        raw = f"Content-Type: {ct}\r\n\r\n".encode() + req.content
        msg = email.parser.BytesParser().parsebytes(raw)
        parts = msg.get_payload()
        if not isinstance(parts, list):
            return
        for part in parts:
            if isinstance(part, str):
                continue
            name = part.get_param("name", header="content-disposition")
            filename = part.get_param("filename", header="content-disposition")
            if name is None:
                continue
            if filename is not None:
                self.files[name] = filename
            else:
                body = part.get_payload(decode=True)
                self.data[name] = body.decode() if isinstance(body, bytes) else ""

    def _parse_urlencoded(self, req: httpx.Request) -> None:
        for key, value in parse_qsl(req.content.decode()):
            self.data[key] = value


class RespondingConnection:
    """Injects a fixed JSON payload into the TypeAdapter on every call.

    Mirrors what HTTPConnection does after unwrapping the Telegram API envelope.
    """

    def __init__(self, data: bytes) -> None:
        self._data = data

    def do(self, _method: Method, _payload: Payload, adapter: TypeAdapter[T]) -> T:
        return adapter.validate_python(json.loads(self._data))
