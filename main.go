package main

import (
	"log"

	"github.com/ThomasK33/go-underlords/lib"
)

// Runtime version and build number
var (
	Version string
	Build   string
)

func main() {
	log.Println("Starting " + Version + " (" + Build + ")")

	code := lib.ShareCodeV8{}
	// code.DebugPrintSizes()

	code.BoardUnitIDs[0][0] = 46 // Alchemist
	code.BoardUnitIDs[6][6] = 11 // Antimage
	code.BoardUnitIDs[6][3] = 52 // Lich

	code.BoardUnitIDs[4][4] = 255 // Underlord unit
	code.UnderlordIDs[1] = 4      // Hobgen
	code.UnderlordRanks[1] = 4    // Underlords rank

	code.UnitItems[0][0] = lib.NewEquippedItem3Bytes(lib.EquippedItem{
		ItemID: 10171,
	})

	code.UnequippedItems[0][0] = lib.NewEquippedItem3Bytes(lib.EquippedItem{
		ItemID: 10170,
	})

	code.PackedUnitRanks[0] = lib.PackUnitRanks([]uint8{2})
	code.PackedUnitRanks[6] = lib.PackUnitRanks([]uint8{0, 0, 0, 0, 0, 0, 3, 0})

	successfullShareCode := code.ToBase64String()
	log.Println(successfullShareCode)

	// testBoardCode := successfullShareCode[1:]
	testBoardCode := "8qAMMALsnAP4BAP4BAPIBAAAuir4AAP82JAAMNAAACxkSCEcAVxkNUgEAAAI2FgAMAQAABhESAAMRCYIBABQEAAa6JwAVA74BAA=="
	newShareCode := lib.ShareCodeFromBase64(testBoardCode)
	log.Printf("Unit at 6x3: %d", newShareCode.BoardUnitIDs[6][3])
}
