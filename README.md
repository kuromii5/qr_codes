# QR Code Encoder and Decoder

This project provides a Go-based solution for encoding and decoding QR codes. It supports multiple encoding modes (alphanumeric, byte, numeric) and different levels of error correction (L, M, Q, H).

## Features

- **Encode Strings**: Convert input strings into QR codes.
- **Decode QR Codes**: Extract and decode data from QR codes.
- **Multiple Encoding Modes**: Supports alphanumeric, byte, and numeric modes.
- **Error Correction Levels**: Supports Low (L), Medium (M), Quartile (Q), and High (H) error correction levels.
- **Save as PNG**: Save generated QR codes as PNG images.

## Installation

To use this project, you need to have Go installed on your system. You can download and install Go from [here](https://golang.org/dl/).

1. Clone the repository:
    ```bash
    git clone https://github.com/yourusername/qr-codes.git
    cd qr-codes
    ```

2. Install the necessary dependencies:
    ```bash
    go mod tidy
    ```

## Usage

### Encoding a QR Code

To encode a string into a QR code and save it as a PNG image, run the `main.go` file. You will be prompted to enter the string to encode and the desired error correction level.

```bash
go run main.go
