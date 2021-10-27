package main

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Ask user for secret key input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the quest account from the secret key.
	questAccount, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fund and create the quest account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + questAccount.Address())
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
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatal(err)
	}

	// Convert the given clue to string, then hash it.
	secretKey, _ := base64.StdEncoding.DecodeString("S2FuYXllTmV0")
	hash := sha256.Sum256(secretKey)
	hashAddress, err := strkey.Encode(strkey.VersionByteHashX, hash[:])
	if err != nil {
		log.Fatal(err)
	}

	// Build a set options operation to add the signer.
	op := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: hashAddress,
			Weight:  txnbuild.Threshold(1),
		},
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
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
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount)
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

	// Build a set options operation to remove the signer.
	op = txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: hashAddress,
			Weight:  txnbuild.Threshold(0),
		},
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the previously added hash.
	tx, err = tx.SignHashX(secretKey)
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

	// Wait for user input to close.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
