package encoder

type Bits struct {
	bits   []byte
	length int
}

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

func (b *Bits) Append(p []byte) {
	if b.length%8 != 0 {
		panic("fractional byte")
	}
	b.bits = append(b.bits, p...)
	b.length += 8 * len(p)
}

func (b *Bits) Write(v uint, length int) {
	for length > 0 {
		n := length
		if n > 8 {
			n = 8
		}
		if b.length%8 == 0 {
			b.bits = append(b.bits, 0)
		} else {
			m := -b.length & 7
			if n > m {
				n = m
			}
		}
		b.length += n
		sh := uint(length - n)
		b.bits[len(b.bits)-1] |= uint8(v >> sh << uint(-b.length&7))
		v -= v >> sh << sh
		length -= n
	}
}

func (b *Bits) Pad(n int) {
	if n < 0 {
		panic("qr: invalid pad size")
	}
	if n <= 4 {
		b.Write(0, n)
	} else {
		b.Write(0, 4)
		n -= 4
		n -= -b.Length() & 7
		b.Write(0, -b.Length()&7)
		pad := n / 8
		for i := 0; i < pad; i += 2 {
			b.Write(0xec, 8)
			if i+1 >= pad {
				break
			}
			b.Write(0x11, 8)
		}
	}
}

// adds error correction bytes
func (b *Bits) AddECBytes(v Version, l Level) {
	nd := v.DataBytes(l)
	if b.length < nd*8 {
		b.Pad(nd*8 - b.length)
	}
	if b.length != nd*8 {
		panic("qr: too much data")
	}

	dat := b.Bytes()
	vt := &VTable[v]
	lev := &vt.ECLevel[l]
	db := nd / lev.Blocks
	extra := nd % lev.Blocks
	codewords := make([]byte, lev.Codewords)
	rs := NewRSEncoder(GField, lev.Codewords)
	for i := 0; i < lev.Blocks; i++ {
		if i == lev.Blocks-extra {
			db++
		}
		rs.ECC(dat[:db], codewords)
		b.Append(codewords)
		dat = dat[db:]
	}

	if len(b.Bytes()) != vt.Bytes {
		panic("qr: internal error")
	}
}
