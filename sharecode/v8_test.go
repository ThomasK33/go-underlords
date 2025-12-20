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
		item  V8EquippedItem
		bytes [3]byte
	}{
		{item: V8EquippedItem{ItemID: 10170}, bytes: [3]byte{186, 39, 0}},
		{item: V8EquippedItem{ItemID: 10171}, bytes: [3]byte{187, 39, 0}},
		{item: V8EquippedItem{ItemID: 10201}, bytes: [3]byte{217, 39, 0}},
		// Test gem flags
		{item: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true}, bytes: [3]byte{186, 39, 0x01}},
		{item: V8EquippedItem{ItemID: 10170, HasDefensiveGem: true}, bytes: [3]byte{186, 39, 0x02}},
		{item: V8EquippedItem{ItemID: 10170, HasSupportGem: true}, bytes: [3]byte{186, 39, 0x04}},
		{item: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true, HasDefensiveGem: true}, bytes: [3]byte{186, 39, 0x03}},
		{item: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true, HasDefensiveGem: true, HasSupportGem: true}, bytes: [3]byte{186, 39, 0x07}},
	}

	for _, item := range itemTable {
		equippedItem3Bytes := NewV8EquippedItem3Bytes(item.item)

		if !bytes.Equal(equippedItem3Bytes[:], item.bytes[:]) {
			t.Errorf("Equipped item bytes incorrectly encoded: got %v, expected %v", equippedItem3Bytes, item.bytes)
		}
	}
}

func TestToEquippedItem(t *testing.T) {
	itemTable := []struct {
		expected V8EquippedItem
		bytes    V8EquippedItem3Bytes
	}{
		{expected: V8EquippedItem{ItemID: 10170}, bytes: [3]byte{186, 39, 0}},
		{expected: V8EquippedItem{ItemID: 10171}, bytes: [3]byte{187, 39, 0}},
		{expected: V8EquippedItem{ItemID: 10201}, bytes: [3]byte{217, 39, 0}},
		// Test gem flags
		{expected: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true}, bytes: [3]byte{186, 39, 0x01}},
		{expected: V8EquippedItem{ItemID: 10170, HasDefensiveGem: true}, bytes: [3]byte{186, 39, 0x02}},
		{expected: V8EquippedItem{ItemID: 10170, HasSupportGem: true}, bytes: [3]byte{186, 39, 0x04}},
		{expected: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true, HasDefensiveGem: true}, bytes: [3]byte{186, 39, 0x03}},
		{expected: V8EquippedItem{ItemID: 10170, HasOffensiveGem: true, HasDefensiveGem: true, HasSupportGem: true}, bytes: [3]byte{186, 39, 0x07}},
	}

	for _, item := range itemTable {
		equippedItem, err := item.bytes.ToEquippedItem()
		if err != nil {
			t.Errorf("Error converting bytes to equipped item: %s", err)
		}

		if equippedItem.ItemID != item.expected.ItemID {
			t.Errorf("Incorrectly deserialized item id. Got %d instead of %d", equippedItem.ItemID, item.expected.ItemID)
		}
		if equippedItem.HasOffensiveGem != item.expected.HasOffensiveGem {
			t.Errorf("Incorrectly deserialized HasOffensiveGem. Got %v instead of %v", equippedItem.HasOffensiveGem, item.expected.HasOffensiveGem)
		}
		if equippedItem.HasDefensiveGem != item.expected.HasDefensiveGem {
			t.Errorf("Incorrectly deserialized HasDefensiveGem. Got %v instead of %v", equippedItem.HasDefensiveGem, item.expected.HasDefensiveGem)
		}
		if equippedItem.HasSupportGem != item.expected.HasSupportGem {
			t.Errorf("Incorrectly deserialized HasSupportGem. Got %v instead of %v", equippedItem.HasSupportGem, item.expected.HasSupportGem)
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
