package models

// A Version represents a QR version.
// The version specifies the size of the QR code:
// a QR code with version v has 4v+17 pixels on a side.
// Versions number from 1 to 40: the larger the version,
// the more information the code can store.
type Version int

// describes metadata associated with a QR code version
type vData struct {
	AlignPos    int            // position of the first alignment pattern.
	AlignStride int            // distance between alignment patterns.
	Bytes       int            // number of bytes in the data for this version.
	Pattern     int            // pattern (version information bit stream) associated with v. of the QR code.
	ECLevel     [4]ECLevelData // represents the error correction levels for this version.
}

// represents the error correction level for a QR code version.
type ECLevelData struct {
	Blocks    int // number of blocks for this error correction level.
	Codewords int // number of error correction codewords for this error correction level.
}

const MinVersion = 1
const MaxVersion = 40

var VTable = []vData{
	{},
	{100, 100, 26, 0x0, [4]ECLevelData{{1, 7}, {1, 10}, {1, 13}, {1, 17}}},          // 1
	{16, 100, 44, 0x0, [4]ECLevelData{{1, 10}, {1, 16}, {1, 22}, {1, 28}}},          // 2
	{20, 100, 70, 0x0, [4]ECLevelData{{1, 15}, {1, 26}, {2, 18}, {2, 22}}},          // 3
	{24, 100, 100, 0x0, [4]ECLevelData{{1, 20}, {2, 18}, {2, 26}, {4, 16}}},         // 4
	{28, 100, 134, 0x0, [4]ECLevelData{{1, 26}, {2, 24}, {4, 18}, {4, 22}}},         // 5
	{32, 100, 172, 0x0, [4]ECLevelData{{2, 18}, {4, 16}, {4, 24}, {4, 28}}},         // 6
	{20, 16, 196, 0x7c94, [4]ECLevelData{{2, 20}, {4, 18}, {6, 18}, {5, 26}}},       // 7
	{22, 18, 242, 0x85bc, [4]ECLevelData{{2, 24}, {4, 22}, {6, 22}, {6, 26}}},       // 8
	{24, 20, 292, 0x9a99, [4]ECLevelData{{2, 30}, {5, 22}, {8, 20}, {8, 24}}},       // 9
	{26, 22, 346, 0xa4d3, [4]ECLevelData{{4, 18}, {5, 26}, {8, 24}, {8, 28}}},       // 10
	{28, 24, 404, 0xbbf6, [4]ECLevelData{{4, 20}, {5, 30}, {8, 28}, {11, 24}}},      // 11
	{30, 26, 466, 0xc762, [4]ECLevelData{{4, 24}, {8, 22}, {10, 26}, {11, 28}}},     // 12
	{32, 28, 532, 0xd847, [4]ECLevelData{{4, 26}, {9, 22}, {12, 24}, {16, 22}}},     // 13
	{24, 20, 581, 0xe60d, [4]ECLevelData{{4, 30}, {9, 24}, {16, 20}, {16, 24}}},     // 14
	{24, 22, 655, 0xf928, [4]ECLevelData{{6, 22}, {10, 24}, {12, 30}, {18, 24}}},    // 15
	{24, 24, 733, 0x10b78, [4]ECLevelData{{6, 24}, {10, 28}, {17, 24}, {16, 30}}},   // 16
	{28, 24, 815, 0x1145d, [4]ECLevelData{{6, 28}, {11, 28}, {16, 28}, {19, 28}}},   // 17
	{28, 26, 901, 0x12a17, [4]ECLevelData{{6, 30}, {13, 26}, {18, 28}, {21, 28}}},   // 18
	{28, 28, 991, 0x13532, [4]ECLevelData{{7, 28}, {14, 26}, {21, 26}, {25, 26}}},   // 19
	{32, 28, 1085, 0x149a6, [4]ECLevelData{{8, 28}, {16, 26}, {20, 30}, {25, 28}}},  // 20
	{26, 22, 1156, 0x15683, [4]ECLevelData{{8, 28}, {17, 26}, {23, 28}, {25, 30}}},  // 21
	{24, 24, 1258, 0x168c9, [4]ECLevelData{{9, 28}, {17, 28}, {23, 30}, {34, 24}}},  // 22
	{28, 24, 1364, 0x177ec, [4]ECLevelData{{9, 30}, {18, 28}, {25, 30}, {30, 30}}},  // 23
	{26, 26, 1474, 0x18ec4, [4]ECLevelData{{10, 30}, {20, 28}, {27, 30}, {32, 30}}}, // 24
	{30, 26, 1588, 0x191e1, [4]ECLevelData{{12, 26}, {21, 28}, {29, 30}, {35, 30}}}, // 25
	{28, 28, 1706, 0x1afab, [4]ECLevelData{{12, 28}, {23, 28}, {34, 28}, {37, 30}}}, // 26
	{32, 28, 1828, 0x1b08e, [4]ECLevelData{{12, 30}, {25, 28}, {34, 30}, {40, 30}}}, // 27
	{24, 24, 1921, 0x1cc1a, [4]ECLevelData{{13, 30}, {26, 28}, {35, 30}, {42, 30}}}, // 28
	{28, 24, 2051, 0x1d33f, [4]ECLevelData{{14, 30}, {28, 28}, {38, 30}, {45, 30}}}, // 29
	{24, 26, 2185, 0x1ed75, [4]ECLevelData{{15, 30}, {29, 28}, {40, 30}, {48, 30}}}, // 30
	{28, 26, 2323, 0x1f250, [4]ECLevelData{{16, 30}, {31, 28}, {43, 30}, {51, 30}}}, // 31
	{32, 26, 2465, 0x209d5, [4]ECLevelData{{17, 30}, {33, 28}, {45, 30}, {54, 30}}}, // 32
	{28, 28, 2611, 0x216f0, [4]ECLevelData{{18, 30}, {35, 28}, {48, 30}, {57, 30}}}, // 33
	{32, 28, 2761, 0x228ba, [4]ECLevelData{{19, 30}, {37, 28}, {51, 30}, {60, 30}}}, // 34
	{28, 24, 2876, 0x2379f, [4]ECLevelData{{19, 30}, {38, 28}, {53, 30}, {63, 30}}}, // 35
	{22, 26, 3034, 0x24b0b, [4]ECLevelData{{20, 30}, {40, 28}, {56, 30}, {66, 30}}}, // 36
	{26, 26, 3196, 0x2542e, [4]ECLevelData{{21, 30}, {43, 28}, {59, 30}, {70, 30}}}, // 37
	{30, 26, 3362, 0x26a64, [4]ECLevelData{{22, 30}, {45, 28}, {62, 30}, {74, 30}}}, // 38
	{24, 28, 3532, 0x27541, [4]ECLevelData{{24, 30}, {47, 28}, {65, 30}, {77, 30}}}, // 39
	{28, 28, 3706, 0x28c69, [4]ECLevelData{{25, 30}, {49, 28}, {68, 30}, {81, 30}}}, // 40
}

func (v Version) SizeClass() int {
	if v <= 9 {
		return 0
	}
	if v <= 26 {
		return 1
	}
	return 2
}

// DataBytes returns the number of data bytes that can be
// stored in a QR code with the given version and level.
func (v Version) DataBytes(l Level) int {
	vData := &VTable[v]
	level := &vData.ECLevel[l]
	return vData.Bytes - level.Blocks*level.Codewords
}

func (v Version) GetPixelSize() int {
	return (4 * int(v)) + 17
}
