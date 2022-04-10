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
		fmt.Println("Successfully funded quest account.")
	} else {
		fmt.Println("Error funding quest account.")
	}

	// Fund and create the generated account.
	resp, err = http.Get("https://friendbot.stellar.org/?addr=" + generatedKp.Address())
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Get and print the response from friendbot.
	if resp.Status == "200 OK" {
		fmt.Println("Successfully funded generated account.")
	} else {
		fmt.Println("Error funding generated account.")
	}

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the generated account from the network.
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: generatedKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create the asset.
	asset := txnbuild.CreditAsset{
		Code:   "CSTM",
		Issuer: generatedKp.Address(),
	}

	// Create a sell offer from the generated account.
	offerOp := txnbuild.ManageBuyOffer{
		Selling: asset,
		Buying:  txnbuild.NativeAsset{},
		Amount:  "10",
		Price:   "1",
		OfferID: 0,
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&offerOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, generatedKp)
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

	// Build a change trust operation.
	changeTrustAsset, err := asset.ToChangeTrustAsset()
	if err != nil {
		log.Fatal(err)
	}
	trustOp := txnbuild.ChangeTrust{
		Line:  changeTrustAsset,
		Limit: "100",
	}

	// Build a path payment operation.
	paymentOp := txnbuild.PathPaymentStrictReceive{
		Destination: questKp.Address(),
		SendAsset:   txnbuild.NativeAsset{},
		SendMax:     "1",
		DestAsset:   asset,
		DestAmount:  "1",
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&trustOp, &paymentOp},
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
	status, err = client.SubmitTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
