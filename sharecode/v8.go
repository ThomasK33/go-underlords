package sharecode

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

// V8 - The actual share code structure
type V8 struct {
	// total sizeof: 424

	// Intentionally not public, as this field shall be a constant 0
	version              uint8
	UnitItems            [BoardCellNum][BoardCellNum]V8EquippedItem3Bytes     // 24 bits per unit item    (192 bytes)
	BoardUnitIDs         [BoardCellNum][BoardCellNum]uint8                    // 8 bits per unit          (64 bytes)
	SelectedTalents      [SharecodeMaxTalents][2]uint8                        // 128 bits per player      (32 bytes)
	PackedUnitRanks      [BoardCellNum]V8PackedUnitRank                       // 4 bits per unit rank     (32 bytes)
	BenchUnitItems       [BoardCellNum]V8EquippedItem3Bytes                   // 24 bits per unit item    (24 bytes)
	BenchedUnitIDs       [BoardCellNum]uint8                                  // 8 bits per unit          (8 bytes)
	PackedBenchUnitRanks V8PackedUnitRank                                     // 4 bits per unit rank     (4 bytes)
	UnderlordIDs         [2]uint8                                             // 8 bits per player        (2 bytes)
	UnderlordRanks       [2]uint8                                             // 8 bits per player        (2 bytes)
	UnequippedItems      [SharecodeMaxUnequippedItems][2]V8EquippedItem3Bytes // 24 bits per unused item  (60 bytes)
}

// ToString - Alias for ToBase64String
func (sc *V8) ToString() string {
	return sc.ToBase64String()
}

// ToBase64String - Returns base64 encoding of the share code
func (sc *V8) ToBase64String() string {
	const sz = int(unsafe.Sizeof(*sc))
	var asByteSlice []byte = (*(*[sz]byte)(unsafe.Pointer(sc)))[:]

	compressed := snappy.Encode(nil, asByteSlice)
	encodedBoardCode := base64.StdEncoding.EncodeToString(compressed)

	// For a v8 sharecode the version number is always going to be 8
	successfullShareCode := strconv.FormatInt(ShareCodeVersion, 16) + encodedBoardCode

	// Assert correctness by creating new one from string
	return successfullShareCode
}

// PrintBytesString - Print bytes of share code
func (sc *V8) PrintBytesString() {
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
func (sc *V8) DebugPrintSizes() {
	log.Println()
	log.Println("---")
	log.Println("DebugPrintSizes")
	log.Println("---")

	szShareCodeEquippedItem := int(unsafe.Sizeof(V8EquippedItem{}))
	szShareCodeEquippedItem3Bytes := int(unsafe.Sizeof(V8EquippedItem3Bytes{}))
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
func (sc *V8) ReflectAlignments() {
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

// V8FromBase64 - Create a new v8 share code from a byte64 string
func V8FromBase64(sBase64 string) V8 {
	if sBase64[0] == '8' {
		sBase64 = sBase64[1:]
	}

	decodedShareCode, _ := base64.StdEncoding.DecodeString((sBase64))
	uncompressed, _ := snappy.Decode(nil, decodedShareCode)

	newShareCode := V8{}

	const newSZ = int(unsafe.Sizeof(newShareCode))
	var newShareCodeByteSlice []byte = (*(*[newSZ]byte)(unsafe.Pointer(&newShareCode)))[:]

	for i, value := range uncompressed {
		newShareCodeByteSlice[i] = value
	}

	return newShareCode
}

// V8PackedUnitRank - Alias for uint32
type V8PackedUnitRank uint32

// UnpackUnitRanks - Unpack packed unit ranks
func (packedRanks *V8PackedUnitRank) UnpackUnitRanks() []uint8 {
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

// V8PackUnitRanks - Function to pack uint8 array into a uint32, removing the first 4 bits of each uint8
func V8PackUnitRanks(ranks []uint8) V8PackedUnitRank {
	var packedUnitRank V8PackedUnitRank = 0

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

// V8EquippedItem - Equipped item struct
type V8EquippedItem struct {
	ItemID uint16
}

// V8EquippedItem3Bytes - Equipped items use up 24 bits, yet one bit is unused
type V8EquippedItem3Bytes [3]byte

// ToEquippedItem - Convert back 3 bytes array to EquippedItem struct
func (item *V8EquippedItem3Bytes) ToEquippedItem() V8EquippedItem {
	var now []byte = item[:2]
	nowBuffer := bytes.NewReader(now)
	var itemDefIndex uint16
	binary.Read(nowBuffer, binary.LittleEndian, &itemDefIndex)

	return V8EquippedItem{
		ItemID: itemDefIndex,
	}
}

// V8NewEquippedItem3Bytes - Creates a 3 bytes big struct from a ShareCodeEquippedItem
func V8NewEquippedItem3Bytes(item V8EquippedItem) V8EquippedItem3Bytes {
	i := item.ItemID
	var h, l uint8 = uint8(i >> 8), uint8(i & 0xff)
	return V8EquippedItem3Bytes{l, h, 00}
}
