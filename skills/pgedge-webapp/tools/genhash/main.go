// genhash prints the bcrypt hash of a single password argument to stdout.
package main

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: genhash <password>")
		os.Exit(2)
	}
	hash, err := generateHash(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	fmt.Println(hash)
}

func generateHash(password string) (string, error) {
	if password == "" {
		return "", errors.New("password is required")
	}
	h, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return "", err
	}
	return string(h), nil
}
