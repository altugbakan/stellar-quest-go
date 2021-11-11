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
	// Inform the user.
	fmt.Println("Please wait up to 2 minutes if you solved Quest 4 using the previous solution.")

	// Ask user for secret key input.
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

	// Fetch the claimable balances of the quest account from the network.
	claimableBalances, err := client.ClaimableBalances(
		horizonclient.ClaimableBalanceRequest{
			Claimant: questKp.Address(),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get a native claimable balance from the quest account, any will work.
	var balanceID string
	for _, balance := range claimableBalances.Embedded.Records {
		if balance.Asset == "native" {
			balanceID = balance.BalanceID
			break
		}
	}
	if balanceID == "" {
		log.Fatal("no balances to claim")
	}

	// Build a claim claimable balance operation.
	op := txnbuild.ClaimClaimableBalance{
		BalanceID:     balanceID,
		SourceAccount: questKp.Address(),
	}

	// Construct the transaction with both operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
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
