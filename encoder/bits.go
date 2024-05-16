package encoder

import (
	"fmt"

	"github.com/kuromii5/qr_codes/models"
)

type Bits struct {
	bits   []byte
	length int
}

func (b *Bits) Print() {
	paddedLength := b.length + (8-(b.length%8))%8
	for i := 0; i < paddedLength; i++ {
		if b.bits[i/8]&(1<<(7-i%8)) != 0 {
			fmt.Print("1")
		} else {
			fmt.Print("0")
		}
		if (i+1)%8 == 0 {
			fmt.Print(" ")
		}
	}
	fmt.Println()
}

// padding bytes
var padBytes = []uint{0xec, 0x11}

func (b *Bits) Reset() {
	b.bits = b.bits[:0]
	b.length = 0
}

func (b *Bits) Length() int {
	return b.length
}

func (b *Bits) Bytes() []byte {
	if b.length%8 != 0 {
		panic("fractional byte")
	}
	return b.bits
}

// appends byte string to the bit stream
func (b *Bits) Append(data []byte) {
	if b.length%8 != 0 {
		panic("fractional byte")
	}
	b.bits = append(b.bits, data...)
	b.length += 8 * len(data)
}

// writes value of particular length to the bit stream
func (b *Bits) Write(value uint, length int) {
	for length > 0 {
		bitsToWrite := length
		if bitsToWrite > 8 {
			bitsToWrite = 8
		}
		if b.length%8 == 0 {
			b.bits = append(b.bits, 0)
		} else {
			available := -b.length & 7
			if bitsToWrite > available {
				bitsToWrite = available
			}
		}
		b.length += bitsToWrite
		shift := uint(length - bitsToWrite)
		b.bits[len(b.bits)-1] |= uint8(value >> shift << uint(-b.length&7))
		value -= value >> shift << shift
		length -= bitsToWrite
	}
}

// adds padding to bit stream until target length (to max capacity for given version)
func (b *Bits) Pad(targetLength int) {
	if targetLength <= 4 {
		b.Write(0, targetLength)
	} else {
		b.Write(0, 4)
		targetLength -= 4
		targetLength -= -b.Length() & 7
		b.Write(0, -b.Length()&7)
		remainingBytes := targetLength / 8
		for i := 0; i < remainingBytes; i++ {
			b.Write(padBytes[i%2], 8)
		}
	}
}

// adds error correction bytes
func (b *Bits) AddECBytes(v models.Version, l models.Level) {
	// pad bytes if data bytes required for current
	// version are not fullfilled by encoded data
	bytes := v.DataBytes(l)
	if b.length < bytes*8 {
		b.Pad(bytes*8 - b.length)
	}

	data := b.Bytes()
	vInfo := &models.VTable[v]
	level := &vInfo.ECLevel[l]
	bytesPerBlock := bytes / level.Blocks
	extra := bytes % level.Blocks
	codewords := make([]byte, level.Codewords)
	encoder := models.NewRSEncoder(models.GField, level.Codewords)
	for i := 0; i < level.Blocks; i++ {
		if i == level.Blocks-extra {
			bytesPerBlock++
		}
		encoder.AddECC(data[:bytesPerBlock], codewords)
		b.Append(codewords)
		data = data[bytesPerBlock:]
	}

	if len(b.Bytes()) != vInfo.Bytes {
		panic("qr: internal error")
	}
}
