package main

import (
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

	// Fund the quest account.
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

	// Build a dummy payment operation.
	dummyOp := txnbuild.Payment{
		Destination: questAccount.Address(),
		Amount:      "0.1",
		Asset:       txnbuild.NativeAsset{},
	}

	// Increment the source account's sequence, as the pre-authorized transaction
	// will be used in the future, requiring a new sequence number.
	sequenceNumber, err := sourceAccount.GetSequenceNumber()
	if err != nil {
		log.Fatal(err)
	}
	sourceAccount.Sequence = fmt.Sprintf("%d", sequenceNumber+1)

	// Construct the pre-authorized transaction.
	preAuthTx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&dummyOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get the hash and base64 XDR of the pre-authorized transaction.
	txHash, err := preAuthTx.Hash(network.TestNetworkPassphrase)
	if err != nil {
		log.Fatal(err)
	}
	hashAddress, err := strkey.Encode(strkey.VersionByteHashTx, txHash[:])
	if err != nil {
		log.Fatal(err)
	}

	// Construct a set options operation to add the pre-authorized
	// transaction to the account's signers.
	authOp := txnbuild.SetOptions{
		Signer: &txnbuild.Signer{
			Address: hashAddress,
			Weight:  txnbuild.Threshold(1),
		},
	}

	// Revert the source account's sequence to normal.
	sourceAccount.Sequence = fmt.Sprintf("%d", sequenceNumber)

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&authOp},
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

	// Send the previously authorized transaction to the network.
	// Note that this transaction is not signed using the secret key.
	status, err = client.SubmitTransaction(preAuthTx)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to close.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
