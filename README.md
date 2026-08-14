# Markings

Automatically enforce source code banner markings.

## What are banner markings?

Banner markings are text lines that need to appear at either the top or bottom (or oftentimes both) of text documents.

This is common in defense contexts where a document must be marked at the top and the bottom indicating the level of classification[^1].

[^1]: https://www.dodcui.mil/Training/Banner-Line

Code is treated like any other form of text document so for compliant programs you will need to add banner markings to source code files:

```python
#                 CUI

if __name__ == "__main__":
  ...

#                 CUI
```

Banner markings are also used outside of the context of defense to mark documents with license info and other metadata that may be needed for compliance.

## Installation

This tool is mainly intended to be added to CI pipelines and pre-commit hooks but can also be used as a standalone CLI.

[TODO]

## Integrating with Pre-Commit

[TODO]

## Integrating with Github Workflows

[TODO]

## Integrating with GitLab Pipelines

[TODO]
