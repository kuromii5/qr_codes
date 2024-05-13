package encoder

// A Mask describes a mask that is applied to the QR
// code to avoid QR artifacts being interpreted as
// alignment and timing patterns (such as the squares
// in the corners).  Valid masks are integers from 0 to 7.
type Mask int

var maskFunc = []func(int, int) bool{
	func(i, j int) bool {
		return (i+j)%2 == 0
	},
	func(i, j int) bool {
		return i%2 == 0
	},
	func(i, j int) bool {
		return j%3 == 0
	},
	func(i, j int) bool {
		return (i+j)%3 == 0
	},
	func(i, j int) bool {
		return (i/2+j/3)%2 == 0
	},
	func(i, j int) bool {
		return i*j%2+i*j%3 == 0
	},
	func(i, j int) bool {
		return (i*j%2+i*j%3)%2 == 0
	},
	func(i, j int) bool {
		return (i*j%3+(i+j)%2)%2 == 0
	},
}

func (m Mask) Invert(y, x int) bool {
	if m < 0 {
		return false
	}
	return maskFunc[m](y, x)
}
