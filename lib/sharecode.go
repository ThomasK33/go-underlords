package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"log"
	"strconv"
	"unsafe"

	"github.com/golang/snappy"
)

// Share Code Constants
const (
	DacShareCodeVersion             = 8
	DacBoardCellNum                 = 8
	CdacSharecodeMaxTalents         = 16
	CDacSharecodeMaxUnequippedItems = 10
)

// EquippedItem - Equipped item struct
type EquippedItem struct {
	ItemID uint16
}

// EquippedItem3Bytes - Equipped items use up 24 bits, yet one bit is unused
type EquippedItem3Bytes [3]byte

// ToEquippedItem - Convert back 3 bytes array to EquippedItem struct
func (item *EquippedItem3Bytes) ToEquippedItem() EquippedItem {
	var now []byte = item[:2]
	nowBuffer := bytes.NewReader(now)
	var itemDefIndex uint16
	binary.Read(nowBuffer, binary.LittleEndian, &itemDefIndex)

	return EquippedItem{
		ItemID: itemDefIndex,
	}
}

// NewEquippedItem3Bytes - Go hack to create a 3 bytes big struct from a ShareCodeEquippedItem
func NewEquippedItem3Bytes(item EquippedItem) EquippedItem3Bytes {
	i := item.ItemID
	var h, l uint8 = uint8(i >> 8), uint8(i & 0xff)
	return EquippedItem3Bytes{l, h, 00}
}

// ShareCodeV8 - The actual share code structure
type ShareCodeV8 struct {
	// total sizeof: 424

	// Intentionally not public, as this field shall be a constant 0, padding 4 bytes
	version              uint8
	UnitItems            [DacBoardCellNum][DacBoardCellNum]EquippedItem3Bytes   // 24 bits per unit item    (192 bytes)
	BoardUnitIDs         [DacBoardCellNum][DacBoardCellNum]uint8                // 8 bits per unit          (64 bytes)
	SelectedTalents      [CdacSharecodeMaxTalents][2]uint8                      // 128 bits per player      (32 bytes)
	PackedUnitRanks      [DacBoardCellNum]uint32                                // 4 bits per unit rank     (32 bytes)
	BenchUnitItems       [DacBoardCellNum]EquippedItem3Bytes                    // 24 bits per unit item    (24 bytes)
	BenchedUnitIDs       [DacBoardCellNum]uint8                                 // 8 bits per unit          (8 bytes)
	PackedBenchUnitRanks uint32                                                 // 4 bits per unit rank     (4 bytes)
	UnderlordIDs         [2]uint8                                               // 8 bits per player        (2 bytes)
	UnderlordRanks       [2]uint8                                               // 8 bits per player        (2 bytes)
	UnequippedItems      [CDacSharecodeMaxUnequippedItems][2]EquippedItem3Bytes // 24 bits per unused item  (60 bytes)
}

// ToBase64String - Returns base64 encoding of the share code
func (sc *ShareCodeV8) ToBase64String() string {
	const sz = int(unsafe.Sizeof(*sc))
	var asByteSlice []byte = (*(*[sz]byte)(unsafe.Pointer(sc)))[:]

	compressed := snappy.Encode(nil, asByteSlice)
	encodedBoardCode := base64.StdEncoding.EncodeToString(compressed)

	successfullShareCode := strconv.FormatInt(DacShareCodeVersion, 16) + encodedBoardCode

	// Assert correctness by creating new one from string
	return successfullShareCode
}

// DebugPrintSizes - Debug tool printing byte sizes of each field
func (sc *ShareCodeV8) DebugPrintSizes() {
	log.Println()
	log.Println("---")
	log.Println("DebugPrintSizes")
	log.Println("---")

	szShareCodeEquippedItem := int(unsafe.Sizeof(EquippedItem{}))
	szShareCodeEquippedItem3Bytes := int(unsafe.Sizeof(EquippedItem3Bytes{}))
	log.Println()
	log.Println("# Equipped Item")
	log.Printf("szShareCodeEquippedItem Struct: %d bytes", szShareCodeEquippedItem)
	log.Printf("szShareCodeEquippedItem3Bytes Struct: %d bytes", szShareCodeEquippedItem3Bytes)
	log.Println()

	log.Println("# ShareCodeV8")
	szVersion := int(unsafe.Sizeof(sc.version))
	log.Printf("version: %d bytes", szVersion)
	szUnitItems := int(unsafe.Sizeof(sc.UnitItems))
	log.Printf("UnitItems: %d bytes", szUnitItems)
	szBoardUnitIDs := int(unsafe.Sizeof(sc.BoardUnitIDs))
	log.Printf("BoardUnitIDs: %d bytes", szBoardUnitIDs)
	szSelectedTalents := int(unsafe.Sizeof(sc.SelectedTalents))
	log.Printf("SelectedTalents: %d bytes", szSelectedTalents)
	szPackedUnitRanks := int(unsafe.Sizeof(sc.PackedUnitRanks))
	log.Printf("PackedUnitRanks: %d bytes", szPackedUnitRanks)
	szBenchUnitItems := int(unsafe.Sizeof(sc.BenchUnitItems))
	log.Printf("szBenchUnitItems: %d bytes", szBenchUnitItems)
	szBenchedUnitIDs := int(unsafe.Sizeof(sc.BenchedUnitIDs))
	log.Printf("BenchedUnitIDs: %d bytes", szBenchedUnitIDs)
	szPackedBenchUnitRanks := int(unsafe.Sizeof(sc.PackedBenchUnitRanks))
	log.Printf("PackedBenchUnitRanks: %d bytes", szPackedBenchUnitRanks)
	szUnderlordIDs := int(unsafe.Sizeof(sc.UnderlordIDs))
	log.Printf("UnderlordIDs: %d bytes", szUnderlordIDs)
	szUnderlordRanks := int(unsafe.Sizeof(sc.UnderlordRanks))
	log.Printf("UnderlordRanks: %d bytes", szUnderlordRanks)
	szUnequippedItems := int(unsafe.Sizeof(sc.UnequippedItems))
	log.Printf("UnequippedItems: %d bytes", szUnequippedItems)

	szTotal := int(unsafe.Sizeof(*sc))
	log.Println("---")
	log.Printf("ShareCodeV8 total size: %d bytes", szTotal)
	log.Println("---")
	log.Println()
}

// ShareCodeFromBase64 - Create a new share code from a byte64 string
func ShareCodeFromBase64(sBase64 string) ShareCodeV8 {
	decodedShareCode, _ := base64.StdEncoding.DecodeString((sBase64))
	uncompressed, _ := snappy.Decode(nil, decodedShareCode)

	newShareCode := ShareCodeV8{}

	const newSZ = int(unsafe.Sizeof(newShareCode))
	var newShareCodeByteSlice []byte = (*(*[newSZ]byte)(unsafe.Pointer(&newShareCode)))[:]

	for i, value := range uncompressed {
		newShareCodeByteSlice[i] = value
	}

	return newShareCode
}

// PackedUnitRanks - Function to pack uint8 array into a uint32, removing the first 4 bits of each uint8
func PackedUnitRanks(ranks []uint8) uint32 {
	var packedUnitRank uint32 = 0

	for i, rank := range ranks {
		for offsetIndex, offset := range []uint8{1, 2, 4, 8} {
			bit := ((rank & offset) == offset)

			if bit {
				packedUnitRank |= (1 << ((i * 4) + offsetIndex))
			}
		}
	}

	return packedUnitRank
}

// UnpackUnitRanks - Unpack packed unit ranks
func UnpackUnitRanks(packedRanks uint32) []uint8 {
	var ranks []uint8 = make([]uint8, DacBoardCellNum)

	for i := range ranks {
		for offsetIndex, offset := range []int{0, 1, 2, 4} {
			var a uint32 = (1 << ((i * 4) + offsetIndex))
			bit := packedRanks & (1 << ((i * 4) + offsetIndex))

			if bit == a {
				ranks[i] |= (1 << offset)
			}
		}
	}

	return ranks
}
