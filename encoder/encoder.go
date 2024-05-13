package encoder

// Encoding implements a QR data encoding scheme.
// The implementations--Numeric, Alphanumeric, and String--specify
// the character set and the mapping from UTF-8 to code bits.
// The more restrictive the mode, the fewer code bits are needed.
type Encoding interface {
	Check() error
	Bits(v Version) int
	Encode(b *Bits, v Version)
}
