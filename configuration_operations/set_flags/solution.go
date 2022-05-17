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

	// Fund and create the random testnet account.
	resp, err = http.Get("https://friendbot.stellar.org/?addr=" + generatedKp.Address())
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

	// Fetch the generated account from the network.
	client := horizonclient.DefaultTestNetClient
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: generatedKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a set options operation to set the
	// Authorization Required flag.
	setOp := txnbuild.SetOptions{
		SetFlags: []txnbuild.AccountFlag{
			txnbuild.AuthRevocable,
		},
	}

	// Build a change trust operation to allow the
	// quest account to receive the asset.
	asset := txnbuild.CreditAsset{
		Code:   "SQ",
		Issuer: generatedKp.Address(),
	}
	changeOp := txnbuild.ChangeTrust{
		Line:          asset.MustToChangeTrustAsset(),
		SourceAccount: questKp.Address(),
	}

	// Build a set trust line flags operation to
	// add the quest account as a trustor.
	addOp := txnbuild.SetTrustLineFlags{
		Trustor: questKp.Address(),
		Asset:   asset,
		SetFlags: []txnbuild.TrustLineFlag{
			txnbuild.TrustLineAuthorized,
		},
	}

	// Build a payment operation to debit the
	// quest account.
	paymentOp := txnbuild.Payment{
		Destination: questKp.Address(),
		Amount:      "1",
		Asset:       asset,
	}

	// Build another set trust line flags operation to
	// remove the quest account as a trustor.
	removeOp := txnbuild.SetTrustLineFlags{
		Trustor: questKp.Address(),
		Asset:   asset,
		ClearFlags: []txnbuild.TrustLineFlag{
			txnbuild.TrustLineAuthorized,
		},
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
			IncrementSequenceNum: true,
			Operations: []txnbuild.Operation{
				&setOp, &changeOp, &addOp,
				&paymentOp, &removeOp,
			},
			BaseFee:    txnbuild.MinBaseFee,
			Timebounds: txnbuild.NewTimeout(300),
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

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
