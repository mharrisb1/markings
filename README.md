# Markings

Automatically enforce and update source code banner markings.

## What are banner markings?

Banner markings are standardized text blocks injected at the top or bottom of source code files. They are commonly used for compliance, such as applying [classification levels](https://www.dodcui.mil/Training/Banner-Line) in defense contexts or managing open-source licensing headers.

`markings` uses special `markings:managed` boundaries to safely identify, apply, and update these banners without breaking your code.

```go
// markings:managed
//
// Copyright (c) 2026 Markings Contributors
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.
//
// markings:managed

package main
```

## Installation

Currently, `markings` requires a local Go toolchain to install:

```bash
git clone https://github.com/mharrisb1/markings.git
cd markings
go install .
```

## Configuration

Create a `.markings.yaml` file in your repository root to define templates and map them to file types:

> [!WARNING]
> You must use [`migrate-marker`](https://github.com/mharrisb1/markings#usage) if you set a new value for `marker` in config for an existing project.

```yaml
marker: "my-custom-marker:managed" # Optional: defaults to "markings:managed"

templates:
  mit_header: |
    Copyright (c) {{ .Year }} {{ .Author }}
    SPDX-License-Identifier: MIT

data:
  Year: 2026
  Author: Your Name

comment_styles:
  go_line:
    prefix: "// "

rules:
  - id: go-files
    match: "**/*.go"
    header:
      template: mit_header
      newline_after: true
    comment_style: go_line
```

## Usage

You must explicitly pass the files you want to check or fix. Files that do not match any `rules` in your config are ignored.

```bash
# Check if files have correct markings (exits non-zero if missing/incorrect)
markings check main.go utils/*.go

# Automatically inject or update markings in files
markings fix main.go utils/*.go

# Migrate files from an old marker to the new one defined in your config
markings migrate-marker "markings:managed" main.go utils/*.go
```

## Pre-Commit

To automatically enforce or update markings during commits, add the following to your `.pre-commit-config.yaml`.

_(Note: This currently requires Go to be installed on the user's system)._

```yaml
repos:
  - repo: https://github.com/mharrisb1/markings
    rev: v0.1.1
    hooks:
      # Automatically fix markings
      - id: markings
        types: [text]

      # Or just check without modifying
      # - id: markings-check
      #   types: [text]
```
