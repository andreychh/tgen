<!--
SPDX-FileCopyrightText: 2026 Andrey Chernykh
SPDX-License-Identifier: MIT
-->

# tgen

[![REUSE status](https://api.reuse.software/badge/github.com/andreychh/tgen)](https://api.reuse.software/info/github.com/andreychh/tgen)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/andreychh/tgen)](https://github.com/andreychh/tgen/releases)
[![PDD status](https://www.0pdd.com/svg?name=andreychh/tgen)](https://www.0pdd.com/p?name=andreychh/tgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreychh/tgen)](https://goreportcard.com/report/github.com/andreychh/tgen)
<!--
[![codecov](https://codecov.io/gh/andreychh/tgen/graph/badge.svg?token=CRAB598PR3)](https://codecov.io/gh/andreychh/tgen)
-->

**tgen** is a command-line tool that generates ready-to-use API bindings from
the [Telegram Bot API HTML documentation][telegram-api].

Instead of relying on manually updated boilerplate, tgen parses the specification to generate
strongly-typed client code in Go and Python.

> [!WARNING]
> **Work in Progress:** This project is in early development. The generated code is currently
> incomplete and may not be fully functional. The API is subject to breaking changes.

## Features

* **Up-to-Date by Default:** New API methods and standard types are picked up automatically — just
  run `tgen go` or `tgen python` after a new Telegram Bot API release.
* **Spec-Faithful Naming:** Every name from the Telegram Bot API — methods, types, and fields — is
  carried over unchanged, adapted to the target language's naming conventions (`message_id` →
  `MessageID` in Go, unchanged in Python).
* **Spec-Faithful Types:** Ambiguous spec types become real types. The Telegram API describes
  `chat_id` as `Integer or String`. Instead of collapsing this into `any`, tgen generates a proper
  union type with explicit variants:

    ```go
    // address a public channel by username
    ChatID: api.Username("@news"),
    // or a group by its numeric ID — same field, different variant
    ChatID: api.ID(-1001122334455),
    ```

* **Deterministic Builds:** Supports local HTML files for reproducibility or offline work.

## Installation

### Using [mise-en-place][mise]

If you use mise, install the latest release globally:

```bash
mise use -g github:andreychh/tgen
```

### Using Go

With Go available, you can fetch the latest version directly:

```bash
go install github.com/andreychh/tgen@latest
```

### Pre-built Binaries

You can download pre-compiled binaries for your operating system (Linux, macOS, Windows) from
the [Releases page][releases].

## Usage

tgen uses subcommands to target specific languages: `go` and `python`.

### Fetch from the web

Downloads and parses the specification directly from the Telegram website, writing the generated
files to `./api`:

```bash
tgen go --out ./api
tgen python --out ./api
```

### Use a local file

If you have downloaded the HTML specification locally, pass the file path using the `--spec` or `-s`
flag. This is recommended to ensure build reproducibility and avoid network issues:

```bash
# Download the specification
curl -o ./api.html https://core.telegram.org/bots/api

# Generate Go bindings
tgen go -s ./api.html -o ./api

# Generate Python bindings
tgen python -s ./api.html -o ./api
```

## Generated API

### Go

Each Telegram Bot API method is a struct you populate with parameters and call directly on a
`Connection` to get a typed response. You can replace `HTTPConnection` with any implementation to
add retries, proxy requests, or inject `FakeConnection` in tests.

The following script checks bot permissions, sends a release photo with an inline button, and reacts
to the sent message.

```go
package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"my-awesome-bot/api"
)

var (
	token  = "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
	chatID = -1001122334455
)

func main() {
	conn := api.NewHTTPConnection(http.DefaultClient, token)
	ctx := context.Background()

	bot, err := api.GetMeMethod{}.Call(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// ChatID accepts a numeric ID or a channel username interchangeably:
	// chat := api.Username("@mychannel")
	chat := api.ID(chatID)

	member, err := api.GetChatMemberMethod{ChatID: chat, UserID: bot.ID}.Call(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// ChatMember is a sealed union — each variant carries the fields for that role.
	switch member.(type) {
	case api.ChatMemberAdministrator, api.ChatMemberOwner:
		// allowed to post
	default:
		log.Fatalln("bot needs admin rights to post in this chat")
	}

	cover, err := os.Open("cover.jpg")
	if err != nil {
		log.Fatalln(err)
	}
	defer cover.Close()

	msg, err := api.SendPhotoMethod{
		ChatID: chat,
		// Use api.FileID("...") to reuse a photo already on Telegram servers.
		Photo:   api.Upload{Name: "cover.jpg", Reader: cover},
		Caption: new("v2.0 is out! Faster, smaller, better."),
		ReplyMarkup: api.InlineKeyboardMarkup{
			InlineKeyboard: [][]api.InlineKeyboardButton{
				{{Text: "What's new", URL: new("https://github.com/you/proj/releases")}},
			},
		},
	}.Call(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// ReactionTypeEmoji satisfies ReactionType — Type: "emoji" is set automatically.
	_, err = api.SetMessageReactionMethod{
		ChatID:    api.ID(msg.Chat.ID),
		MessageID: msg.MessageID,
		Reaction:  []api.ReactionType{api.ReactionTypeEmoji{Emoji: "🎉"}},
	}.Call(ctx, conn)
	if err != nil {
		log.Fatalln(err)
	}
}
```

#### Testing

`FakeConnection` provides canned responses for unit tests: use `api.Ok(v)` for a successful result
and `api.Err(err)` to simulate a failure.

```go
func TestSendMessage(t *testing.T) {
	conn := api.NewFakeConnection(api.Responses{
		api.MethodSendMessage: api.Ok(api.Message{MessageID: 42}),
	})

	msg, err := api.SendMessageMethod{
		ChatID: api.ID(-1001122334455),
		Text:   "Hello from tgen!",
	}.Call(context.Background(), conn)

	assert.NoError(t, err)
	assert.Equal(t, int64(42), msg.MessageID, "SendMessage must return the sent message")
}

func TestSendMessage_Failure(t *testing.T) {
	conn := api.NewFakeConnection(api.Responses{
		api.MethodSendMessage: api.Err(errors.New("unauthorized")),
	})

	_, err := api.SendMessageMethod{
		ChatID: api.ID(-1001122334455),
		Text:   "Hello from tgen!",
	}.Call(context.Background(), conn)

	assert.ErrorContains(t, err, "unauthorized")
}
```

#### Exhaustiveness checking

All union types are sealed interfaces annotated with `//sumtype:decl`. Use [go-check-sumtype]
(available in [golangci-lint]) to catch type switches that don't cover all variants.

### Python

Each Telegram Bot API method is a pydantic model you instantiate with parameters and call directly
on a `Connection` to get a typed response. You can replace `HTTPConnection` with any implementation
to add retries, proxy requests, or inject `FakeConnection` in tests.

**Dependencies:** `pip install pydantic httpx`

The following script checks bot permissions, sends a release photo with an inline button, and reacts
to the sent message.

```python
import sys

import httpx

from api import (
    ChatMemberAdministrator,
    ChatMemberOwner,
    GetChatMemberMethod,
    GetMeMethod,
    HTTPConnection,
    ID,
    InlineKeyboardButton,
    InlineKeyboardMarkup,
    ReactionTypeEmoji,
    SendPhotoMethod,
    SetMessageReactionMethod,
    Upload,
)

TOKEN = "123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11"
CHAT_ID = -1001122334455


def main():
    conn = HTTPConnection(httpx.Client(timeout=30), TOKEN)

    bot = GetMeMethod().call(conn)

    # ChatID accepts a numeric ID or a channel username interchangeably:
    # chat = Username("@mychannel")
    chat = ID(CHAT_ID)

    member = GetChatMemberMethod(chat_id=chat, user_id=bot.id).call(conn)

    # ChatMember is a discriminated union — each variant carries the fields for that role.
    match member:
        case ChatMemberAdministrator() | ChatMemberOwner():
            pass  # allowed to post
        case _:
            sys.exit("bot needs admin rights to post in this chat")

    with open("cover.jpg", "rb") as cover:
        msg = SendPhotoMethod(
            chat_id=chat,
            # Pass FileID("...") to reuse a photo already on Telegram servers.
            photo=Upload(name="cover.jpg", reader=cover),
            caption="v2.0 is out! Faster, smaller, better.",
            reply_markup=InlineKeyboardMarkup(
                inline_keyboard=[[
                    InlineKeyboardButton(
                        text="What's new →",
                        url="https://github.com/you/proj/releases",
                    )
                ]]
            ),
        ).call(conn)

    # ReactionTypeEmoji is a discriminated union — type="emoji" is set automatically.
    SetMessageReactionMethod(
        chat_id=ID(msg.chat.id),
        message_id=msg.message_id,
        reaction=[ReactionTypeEmoji(emoji="🎉")],
    ).call(conn)


if __name__ == "__main__":
    main()
```

#### Async client

An async equivalent is available in the `api.asyncio` subpackage. Swap `api.HTTPConnection` for
`api.asyncio.HTTPConnection` (backed by `httpx.AsyncClient`) and `await` each `.call()`.

```python
from api.asyncio import HTTPConnection


async def main():
    conn = HTTPConnection(httpx.AsyncClient(timeout=30), TOKEN)
    bot = await GetMeMethod().call(conn)
    ...
```

#### Testing

`FakeConnection` provides canned responses for unit tests. Pass a value directly for a successful
result or an exception instance to simulate a failure.

```python
def test_send_message():
    conn = FakeConnection({
        Method.SendMessage: Message(message_id=42),
    })

    msg = SendMessageMethod(
        chat_id=ID(-1001122334455),
        text="Hello from tgen!",
    ).call(conn)

    assert msg.message_id == 42, "SendMessage must return the sent message"


def test_send_message_failure():
    conn = FakeConnection({
        Method.SendMessage: TelegramError("unauthorized"),
    })

    with pytest.raises(TelegramError, match="unauthorized"):
        SendMessageMethod(
            chat_id=ID(-1001122334455),
            text="Hello from tgen!",
        ).call(conn)
```

## Contributing

Contributions are welcome! As the project evolves, help with refining the HTML parser and generation
templates is highly appreciated. Whether you're fixing a bug, improving documentation, or adding a
new feature — feel free to open a pull request.

### Getting started

1. Install [mise-en-place][mise] — the toolchain manager we use.

2. Set up the toolchain — installs Go, Task, and all other tools declared in `.config/mise.toml`:

   ```bash
   mise install
   ```

3. Explore available tasks:

   ```bash
   task --list
   ```

4. Before submitting a PR, make sure linters and tests pass:

   ```bash
   task lint test
   ```

   If this fails in your environment for reasons unrelated to your changes,
   please [open an issue][new-issue].

> [!TIP]
> [Renovate][renovate] automatically keeps all dependencies up to date. Once its PRs are merged, run
> the following to update your local toolchain:
>
> ```bash
> mise upgrade
> ```

[telegram-api]: https://core.telegram.org/bots/api
[new-issue]: https://github.com/andreychh/tgen/issues/new
[releases]: https://github.com/andreychh/tgen/releases
[mise]: https://mise.jdx.dev/
[renovate]: https://github.com/apps/renovate
[go-check-sumtype]: https://github.com/alecthomas/go-check-sumtype
[golangci-lint]: https://golangci-lint.run
