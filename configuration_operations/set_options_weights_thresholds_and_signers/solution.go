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

	// Fetch the account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a set options operation to set the weights.
	setOp := txnbuild.SetOptions{
		MasterWeight:    txnbuild.NewThreshold(1),
		LowThreshold:    txnbuild.NewThreshold(5),
		MediumThreshold: txnbuild.NewThreshold(5),
		HighThreshold:   txnbuild.NewThreshold(5),
	}

	// Build two more set options operations to add two
	// more signers.
	signer1, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	signer2, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	addOp1 := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: signer1.Address(),
			Weight:  2,
		},
	}
	addOp2 := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: signer2.Address(),
			Weight:  2,
		},
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&setOp, &addOp1, &addOp2},
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

	// Build a dummy payment operation.
	dummyOp := txnbuild.Payment{
		Destination: questKp.Address(),
		Amount:      "1",
		Asset:       txnbuild.NativeAsset{},
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&dummyOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction using 3 signers.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp, signer1, signer2)
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
