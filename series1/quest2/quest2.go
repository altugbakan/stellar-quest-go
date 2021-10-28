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
	// Get the secret key from user input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

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

	// Fund the generated account.
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

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a payment operation.
	op := txnbuild.Payment{
		Destination: generatedKp.Address(),
		Asset:       txnbuild.NativeAsset{},
		Amount:      "10",
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

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction to the network.
	status, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
