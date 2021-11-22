package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the secret key from user input.
	var address, personal string
	fmt.Printf("Please enter the quest account's public key: ")
	fmt.Scanln(&address)
	fmt.Printf("Please enter your personal account's public key: ")
	fmt.Scanln(&personal)

	// Fetch the account from the network.
	client := horizonclient.DefaultPublicNetClient
	account, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: address,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a set options operation to add the signer.
	op := txnbuild.SetOptions{
		LowThreshold:    txnbuild.NewThreshold(5),
		MediumThreshold: txnbuild.NewThreshold(5),
		HighThreshold:   txnbuild.NewThreshold(5),
		HomeDomain:      new(string),
		Signer: &txnbuild.Signer{
			Address: personal,
			Weight:  10,
		},
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              250,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print the XDR.
	xdr, err := tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTransaction XDR: %v\n\n", string(xdr))
	fmt.Scanln()
}
