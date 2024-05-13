package encoder

import (
	"fmt"
)

// A Template describes how to construct a QR code
// with a specific version, level, and mask.
type Template struct {
	Version Version
	Level   Level
	Mask    Mask

	DataBytes  int // number of data bytes
	CheckBytes int // number of error correcting (checksum) bytes
	Blocks     int // number of data blocks

	Grid [][]Pixel // pixel grid
}

func (template *Template) Encode(text ...Encoding) (*QRCode, error) {
	var bits Bits
	for _, t := range text {
		if err := t.Check(); err != nil {
			return nil, err
		}
		t.Encode(&bits, template.Version)
	}
	if bits.Length() > template.DataBytes*8 {
		return nil, fmt.Errorf("cannot encode %d bits into %d-bit code", bits.Length(), template.DataBytes*8)
	}
	bits.AddECBytes(template.Version, template.Level)
	bytes := bits.Bytes()

	// Now we have the error correction bytes and the data bytes.
	// Construct the actual code.
	qr := &QRCode{Size: len(template.Grid), Stride: (len(template.Grid) + 7) &^ 7}
	qr.Bitmap = make([]byte, qr.Stride*qr.Size)
	crow := qr.Bitmap
	for _, row := range template.Grid {
		for x, pix := range row {
			switch pix.Type() {
			case Data, EC:
				o := pix.Offset()
				if bytes[o/8]&(1<<uint(7-o&7)) != 0 {
					pix ^= Black
				}
			}
			if pix&Black != 0 {
				crow[x/8] |= 1 << uint(7-x&7)
			}
		}
		crow = crow[qr.Stride:]
	}
	return qr, nil
}

// NewTemplate returns a template (or "blueprint") for a QR code with the given
// version, level, and mask.
func NewTemplate(version Version, level Level, mask Mask) *Template {
	qrTemplate := createVersionTemplate(version)
	createFormatInfo(level, mask, qrTemplate)
	createData(version, level, qrTemplate)
	applyMask(mask, qrTemplate)
	return qrTemplate
}

// createVersionTemplate creates a Template for the given version.
// It draws timing patterns, finder patterns, alignment patterns
// and version patterns.
func createVersionTemplate(v Version) *Template {
	template := &Template{Version: v}

	size := 17 + int(v)*4
	grid := grid(size)
	template.Grid = grid

	// Timing patterns (overwritten by boxes).
	const tPos = 6 // timing is in row/column 6 (counting from 0)
	for i := range grid {
		p := Timing.Pixel()

		// Check if even
		if i&1 == 0 {
			p |= Black
		}

		grid[i][tPos] = p
		grid[tPos][i] = p
	}

	// Finder patterns (boxes in the top left, top-right, bottom-left corners).
	finderPattern(grid, 0, 0)
	finderPattern(grid, size-7, 0)
	finderPattern(grid, 0, size-7)

	// Alignment patterns (small boxes that help scanner to identify the grid).
	vInfo := &VTable[v]
	for x := 4; x+5 < size; {
		for y := 4; y+5 < size; {
			// don't overwrite timing markers
			if (x < 7 && y < 7) || (x < 7 && y+5 >= size-7) || (x+5 >= size-7 && y < 7) {
			} else {
				alignPattern(grid, x, y)
			}
			if y == 4 {
				y = vInfo.AlignPos
			} else {
				y += vInfo.AlignStride
			}
		}
		if x == 4 {
			x = vInfo.AlignPos
		} else {
			x += vInfo.AlignStride
		}
	}

	// Version patterns.
	vPattern := VTable[v].Pattern
	if vPattern != 0 {
		v := vPattern
		for x := 0; x < 6; x++ {
			for y := 0; y < 3; y++ {
				p := PVersion.Pixel()
				if v&1 != 0 {
					p |= Black
				}
				grid[size-11+y][x] = p
				grid[x][size-11+y] = p
				v >>= 1
			}
		}
	}

	// One lonely black pixel
	grid[size-8][8] = Unused.Pixel() | Black
	return template
}

// createFormatInfo adds the format pixels
func createFormatInfo(level Level, mask Mask, template *Template) {
	formatData := uint32(level^1) << 13 // levels: L=01, M=00, Q=11, H=10
	formatData |= uint32(mask) << 10    // apply given mask

	const formatPoly = 0x537 // polynomial G(x) = x^10 + x^8 + x^5 + x^4 + x^2 + x + 1

	// add the remainder of polynomial division to the data
	remainder := formatData
	for i := 14; i >= 10; i-- {
		if remainder&(1<<uint(i)) != 0 {
			remainder ^= formatPoly << uint(i-10)
		}
	}
	formatData |= remainder

	invert := uint32(0x5412) // mask for format info
	size := len(template.Grid)
	for i := uint(0); i < 15; i++ {
		pix := Format.Pixel() + OffsetPixel(i)

		// if data bit equals 1 make pixel black
		if (formatData>>i)&1 == 1 {
			pix |= Black
		}

		// if mask bit equals 1 then invert pixel color
		if (invert>>i)&1 == 1 {
			pix ^= White | Black
		}

		// top left
		switch {
		case i < 6:
			template.Grid[i][8] = pix
		case i < 8:
			template.Grid[i+1][8] = pix
		case i < 9:
			template.Grid[8][7] = pix
		default:
			template.Grid[8][14-i] = pix
		}
		// bottom right
		switch {
		case i < 8:
			template.Grid[8][size-1-int(i)] = pix
		default:
			template.Grid[size-1-int(14-i)][8] = pix
		}
	}
}

// createData edits a version-only template
// to add error correction level info.
func createData(v Version, level Level, template *Template) {
	template.Level = level

	codewords := VTable[v].Bytes                            // total number of codewords
	blocks := VTable[v].ECLevel[level].Blocks               // total number of EC blocks
	numberECC := VTable[v].ECLevel[level].Codewords         // total number of Error Correction Codewords
	numberDataEC := (codewords - numberECC*blocks) / blocks // total number of data error codewords per block
	extra := (codewords - numberECC*blocks) % blocks        // extra data error codewords
	dataBits := (numberDataEC*blocks + extra) * 8           // total number of data bits (or pixels)
	checkBits := numberECC * blocks * 8                     // total number of error correction bits

	template.DataBytes = codewords - numberECC*blocks
	template.CheckBytes = numberECC * blocks
	template.Blocks = blocks

	// Make data + checksum pixels.
	data := make([]Pixel, dataBits)
	for i := range data {
		data[i] = Data.Pixel() | OffsetPixel(uint(i))
	}
	check := make([]Pixel, checkBits)
	for i := range check {
		check[i] = EC.Pixel() | OffsetPixel(uint(i+dataBits))
	}

	// Split all data into 8 bit blocks.
	dataList := make([][]Pixel, blocks)
	checkList := make([][]Pixel, blocks)
	for i := 0; i < blocks; i++ {
		// The last few blocks have an extra data byte (8 pixels).
		nd := numberDataEC
		if i >= blocks-extra {
			nd++
		}
		dataList[i], data = data[0:nd*8], data[nd*8:]
		checkList[i], check = check[0:numberECC*8], check[numberECC*8:]
	}

	// Build up bit sequence, taking first byte of each block,
	// then second byte, and so on. Then checksums.
	bits := make([]Pixel, dataBits+checkBits)
	dst := bits
	for i := 0; i < numberDataEC+1; i++ {
		for _, b := range dataList {
			if i*8 < len(b) {
				copy(dst, b[i*8:(i+1)*8])
				dst = dst[8:]
			}
		}
	}
	for i := 0; i < numberECC; i++ {
		for _, b := range checkList {
			if i*8 < len(b) {
				copy(dst, b[i*8:(i+1)*8])
				dst = dst[8:]
			}
		}
	}

	// Sweep up pair of columns,
	// then down, assigning to right then left pixel.
	size := len(template.Grid)
	remPixels := make([]Pixel, 7)
	for i := range remPixels {
		remPixels[i] = Extra.Pixel()
	}
	src := append(bits, remPixels...)
	for x := size; x > 0; {
		for y := size - 1; y >= 0; y-- {
			if template.Grid[y][x-1].Type() == 0 {
				template.Grid[y][x-1], src = src[0], src[1:]
			}
			if template.Grid[y][x-2].Type() == 0 {
				template.Grid[y][x-2], src = src[0], src[1:]
			}
		}
		x -= 2
		if x == 7 { // vertical timing pattern
			x--
		}
		for y := 0; y < size; y++ {
			if template.Grid[y][x-1].Type() == 0 {
				template.Grid[y][x-1], src = src[0], src[1:]
			}
			if template.Grid[y][x-2].Type() == 0 {
				template.Grid[y][x-2], src = src[0], src[1:]
			}
		}
		x -= 2
	}
}

// applyMask edits a version+level-only template to add the mask.
func applyMask(mask Mask, template *Template) {
	template.Mask = mask
	for y, row := range template.Grid {
		for x, pix := range row {
			if pType := pix.Type(); (pType == Data || pType == EC || pType == Extra) && template.Mask.Invert(y, x) {
				row[x] ^= Black | White
			}
		}
	}
}

// finderPattern draws a box at given x, y
func finderPattern(m [][]Pixel, x, y int) {
	pos := Finder.Pixel()
	// box
	for dy := 0; dy < 7; dy++ {
		for dx := 0; dx < 7; dx++ {
			p := pos
			if dx == 0 || dx == 6 || dy == 0 || dy == 6 || 2 <= dx && dx <= 4 && 2 <= dy && dy <= 4 {
				p |= Black
			}
			m[y+dy][x+dx] = p
		}
	}

	// separator
	for dy := -1; dy < 8; dy++ {
		if 0 <= y+dy && y+dy < len(m) {
			if x > 0 {
				m[y+dy][x-1] = pos
			}
			if x+7 < len(m) {
				m[y+dy][x+7] = pos
			}
		}
	}

	for dx := -1; dx < 8; dx++ {
		if 0 <= x+dx && x+dx < len(m) {
			if y > 0 {
				m[y-1][x+dx] = pos
			}
			if y+7 < len(m) {
				m[y+7][x+dx] = pos
			}
		}
	}
}

// alignPattern draw an alignment (small) box
func alignPattern(m [][]Pixel, x, y int) {
	// box
	align := Alignment.Pixel()
	for dy := 0; dy < 5; dy++ {
		for dx := 0; dx < 5; dx++ {
			p := align
			if dx == 0 || dx == 4 || dy == 0 || dy == 4 || dx == 2 && dy == 2 {
				p |= Black
			}
			m[y+dy][x+dx] = p
		}
	}
}
