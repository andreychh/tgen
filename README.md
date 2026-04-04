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
the [Telegram Bot API HTML documentation](https://core.telegram.org/bots/api).

Instead of relying on manually updated boilerplate, `tgen` parses the specification to generate
strongly-typed client code.

> [!WARNING]
> **Work in Progress:** This project is in early development. The generated code is currently
> incomplete and may not be fully functional. The API is subject to breaking changes.

## Features

* **Up-to-Date by Default:** New API methods and standard types are picked up automatically — just
  run `tgen go` after a new Telegram Bot API release.
* **Spec-Faithful Types:** Ambiguous spec types become real Go types. The Telegram API describes
  `chat_id` as `Integer or String`. Instead of collapsing this into `any`, tgen generates a proper
  union type with explicit variants:

    ```go
    // address a public channel by username
    ChatID: api.ChatID{Username: new("@news")},
    // or a group by its numeric ID — same field, different variant
    ChatID: api.ChatID{ID: new(int64(-1001122334455))},
    ```

* **Deterministic Builds:** Supports local HTML files for reproducibility or offline work.

## Installation

### Using [mise-en-place](https://mise.jdx.dev/)

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
the [Releases page](https://github.com/andreychh/tgen/releases).

## Usage

`tgen` uses subcommands to target specific languages — currently only `go` is available.

### Fetch from the web

Downloads and parses the specification directly from the Telegram website, writing the generated
files to `./api`:

```bash
tgen go --out ./api
```

### Use a local file

If you have downloaded the HTML specification locally, pass the file path using the `--spec` or `-s`
flag. This is recommended to ensure build reproducibility and avoid network issues:

```bash
# Download the specification
curl -o api.html https://core.telegram.org/bots/api

# Generate the code
tgen go -s ./api.html -o ./api
```

## Generated API

### Go

The generated API follows a consistent pattern: each Telegram Bot API method is a struct you
populate with parameters and call directly on a `Connection` to get a typed response. You can
replace `HTTPConnection` with any implementation to add retries, proxy requests, or inject
`FakeConnection` in tests:

```go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"awesome-bot/api"
)

func main() {
	conn := api.NewHTTPConnection(http.DefaultClient, os.Getenv("BOT_TOKEN"))
	ctx := context.Background()

	msg, err := api.SendMessageMethod{
		ChatID: api.ChatID{Username: new("@news")},
		Text:   "Hello from tgen!",
	}.Call(ctx, conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Sent message %d\n", msg.MessageID)
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
		ChatID: api.ChatID{Username: new("@news")},
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
		ChatID: api.ChatID{Username: new("@news")},
		Text:   "Hello from tgen!",
	}.Call(context.Background(), conn)

	assert.ErrorContains(t, err, "unauthorized")
}
```

## Contributing

Contributions are welcome! As the project evolves, help with refining the HTML parser and generation
templates is highly appreciated. Whether you're fixing a bug, improving documentation, or adding a
new feature — feel free to open a pull request.

### Getting started

1. Install [mise-en-place](https://mise.jdx.dev/) — the toolchain manager we use.

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
   please [open an issue](https://github.com/andreychh/tgen/issues/new).

> [!TIP]
> [Renovate](https://www.mend.io/renovate/) automatically keeps all dependencies up to date. Once
> its PRs are merged, run the following to update your local toolchain:
>
> ```bash
> mise upgrade
> ```
