//go:build ignore

package main

import (
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := "password123"
	if len(os.Args) > 1 {
		password = os.Args[1]
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		os.Exit(1)
	}
	fmt.Print(string(h))
}
