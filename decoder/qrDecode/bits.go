package qrdecode

// BitsReader is a helper to read bits from the encoded string
type BitsReader struct {
	data   []byte
	offset int
}

// NewBitsReader creates a new BitsReader
func NewBitsReader(data []byte) *BitsReader {
	return &BitsReader{data: data}
}

// Read reads the specified number of bits and returns the corresponding integer value
func (br *BitsReader) Read(n int) uint {
	var value uint
	for i := 0; i < n; i++ {
		if br.offset >= len(br.data)*8 {
			break
		}
		byteIndex := br.offset / 8
		bitIndex := br.offset % 8
		bit := (br.data[byteIndex] >> (7 - bitIndex)) & 1
		value = (value << 1) | uint(bit)
		br.offset++
	}
	return value
}
