package main

import (
	"log"

	"github.com/ThomasK33/go-underlords/sharecode"
)

func main() {
	code := sharecode.V8{}

	code.BoardUnitIDs[0][0] = 46 // Alchemist
	code.BoardUnitIDs[6][6] = 11 // Antimage
	code.BoardUnitIDs[6][3] = 52 // Lich

	code.BoardUnitIDs[4][4] = 255 // Underlord unit
	code.UnderlordIDs[1] = 4      // Hobgen
	code.UnderlordRanks[1] = 4    // Underlords rank

	code.UnitItems[0][0] = sharecode.NewV8EquippedItem3Bytes(sharecode.V8EquippedItem{
		ItemID: 10171,
	})

	code.UnequippedItems[0][0] = sharecode.NewV8EquippedItem3Bytes(sharecode.V8EquippedItem{
		ItemID: 10170,
	})

	code.PackedUnitRanks[0] = sharecode.NewV8PackedUnitRanks([]uint8{2})
	code.PackedUnitRanks[6] = sharecode.NewV8PackedUnitRanks([]uint8{0, 0, 0, 0, 0, 0, 3, 0})

	successfullShareCode := code.ToBase64String()
	log.Println(successfullShareCode)

	testShareCode := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	newShareCode := sharecode.NewV8FromCode(testShareCode)
	log.Print(newShareCode.BoardUnitIDs)
	log.Printf("Unit at 6x6: %d", newShareCode.BoardUnitIDs[6][6])
}
