package utils

import (
	"fmt"
	"strconv"

	"github.com/kuromii5/qr_codes/models"
)

func BytesToUInt32(bytes []byte) uint32 {
	var result uint32
	for _, b := range bytes {
		result = (result << 1) | (uint32(b) & 1)
	}
	return result
}

// PNG returns a PNG image displaying the code.
//
// PNG uses a custom encoder tailored to QR codes.
// Its compressed size is about 2x away from optimal
func PNG(c *models.QRCode) []byte {
	var p models.PNGWriter
	return p.Encode(c)
}

// Преобразует строку битов в байтовый слайс
func StringToBytes(s string) []byte {
	var bytes []byte
	for i := 0; i < len(s); i += 8 {
		byteVal, _ := strconv.ParseUint(s[i:i+8], 2, 8)
		bytes = append(bytes, byte(byteVal))
	}
	return bytes
}

// Преобразует байтовый слайс в строку битов
func BytesToString(bytes []byte) string {
	var s string
	for _, b := range bytes {
		s += fmt.Sprintf("%08b", b)
	}
	return s
}
