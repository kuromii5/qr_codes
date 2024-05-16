package models

// A QRCode is a square pixel grid.
type QRCode struct {
	Bitmap []byte // 1 is black, 0 is white
	Size   int    // number of pixels on a side
	Stride int    // number of bytes per row
	Scale  int    // number of bits per QR pixel
}

// Black returns true if the pixel at (x,y) is black.
func (c *QRCode) Black(x, y int) bool {
	return 0 <= x && x < c.Size && 0 <= y && y < c.Size &&
		c.Bitmap[y*c.Stride+x/8]&(1<<uint(7-x&7)) != 0
}

func (qr *QRCode) CreateBitmapMatrix() [][]byte {
	matrix := make([][]byte, qr.Size)
	for i := 0; i < qr.Size; i++ {
		matrix[i] = make([]byte, qr.Size)
		for j := 0; j < qr.Size; j++ {
			if qr.Bitmap[i*qr.Stride+j/8]&(1<<uint(7-j%8)) != 0 {
				matrix[i][j] = 1 // Black
			} else {
				matrix[i][j] = 0 // White
			}
		}
	}
	return matrix
}
