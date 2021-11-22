package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/keypair"
)

func main() {
	// Get the secret key from user input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the account from the secret key.
	kp, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Print the public key.
	fmt.Printf("Your public key is %v\n", kp.Address())

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
