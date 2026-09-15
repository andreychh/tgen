# SPDX-FileCopyrightText: 2026 Andrey Chernykh
# SPDX-License-Identifier: MIT
import api


def test_package_exports_every_name_it_declares() -> None:
    missing = [name for name in api.__all__ if not hasattr(api, name)]
    assert missing == []
