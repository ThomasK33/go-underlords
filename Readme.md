# Go-Underlords

[![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/ThomasK33/go-underlords)](https://github.com/ThomasK33/go-underlords/blob/master/go.mod)
[![Go Report Card](https://goreportcard.com/badge/github.com/ThomasK33/go-underlords)](https://goreportcard.com/report/github.com/ThomasK33/go-underlords)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/ThomasK33/go-underlords/)](https://pkg.go.dev/github.com/ThomasK33/go-underlords)
[![Go CI](https://github.com/ThomasK33/go-underlords/workflows/Go%20CI/badge.svg)](https://github.com/ThomasK33/go-underlords/actions?query=workflow%3A%22Go+CI%22)
[![Coverage Status](https://coveralls.io/repos/github/ThomasK33/go-underlords/badge.svg?branch=master)](https://coveralls.io/github/ThomasK33/go-underlords?branch=master)

Go library for Dota Underlords.

---

## Features

- Share code parsing
- Share code encoding

## Installation

This is a [Go](https://golang.org/) module.

Before installing, please download and install a Go version greater than or equal to 1.11.

Installation of the module is done using:

```bash
go get github.com/ThomasK33/go-underlords
```

## Quick Start

```go
package main

import (
  "log"

  "github.com/ThomasK33/go-underlords/sharecode"
)

func main() {
  // Create a clean v8 share code object
  code := sharecode.V8{}

  code.BoardUnitIDs[0][0] = 46 // Alchemist
  code.BoardUnitIDs[6][6] = 11 // Antimage
  code.BoardUnitIDs[6][3] = 52 // Lich

  code.BoardUnitIDs[4][4] = 255 // Underlord unit
  code.UnderlordIDs[1] = 4      // Hobgen
  code.UnderlordRanks[1] = 4    // Underlords rank

  // Items are encoded in 3 bytes, thus a remapping happens here
  code.UnitItems[0][0] = sharecode.NewV8EquippedItem3Bytes(sharecode.V8EquippedItem{
    ItemID: 10171,
  })

  code.UnequippedItems[0][0] = sharecode.NewV8EquippedItem3Bytes(sharecode.V8EquippedItem{
    ItemID: 10170,
  })

  code.PackedUnitRanks[0] = sharecode.NewV8PackedUnitRanks([]uint8{2})
  code.PackedUnitRanks[6] = sharecode.NewV8PackedUnitRanks([]uint8{0, 0, 0, 0, 0, 0, 3, 0})

  // code.ToString() is equivalent to code.ToBase64String()
  successfulShareCode := code.ToBase64String()
  log.Println(successfulShareCode)

  // Instantiate a v8 share code from a base64 string
  testShareCode := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
  newShareCode := sharecode.NewV8FromCode(testShareCode)
  log.Print(newShareCode.BoardUnitIDs)
  log.Printf("Unit at 6x6: %d", newShareCode.BoardUnitIDs[6][6])
}

```

## Docs & Community

- [GoDoc](https://godoc.org/github.com/ThomasK33/go-underlords)
- [Pkg Go Dev](https://pkg.go.dev/mod/github.com/ThomasK33/go-underlords)
- [Discord Server](https://discord.gg/u9qJxzQ) (I'm Grey#1214)

## Contributing

[Contributing Guide](https://github.com/ThomasK33/go-underlords/blob/master/CONTRIBUTING.md)

## License

[MIT](https://github.com/ThomasK33/go-underlords/blob/master/LICENSE)

## Disclaimer

This project is not affiliated with Valve Corporation.
Dota Underlords, Dota and Steam are registered trademarks of Valve Corporation.
