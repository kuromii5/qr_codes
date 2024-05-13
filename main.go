package main

import (
	"fmt"

	"github.com/kuromii5/qr_codes/encoder"
)

func main() {
	_, err := encoder.Encode("https://www.youtube.com/watch?v=dQw4w9WgXcQ&ab_channel=RickAstley", encoder.L)
	if err != nil {
		fmt.Println("lol")
	}
}
