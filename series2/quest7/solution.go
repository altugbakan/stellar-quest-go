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
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the quest account from the secret key.
	questKp, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the accounts sponsored by the quest account.
	sponsoredAccounts, err := client.Accounts(horizonclient.AccountsRequest{
		Sponsor: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get a sponsored account, any will work.
	if len(sponsoredAccounts.Embedded.Records) == 0 {
		log.Fatal("no sponsored account found")
	}
	sponsoredAccount := sponsoredAccounts.Embedded.Records[0].AccountID

	// Build a payment operation to ensure the sponsored account's
	// balance does not fall below the minimum balance.
	paymentOp := txnbuild.Payment{
		Destination: sponsoredAccount,
		Amount:      "5",
		Asset:       txnbuild.NativeAsset{},
	}

	// Build a revoke sponsorship operation.
	revokeOp := txnbuild.RevokeSponsorship{
		SponsorshipType: txnbuild.RevokeSponsorshipTypeAccount,
		Account:         &sponsoredAccount,
	}

	// Construct the transaction with both operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
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
