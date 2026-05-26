# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import pytest
import httpx
from pydantic import TypeAdapter

from api.client import Call, Error, FakeConnection, HTTPConnection, SeqCallQueue
from api.method import Method


class _EmptyPayload:
    def request(self, method: str, url: str) -> httpx.Request:
        return httpx.Request(method, url)


def _static_transport(body: dict) -> httpx.MockTransport:
    def handler(req: httpx.Request) -> httpx.Response:
        return httpx.Response(200, json=body)
    return httpx.MockTransport(handler)


def test_error_str_returns_description() -> None:
    err = Error(400, "Bad Request: message text is empty")
    assert str(err) == "Bad Request: message text is empty", \
        "str(Error) must return the description so it integrates naturally with Python exception messages"


def test_error_stores_code_and_description() -> None:
    err = Error(401, "Unauthorized")
    assert err.code == 401 and err.description == "Unauthorized", \
        "Error must expose code and description as readable attributes"


def test_http_connection_do_decodes_result_on_ok_true_envelope() -> None:
    conn = HTTPConnection(
        httpx.Client(transport=_static_transport({"ok": True, "result": True})),
        "test-token",
    )
    result = conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))
    assert result is True, \
        "HTTPConnection must decode the result field into the return value on ok:true envelope"


@pytest.mark.parametrize("body,want_code,want_desc", [
    ({"ok": False, "error_code": 401, "description": "Unauthorized"}, 401, "Unauthorized"),
    ({"ok": False}, 0, "<no description>"),
])
def test_http_connection_do_raises_error_on_ok_false_envelope(
    body: dict, want_code: int, want_desc: str,
) -> None:
    conn = HTTPConnection(
        httpx.Client(transport=_static_transport(body)),
        "test-token",
    )
    with pytest.raises(Error) as info:
        conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))
    assert info.value.code == want_code and info.value.description == want_desc, \
        "Error must carry error_code and description from the response envelope, defaulting to 0 and '<no description>' when absent"


def test_http_connection_do_propagates_transport_error() -> None:
    def handler(req: httpx.Request) -> httpx.Response:
        raise httpx.ConnectError("refused")
    conn = HTTPConnection(httpx.Client(transport=httpx.MockTransport(handler)), "test-token")
    with pytest.raises(httpx.ConnectError):
        conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))


def test_error_raised_by_do_is_catchable_as_error_type() -> None:
    conn = HTTPConnection(
        httpx.Client(transport=_static_transport({"ok": False, "error_code": 403, "description": "Forbidden"})),
        "test-token",
    )
    caught = None
    try:
        conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))
    except Exception as exc:
        caught = exc
    assert isinstance(caught, Error), \
        "Error raised by do must be an api.Error instance so callers can catch it with 'except Error'"


def test_fake_connection_do_propagates_exception_response() -> None:
    queue = SeqCallQueue(Call(Method.GetChat, RuntimeError("bang")))
    conn = FakeConnection(queue)
    with pytest.raises(RuntimeError, match="bang"):
        conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))


def test_fake_connection_do_decodes_ok_value_into_return_value() -> None:
    queue = SeqCallQueue(Call(Method.GetChat, True))
    conn = FakeConnection(queue)
    result = conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))
    assert result is True, \
        "FakeConnection must decode the response value into the return type via the TypeAdapter"


def test_fake_connection_do_raises_runtime_error_when_queue_is_exhausted() -> None:
    conn = FakeConnection(SeqCallQueue())
    with pytest.raises(RuntimeError):
        conn.do(Method.GetChat, _EmptyPayload(), TypeAdapter(bool))


def test_fake_connection_do_raises_runtime_error_on_method_mismatch() -> None:
    queue = SeqCallQueue(Call(Method.GetChat, True))
    conn = FakeConnection(queue)
    with pytest.raises(RuntimeError):
        conn.do(Method.SendMessage, _EmptyPayload(), TypeAdapter(bool))


def test_seq_call_queue_next_raises_runtime_error_when_exhausted() -> None:
    queue = SeqCallQueue()
    with pytest.raises(RuntimeError):
        queue.next()