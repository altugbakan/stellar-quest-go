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
	// Enter your quest account secret key below.
	secret := "SAGXUH5I7IMSDT6RLCF7HSP4UISLYTF6FVAYLTRYX6KSNLJQN266JHPK"
	// ..........................................

	// Get the keypair of the quest account from the secret key.
	questAccount, _ := keypair.Parse(secret)

	// Generate a random testnet account.
	pair, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("The generated secret key is %v\n", pair.Seed())
	fmt.Printf("The generated public key is %v\n", pair.Address())

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	// Build a begin sponsorship operation.
	beginOp := txnbuild.BeginSponsoringFutureReserves{
		SponsoredID: pair.Address(),
	}

	// Build a create account operation.
	createOp := txnbuild.CreateAccount{
		Destination: pair.Address(),
		Amount:      "0",
	}

	// Build an end sponsorship operation.
	endOp := txnbuild.EndSponsoringFutureReserves{
		SourceAccount: pair.Address(),
	}

	// Construct the transaction with multiple operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&beginOp, &createOp, &endOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Sign the transaction with both keys.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full), pair)
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network.
	status, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the response.
	fmt.Println(status)
}