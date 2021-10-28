package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Ask user for secret key input.
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

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a begin sponsorship operation.
	beginOp := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID: generatedKp.Address(),
	}

	// Build a create account operation.
	createOp := txnbuild.CreateAccount{
		Destination: generatedKp.Address(),
		Amount:      "0",
	}

	// Build an end sponsorship operation.
	endOp := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: generatedKp.Address(),
	}

	// Construct the transaction with multiple operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&beginOp, &createOp, &endOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with both keys.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp, generatedKp)
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

	// Inform the user and wait for user input to exit.
	fmt.Printf("The public key of the account you sponsored is \"%v\".\n", generatedKp.Address())
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
