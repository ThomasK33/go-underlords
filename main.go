package main

import (
	"log"

	"github.com/ThomasK33/go-underlords/sharecode"
)

// Runtime version and build number
var (
	Version string
	Build   string
)

func main() {
	log.Println("Starting " + Version + " (" + Build + ")")

	code := sharecode.V8{}
	// code.DebugPrintSizes()
	// code.ReflectAlignments()

	code.BoardUnitIDs[0][0] = 46 // Alchemist
	code.BoardUnitIDs[6][6] = 11 // Antimage
	code.BoardUnitIDs[6][3] = 52 // Lich

	code.BoardUnitIDs[4][4] = 255 // Underlord unit
	code.UnderlordIDs[1] = 4      // Hobgen
	code.UnderlordRanks[1] = 4    // Underlords rank

	code.UnitItems[0][0] = sharecode.V8NewEquippedItem3Bytes(sharecode.V8EquippedItem{
		ItemID: 10171,
	})

	code.UnequippedItems[0][0] = sharecode.V8NewEquippedItem3Bytes(sharecode.V8EquippedItem{
		ItemID: 10170,
	})

	code.PackedUnitRanks[0] = sharecode.V8PackUnitRanks([]uint8{2})
	code.PackedUnitRanks[6] = sharecode.V8PackUnitRanks([]uint8{0, 0, 0, 0, 0, 0, 3, 0})

	successfullShareCode := code.ToBase64String()
	log.Println(successfullShareCode)

	// testBoardCode := successfullShareCode
	testBoardCode := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	newShareCode := sharecode.V8FromBase64(testBoardCode)
	log.Print(newShareCode.BoardUnitIDs)
	log.Printf("Unit at 6x3: %d", newShareCode.BoardUnitIDs[6][6])
}
