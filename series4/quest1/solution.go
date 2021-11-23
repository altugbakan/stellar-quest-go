package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the secret key and emoji from user input.
	var secret, emoji string
	fmt.Printf("Please enter the given secret key: ")
	fmt.Scanln(&secret)
	fmt.Printf("Please enter the given emoji: ")
	fmt.Scanln(&emoji)

	// Get the keypair of the quest account from the secret key.
	questKp, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Generate a random testnet account.
	generatedKp, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The generated secret key is %v\n", generatedKp.Seed())
	fmt.Printf("The generated public key is %v\n", generatedKp.Address())

	// Fund and create the generated account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + generatedKp.Address())
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Get and print the response from friendbot.
	if resp.Status == "200 OK" {
		fmt.Println("Successfully funded account.")
	} else {
		fmt.Println("Error funding account.")
	}

	// Fetch the account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build an account merge operation.
	op := txnbuild.AccountMerge{
		Destination: generatedKp.Address(),
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the secret key.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the emoji.
	tx, err = tx.SignHashX([]byte(emoji))
	if err != nil {
		log.Fatal(err)
	}

	// Print the XDR.
	xdr, err := tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTransaction XDR: %v\n\n", string(xdr))

	// Inform the user and wait for user input to exit.
	fmt.Printf("Account will be merged to %v, "+
		"which has the secret key %v.\nPress \"Enter\" to exit.\n",
		generatedKp.Address(), generatedKp.Seed())
	fmt.Scanln()
}
