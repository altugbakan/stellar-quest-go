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
	// Get the user's choice of which part of the quest they will solve.
	var part int
	fmt.Printf("1 to create claimable balance, 2 to claim claimable balance: ")
	fmt.Scanln(&part)

	if part == 1 {
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

		// Fetch the quest account from the network.
		client := horizonclient.DefaultTestNetClient
		questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
			AccountID: questKp.Address(),
		})
		if err != nil {
			log.Fatal(err)
		}

		// Create the condition.
		predicate := txnbuild.NotPredicate(txnbuild.BeforeRelativeTimePredicate(300))

		// Build an account creation operation to create the random account first.
		createOp := txnbuild.CreateAccount{
			Destination: generatedKp.Address(),
			Amount:      "100",
		}

		// Build a create claimable balance operation.
		claimOp := txnbuild.CreateClaimableBalance{
			Amount: "100",
			Asset:  txnbuild.NativeAsset{},
			Destinations: []txnbuild.Claimant{
				txnbuild.NewClaimant(generatedKp.Address(), &predicate),
			},
		}

		// Construct the transaction.
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &questAccount,
				IncrementSequenceNum: true,
				Operations:           []txnbuild.Operation{&createOp, &claimOp},
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

		// Inform the user and wait for user input to exit.
		fmt.Println("A claimable balance is created to be claimed after 5 minutes.\nPress \"Enter\" to exit.")
		fmt.Scanln()
	} else if part == 2 {
		// Inform the user.
		fmt.Println("Please wait up to 5 minutes if you created the balance using this program.")

		// Get the secret key from user input.
		var secret string
		fmt.Printf("Please enter the generated account's secret key: ")
		fmt.Scanln(&secret)

		// Get the keypair of the generated account from the secret key.
		generatedKp, err := keypair.ParseFull(secret)
		if err != nil {
			log.Fatal(err)
		}

		// Fetch the generated account from the network.
		client := horizonclient.DefaultTestNetClient
		generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
			AccountID: generatedKp.Address(),
		})
		if err != nil {
			log.Fatal(err)
		}

		// Fetch the claimable balances of the quest account from the network.
		claimableBalances, err := client.ClaimableBalances(
			horizonclient.ClaimableBalanceRequest{
				Claimant: generatedKp.Address(),
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
			SourceAccount: generatedKp.Address(),
		}

		// Construct the transaction with both operations.
		tx, err := txnbuild.NewTransaction(
			txnbuild.TransactionParams{
				SourceAccount:        &generatedAccount,
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

		// Wait for user input to exit.
		fmt.Println("Press \"Enter\" to exit.")
		fmt.Scanln()
	} else {
		// Inform the user and wait for user input to exit.
		fmt.Println("Invalid input.\nPress \"Enter\" to exit.")
		fmt.Scanln()
	}
}
