package main

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	password := []byte("123456")
	hash, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error generando hash: %v\n", err)
		return
	}
	fmt.Printf("Hash generado: %s\n", hash)
}
