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
	fmt.Printf("Please enter the quest account's secret key: ")
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

	// Fund and create the quest account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + questKp.Address())
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

	// Build an account creation operation to create the random account first.
	createOp := txnbuild.CreateAccount{
		Destination: generatedKp.Address(),
		Amount:      "2000",
	}

	// Build a set options operation to allow clawbacks from the issuer.
	allowOp := txnbuild.SetOptions{
		// Note that the revocable flag must be set to enable clawbacks.
		SetFlags: []txnbuild.AccountFlag{txnbuild.AuthRevocable, txnbuild.AuthClawbackEnabled},
	}

	// Create the asset
	asset := txnbuild.CreditAsset{
		Code:   "ASSET",
		Issuer: questKp.Address(),
	}

	// Build a change trust operation to allow assets from the issuer.
	trustOp := txnbuild.ChangeTrust{
		Line:          asset.MustToChangeTrustAsset(),
		SourceAccount: generatedKp.Address(),
	}

	// Build a payment operation.
	paymentOp := txnbuild.Payment{
		Destination: generatedKp.Address(),
		Amount:      "100",
		Asset:       asset,
	}

	// Build a clawback operation.
	clawbackOp := txnbuild.Clawback{
		From:   generatedKp.Address(),
		Amount: "50",
		Asset:  asset,
	}

	// Construct the transaction from the issuer account.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&createOp, &allowOp, &trustOp, &paymentOp, &clawbackOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction.
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

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
