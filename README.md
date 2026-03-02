<!--
SPDX-FileCopyrightText: 2026 Andrey Chernykh
SPDX-License-Identifier: MIT
-->

# tgen

[![REUSE status](https://api.reuse.software/badge/github.com/andreychh/tgen)](https://api.reuse.software/info/github.com/andreychh/tgen)
[![codecov](https://codecov.io/gh/andreychh/tgen/graph/badge.svg?token=CRAB598PR3)](https://codecov.io/gh/andreychh/tgen)
[![GitHub release (latest by date)](https://img.shields.io/github/v/release/andreychh/tgen)](https://github.com/andreychh/tgen/releases)
[![PDD status](https://www.0pdd.com/svg?name=andreychh/tgen)](https://www.0pdd.com/p?name=andreychh/tgen)
[![Go Report Card](https://goreportcard.com/badge/github.com/andreychh/tgen)](https://goreportcard.com/report/github.com/andreychh/tgen)

**tgen** is a command-line tool that turns
the [Telegram Bot API HTML documentation](https://core.telegram.org/bots/api) into ready-to-use API
bindings.

Instead of relying on manually updated boilerplate, `tgen` parses the specification to generate
strongly-typed client code.

> [!WARNING]
> **Work in Progress:** This project is in early development. The generated code is currently
> incomplete and may not be fully functional. The API is subject to breaking changes.

## Features

* **Automated & Up-to-Date:** Parses documentation directly from the Telegram website to ensure your
  bindings match the latest API changes.
* **Fluent Output:** Aims to generate developer-friendly interfaces that follow language-specific
  idioms. For example, in Go, it targets syntax like:

    ```go
    resp, err := api.SendMessage(client).Call(api.SendMessageRequest{
        ChatID: 123456789,
        Text:   "Hello from generated code!",
    })
    ```

* **Deterministic Builds:** Supports local HTML files for reproducibility or offline work.
* **Standards-Driven:** Built following [clig.dev](https://clig.dev/) guidelines and strictly
  adheres to [REUSE](https://reuse.software/) compliance.

## Installation

### Using Go

If you have Go installed, you can install the latest version directly:

```bash
go install github.com/andreychh/tgen@latest
```

### Pre-built Binaries

You can download pre-compiled binaries for your operating system (Linux, macOS, Windows) from
the [Releases page](https://github.com/andreychh/tgen/releases).

## Usage

`tgen` uses subcommands to target specific programming languages. Currently, only the `go` command
is available.

### Fetch from the web

Download and parse the specification directly from the Telegram website, outputting the generated
files to the `./api` directory:

```bash
tgen go --out ./api
```

### Use a local file

If you have downloaded the HTML specification locally, pass the file path using the `--spec` flag.
This is recommended to ensure build reproducibility and avoid network issues:

```bash
# Download the specification
curl -o api.html https://core.telegram.org/bots/api

# Generate the code
tgen go --spec ./api.html --out ./api
```

## Contributing

Contributions are welcome! As the project evolves, help with refining the HTML parser and generation
templates is highly appreciated.

We use [Task](https://taskfile.dev/) for development automation.

Running linting tasks locally requires the following tools: `golangci-lint`, `yamllint`,
`markdownlint-cli2`, `reuse`, `typos`, and `pdd`. Note that all checks are also automatically
enforced by GitHub Actions on every Pull Request.

1. Clone the repository.
2. List available tasks:

    ```bash
    task --list
    ```

3. Run linters and tests before submitting a PR:

    ```bash
    task lint test
    ```
