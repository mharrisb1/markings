# // Markings

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

Download the latest pre-compiled binary for your OS and architecture from the [GitHub Releases](https://github.com/mharrisb1/markings/releases) page.

Alternatively, if you have a local Go toolchain installed, you can build from source:

```bash
go install github.com/mharrisb1/markings@latest
```

## Configuration

Create a `.markings.yaml` file in your repository root to define templates and map them to file types:

> [!WARNING]
> You must use [`migrate-marker`](https://github.com/mharrisb1/markings#usage) if you set a new value for `marker` in config for an existing project.

```yaml
marker: "my-custom-marker:managed" # Optional: defaults to "markings:managed"

exclude:
  - "**/.git"
  - "**/node_modules"

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

## Template Functions

You can use built-in macros and logic inside your templates.

```yaml
templates:
  dynamic_header: |
    File: {{ filename }}
    Copyright (c) {{ year }} {{ upper .Company }}
```

### Custom Macros

| Function                          | Description                     | Example                           |
| --------------------------------- | ------------------------------- | --------------------------------- |
| `upper`                           | Returns uppercase string        | `{{ upper .Name }}`               |
| `lower`                           | Returns lowercase string        | `{{ lower .Name }}`               |
| `title`                           | Returns titlecase string        | `{{ title .Name }}`               |
| `trim` / `trimLeft` / `trimRight` | Removes whitespace              | `{{ trim .Name }}`                |
| `trimPrefix` / `trimSuffix`       | Removes prefix or suffix        | `{{ trimPrefix .Name "v" }}`      |
| `replace`                         | Replaces substring              | `{{ replace .Name "old" "new" }}` |
| `contains`                        | Checks if substring exists      | `{{ if contains .Name "abc" }}`   |
| `hasPrefix` / `hasSuffix`         | Checks string boundaries        | `{{ if hasPrefix .Name "test" }}` |
| `split`                           | Splits string by separator      | `{{ split .Name "," }}`           |
| `join`                            | Joins array into string         | `{{ join .List "," }}`            |
| `time`                            | Formats current time (Go style) | `{{ time "02/01/06" }}`           |
| `date`                            | Current date (YYYY-MM-DD)       | `{{ date }}`                      |
| `year`                            | Current YYYY year               | `{{ year }}`                      |
| `month`                           | Current MM month                | `{{ month }}`                     |
| `day`                             | Current DD day                  | `{{ day }}`                       |
| `filename`                        | Current filename                | `{{ filename }}`                  |

### Standard Macros

| Function             | Description                     | Example                     |
| -------------------- | ------------------------------- | --------------------------- |
| `eq` / `ne`          | Equal / Not equal               | `{{ if eq .Role "admin" }}` |
| `lt` / `le`          | Less than / Less or equal       | `{{ if lt .Age 18 }}`       |
| `gt` / `ge`          | Greater than / Greater or equal | `{{ if gt .Score 90 }}`     |
| `and` / `or` / `not` | Logical AND / OR / NOT          | `{{ if and .A .B }}`        |
| `len`                | Length of string, slice, map    | `{{ len .Users }}`          |
| `index`              | Value from map/slice by index   | `{{ index .List 0 }}`       |
| `slice`              | Slices array or string          | `{{ slice .Name 0 3 }}`     |

## Usage

You can pass specific files, directories, or nothing at all (which defaults to scanning the current directory). Files are automatically recursively discovered and filtered against your `exclude` and `include` rules.

```bash
# Check all files in the current directory recursively
markings check

# Automatically inject or update markings in a specific folder
markings fix src/

# Migrate the entire repository to a new marker
markings migrate-marker "markings:managed"
```

## CI Integrations

You can easily integrate `markings` into your CI/CD pipelines to automatically check or fix markings. By default, it looks for `.markings.yaml` in the root, but you can override this.

### GitHub Actions

We provide native GitHub Actions to run markings.

```yaml
steps:
  - uses: actions/checkout@v4

  # Optional: Explicitly install the CLI to your PATH for custom scripts
  # - uses: mharrisb1/markings/setup@v1

  # Automatically check for correct markings
  - uses: mharrisb1/markings/check@v1
    with:
      args: "--config custom/path/.markings.yaml" # Optional arguments


  # Or automatically fix markings
  # - uses: mharrisb1/markings/fix@v1
```

### Pre-Commit

To automatically enforce or update markings during commits, add the following to your `.pre-commit-config.yaml`.

```yaml
repos:
  - repo: https://github.com/mharrisb1/markings
    rev: v0.1.1
    hooks:
      # Automatically fix markings
      - id: markings
        types: [text]
        args: ["--config", "custom/path/.markings.yaml"] # Optional arguments


      # Or just check without modifying
      # - id: markings-check
      #   types: [text]
```
