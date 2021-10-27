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
	// Ask user for secret key and the account sponsored input.
	var secret, sponsoredID string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)
	fmt.Printf("Please enter the public key of the account you " +
		"sponsored on Quest 6: ")
	fmt.Scanln(&sponsoredID)

	// Get the keypair of the quest account from the secret key.
	questAccount, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatal(err)
	}

	// Build a payment operation to ensure the sponsored account's
	// balance does not fall below the minimum balance.
	paymentOp := txnbuild.Payment{
		Destination: sponsoredID,
		Amount:      "5",
		Asset:       txnbuild.NativeAsset{},
	}

	// Build a revoke sponsorship operation.
	revokeOp := txnbuild.RevokeSponsorship{
		SponsorshipType: txnbuild.RevokeSponsorshipTypeAccount,
		Account:         &sponsoredID,
	}

	// Construct the transaction with both operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&paymentOp, &revokeOp},
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

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
