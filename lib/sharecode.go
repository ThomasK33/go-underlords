package lib

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"unsafe"

	"github.com/golang/snappy"
)

// Share Code Constants
const (
	ShareCodeVersion            = 8
	BoardCellNum                = 8
	SharecodeMaxTalents         = 16
	SharecodeMaxUnequippedItems = 10
)

// ShareCodeV8 - The actual share code structure
type ShareCodeV8 struct {
	// total sizeof: 424

	// Intentionally not public, as this field shall be a constant 0
	version              uint8
	UnitItems            [BoardCellNum][BoardCellNum]EquippedItem3Bytes     // 24 bits per unit item    (192 bytes)
	BoardUnitIDs         [BoardCellNum][BoardCellNum]uint8                  // 8 bits per unit          (64 bytes)
	SelectedTalents      [SharecodeMaxTalents][2]uint8                      // 128 bits per player      (32 bytes)
	PackedUnitRanks      [BoardCellNum]PackedUnitRank                       // 4 bits per unit rank     (32 bytes)
	BenchUnitItems       [BoardCellNum]EquippedItem3Bytes                   // 24 bits per unit item    (24 bytes)
	BenchedUnitIDs       [BoardCellNum]uint8                                // 8 bits per unit          (8 bytes)
	PackedBenchUnitRanks PackedUnitRank                                     // 4 bits per unit rank     (4 bytes)
	UnderlordIDs         [2]uint8                                           // 8 bits per player        (2 bytes)
	UnderlordRanks       [2]uint8                                           // 8 bits per player        (2 bytes)
	UnequippedItems      [SharecodeMaxUnequippedItems][2]EquippedItem3Bytes // 24 bits per unused item  (60 bytes)
}

// ToBase64String - Returns base64 encoding of the share code
func (sc *ShareCodeV8) ToBase64String() string {
	const sz = int(unsafe.Sizeof(*sc))
	var asByteSlice []byte = (*(*[sz]byte)(unsafe.Pointer(sc)))[:]

	compressed := snappy.Encode(nil, asByteSlice)
	encodedBoardCode := base64.StdEncoding.EncodeToString(compressed)

	successfullShareCode := strconv.FormatInt(ShareCodeVersion, 16) + encodedBoardCode

	// Assert correctness by creating new one from string
	return successfullShareCode
}

// PrintBytesString - Print bytes of share code
func (sc *ShareCodeV8) PrintBytesString() {
	const sz = int(unsafe.Sizeof(*sc))
	var asByteSlice []byte = (*(*[sz]byte)(unsafe.Pointer(sc)))[:]

	fmt.Println()
	fmt.Printf("% x", asByteSlice)
	fmt.Println()
	fmt.Println()

	for _, byteEntry := range asByteSlice {
		fmt.Print(int(byteEntry))
		fmt.Print(" ")
	}
	fmt.Println()
	fmt.Println()
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

// ReflectAlignments - Debug function to view memory usage and layout
func (sc *ShareCodeV8) ReflectAlignments() {
	// First ask Go to give us some information about the MyData type
	typ := reflect.TypeOf(*sc)
	log.Println()
	log.Printf("Struct is %d bytes long\n", typ.Size())
	// We can run through the fields in the structure in order
	n := typ.NumField()
	for i := 0; i < n; i++ {
		field := typ.Field(i)
		log.Printf("%s at offset %v, size=%d, align=%d\n",
			field.Name, field.Offset, field.Type.Size(),
			field.Type.Align())
	}

	log.Println()
}

// ShareCodeFromBase64 - Create a new share code from a byte64 string
func ShareCodeFromBase64(sBase64 string) ShareCodeV8 {
	if sBase64[0] == '8' {
		sBase64 = sBase64[1:]
	}

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

// PackedUnitRank - Alias for uint32
type PackedUnitRank uint32

// UnpackUnitRanks - Unpack packed unit ranks
func (packedRanks *PackedUnitRank) UnpackUnitRanks() []uint8 {
	var ranks []uint8 = make([]uint8, BoardCellNum)

	for i := range ranks {
		for offsetIndex, offset := range []int{0, 1, 2, 4} {
			var value uint32 = (1 << ((i * 4) + offsetIndex))
			packedRankAnd := uint32(*packedRanks) & (value)

			if packedRankAnd == value {
				ranks[i] |= (1 << offset)
			}
		}
	}

	return ranks
}

// PackUnitRanks - Function to pack uint8 array into a uint32, removing the first 4 bits of each uint8
func PackUnitRanks(ranks []uint8) PackedUnitRank {
	var packedUnitRank PackedUnitRank = 0

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
