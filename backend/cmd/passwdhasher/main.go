package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		os.Exit(1)
	}
	hasher := sha256.New()
	hasher.Write([]byte(os.Args[1]))
	fmt.Println(hex.EncodeToString(hasher.Sum(nil)))
}
