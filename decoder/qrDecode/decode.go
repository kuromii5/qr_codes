package qrdecode

import (
	"fmt"

	"github.com/kuromii5/qr_codes/models"
)

type Decoder interface {
	Decode(bits *BitsReader, version models.Version) (string, error)
}

func ReadQR(qr *models.QRCode) string {
	qrMatrix := qr.CreateBitmapMatrix()

	// extract format data
	level, mask := getFormatInfo(qrMatrix)

	// extract version data
	var version models.Version // 0 by default
	if len(qrMatrix) >= 45 {
		version = getVersionInfo(qrMatrix)
	} else {
		version = models.Version((len(qrMatrix) - 17) / 4)
	}
	if version < 1 {
		return "Wrong size of QR-code"
	}

	// mark patterns so we can properly unmask data
	qrMatrix = MarkPatterns(version, qrMatrix)

	// unmask qr data
	unmaskQRCode(qrMatrix, mask)

	// read encoded data and encoding mode
	bits := readData(qrMatrix)
	bitReader := NewBitsReader(bits)
	mode := bitReader.Read(4)

	// check mode
	var dec Decoder
	switch mode {
	case 1:
		dec = Numeric{}
	case 2:
		dec = Alphanumeric{}
	case 4:
		dec = Byte{}
	default:
		return "Unsupported mode"
	}

	// decode the data with picked mode
	decoded, err := dec.Decode(bitReader, version)
	if err != nil {
		return fmt.Sprintf("Error decoding: %v", err)
	}

	// correct errors with given error correction level
	fmt.Println(level)

	// return data
	return decoded
}

func readData(qrMatrix [][]byte) []byte {
	size := len(qrMatrix)
	var dataBits []byte

	readBit := func(x, y int) byte {
		return qrMatrix[y][x] & 1
	}

	// traverse the QR matrix in zigzag pattern to read the data bits
	for x := size - 1; x > 0; x -= 2 {
		if x == 6 { // skip the vertical timing pattern
			x--
		}
		for y := size - 1; y >= 0; y-- {
			if qrMatrix[y][x-1] != 2 {
				dataBits = append(dataBits, readBit(x, y))
				dataBits = append(dataBits, readBit(x-1, y))
			}
		}
		x -= 2
		if x == 6 { // skip the vertical timing pattern
			x--
		}
		for y := 0; y < size; y++ {
			if qrMatrix[y][x-1] != 2 {
				dataBits = append(dataBits, readBit(x, y))
				dataBits = append(dataBits, readBit(x-1, y))
			}
		}
	}

	// convert bits to bytes
	var dataBytes []byte
	for i := 0; i < len(dataBits); i += 8 {
		byteVal := byte(0)
		for j := 0; j < 8 && i+j < len(dataBits); j++ {
			byteVal = (byteVal << 1) | dataBits[i+j]
		}
		dataBytes = append(dataBytes, byteVal)
	}

	return dataBytes
}

func unmaskQRCode(matrix [][]byte, mask models.Mask) {
	size := len(matrix)
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			if matrix[y][x] != 2 && mask.Invert(y, x) {
				matrix[y][x] ^= 1
			}
		}
	}
}
