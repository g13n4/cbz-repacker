# CBZ Chapter Aggregator & Repacker

A simple utility for aggregation and consolidation cbz files in one

---

## Features

* **Custom cbz volume size:** You can choose how many cbz files will be packed in one
* **Custom regex expression support:** Warns you if your regex doesn't match a cbz file

---

| Flag                  | Type     | Default Value               | Description                                                                                 |
|:----------------------|:---------|:----------------------------|:--------------------------------------------------------------------------------------------|
| `-input`              | `string` | `.`                         | Target directory where the source `.cbz` files are located.                                 |
| `-output`             | `string` | `.`                         | Target directory where your new repacked `.cbz` volumes will be generated.                  |
| `-size`               | `int`    | `7`                         | Maximum number of chapters to compress inside a single combined `.cbz` archive.             |
| `-pattern`            | `string` | `[Cc]hapter[_\-\s\t]*(\d+)` | Regex rule used to find and extract the chapter integer identifier.                         |
| `-ignore_not_matched` | `bool`   | `false`                     | If set to `true`, suppresses console warnings for files that don't match the regex pattern. |
| `-ignore_order_check` | `bool`   | `false`                     | If set to `true`, suppresses error on missing or repeating files.                           |

---

## Installation

Ensure you have [Go](https://go.dev/) installed on your system, then clone the repository and build the binary:

```bash
git clone [https://github.com/g13n4/cbz-repacker.git](https://github.com/g13n4/cbz-repacker.git)
cd cbz-repacker
go build -o cbz-repacker main.go