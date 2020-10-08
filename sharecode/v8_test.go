package sharecode

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestNewV8FromCode(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	if shareCode.BoardUnitIDs[6][6] != 6 {
		t.Errorf("Board Unit at 6x6 is %d instead of 6", shareCode.BoardUnitIDs[6][6])
	}
}

func TestToString(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	if code != shareCode.ToString() {
		t.Errorf("Struct did not serialize back to initial share code")
	}

	if code != shareCode.ToBase64String() {
		t.Errorf("Struct did not serialize back to initial base64 share code")
	}
}

func TestPrintBytesString(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	shareCode.PrintBytesString()
}

func TestDebugPrintSizes(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	shareCode.DebugPrintSizes()
}

func TestReflectAlignments(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	shareCode.ReflectAlignments()
}

func TestUnpackUnitRanks(t *testing.T) {
	code := "8qAMAAP4BAK4BAATjJ/5uAEZuAAAgEVM0LgAAAG0AbQAACwAAAP8BDAABCRsI/wAJARcBAQAOAQUBAQAGES0QbUBHOlcBEmoBAAFIACABaBABAyAAEAEpLAIgIAAwAAAGAgEgAAWCAHUR2gB0EQkBAQRjAAVyLBAAAgABBAMGdycAdy4fAK4BAA=="
	shareCode := NewV8FromCode(code)

	if shareCode.PackedUnitRanks[1] != 8224 {
		t.Errorf("Packed unit ranks at 1 is %d instead of 8224", shareCode.PackedUnitRanks[1])
	}

	if !testEq(shareCode.PackedUnitRanks[1].UnpackUnitRanks(), []uint8{0, 2, 0, 2, 0, 0, 0, 0}) {
		t.Errorf("Unpacked unit ranks were not equal to expected values")
	}
}

func TestNewV8PackedUnitRanks(t *testing.T) {
	unitRanks := []uint8{0, 2, 0, 2, 0, 0, 0, 0}
	packedUnitRank := NewV8PackedUnitRanks(unitRanks)

	if packedUnitRank != 8224 {
		t.Errorf("Expected a packed unit rank of 8224, instead got %d", packedUnitRank)
	}

	if !testEq(packedUnitRank.UnpackUnitRanks(), unitRanks) {
		t.Errorf("Unpacked unit ranks not equal to initial unit ranks")
	}
}

func TestNewV8EquippedItem3Bytes(t *testing.T) {
	itemTable := []struct {
		itemID uint16
		bytes  [3]byte
	}{
		{itemID: 10170, bytes: [3]byte{186, 39, 0}},
		{itemID: 10171, bytes: [3]byte{187, 39, 0}},
		{itemID: 10201, bytes: [3]byte{217, 39, 0}},
	}

	for _, item := range itemTable {
		equippedItem3Bytes := NewV8EquippedItem3Bytes(V8EquippedItem{ItemID: item.itemID})

		if !bytes.Equal(equippedItem3Bytes[:], item.bytes[:]) {
			t.Errorf("Equipped item bytes incorrectly encoded")
		}
	}
}

func TestToEquippedItem(t *testing.T) {
	itemTable := []struct {
		itemID uint16
		bytes  V8EquippedItem3Bytes
	}{
		{itemID: 10170, bytes: [3]byte{186, 39, 0}},
		{itemID: 10171, bytes: [3]byte{187, 39, 0}},
		{itemID: 10201, bytes: [3]byte{217, 39, 0}},
	}

	for _, item := range itemTable {
		equippedItemID := item.bytes.ToEquippedItem().ItemID

		if equippedItemID != item.itemID {
			t.Errorf("Incorrectly deserialized item id. Got %d instead of %d", equippedItemID, item.itemID)
		}
	}
}

// --- Helper functions ---
func testEq(a, b []uint8) bool {

	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
