# Markings

Automatically enforce source code banner markings.

## What are banner markings?

Banner markings are text lines that need to appear at either the top or bottom (or oftentimes both) of text documents.

This is common in defense contexts where a document must be marked at the top and the bottom indicating the level of classification[^1].

[^1]: https://www.dodcui.mil/Training/Banner-Line

Code is treated like any other form of text document so for compliant programs you will need to add banner markings to source code files:

```python
#                 CUI
#
#         UNCLASSIFIED//FOUO

import sys

def process_data(data):
    """Process sensitive data."""
    pass

if __name__ == "__main__":
    print("Starting process...")
    process_data([])
    sys.exit(0)

#                 CUI
#
#         UNCLASSIFIED//FOUO
```

Banner markings are also used outside of the context of defense to mark documents with license info and other metadata that may be needed for compliance.

```go
// Copyright (c) 2026 Markings Contributors
// SPDX-License-Identifier: MIT
//
// This software is released under the MIT License.

package main
```

## Installation

This tool is mainly intended to be added to CI pipelines and pre-commit hooks but can also be used as a standalone CLI. Currently, there are no pre-built binaries.

```bash
git clone https://github.com/mharrisb1/markings.git
cd markings
go install .
```

Then to see usage docs:

```bash
markings --help
```

## Integrating with Pre-Commit

To automatically add or update markings add the following to you `.pre-commit-config.yaml`:

```yaml
repos:
  - repo: https://github.com/mharrisb1/markings
    rev: v0.1.0
    hooks:
      - id: markings
```

To just check (no fix) you can add:

```yaml
repos:
  - repo: https://github.com/mharrisb1/markings
    rev: v0.1.0
    hooks:
      - id: markings-check
```

## Integrating with Github Workflows

[TODO]

## Integrating with GitLab Pipelines

[TODO]
