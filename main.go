package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/kuromii5/qr_codes/encoder"
	"github.com/kuromii5/qr_codes/models"
	"github.com/kuromii5/qr_codes/utils"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	// Get input string from the user
	fmt.Print("Enter the string to encode as QR: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	input = strings.TrimSpace(input)

	// Get error correction level from the user
	fmt.Print("Enter the error correction level (L, M, Q, H): ")
	levelInput, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading input:", err)
		return
	}
	levelInput = strings.TrimSpace(levelInput)

	var level models.Level
	switch strings.ToUpper(levelInput) {
	case "L":
		level = models.L
	case "M":
		level = models.M
	case "Q":
		level = models.Q
	case "H":
		level = models.H
	default:
		fmt.Println("Invalid error correction level. Choose from L, M, Q, H.")
		return
	}

	// Encode the string into a QR code
	qrCode, err := encoder.Encode(input, level)
	if err != nil {
		fmt.Println("Error encoding QR code:", err)
		return
	}

	// Convert QR code to PNG data
	pngData := utils.PNG(qrCode)

	// Save the PNG data to a file
	fileName := "qr.png"
	err = os.WriteFile(fileName, pngData, 0666)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("QR code saved to", fileName)
}
