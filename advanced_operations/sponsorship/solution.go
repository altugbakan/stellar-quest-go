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

	// Fund and create the generated account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + generatedKp.Address())
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

	// Fetch the quest account from the network if it exists.
	client := horizonclient.DefaultTestNetClient
	questAccount, _ := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})

	// Fetch the generated account from the network.
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: generatedKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get the balance of the quest account.
	accountBalance := "0.0000000"
	for _, balance := range questAccount.Balances {
		if balance.Asset.Type == "native" {
			accountBalance = balance.Balance
			break
		}
	}

	// If account balance is 0, do Part 1 of the solution.
	var ops []txnbuild.Operation
	if accountBalance == "0.0000000" {
		// Build begin sponsoring future reserves operation.
		ops = append(ops, &txnbuild.BeginSponsoringFutureReserves{
			SponsoredID: questKp.Address(),
		})

		// Build create account operation.
		ops = append(ops, &txnbuild.CreateAccount{
			Destination: questKp.Address(),
			Amount:      "0",
		})

		// Build end sponsoring future reserves operation.
		ops = append(ops, &txnbuild.EndSponsoringFutureReserves{
			SourceAccount: questKp.Address(),
		})
	} else {
		// Else, do Part 2 of the solution

		// Build begin sponsoring future reserves operation.
		ops = append(ops, &txnbuild.BeginSponsoringFutureReserves{
			SponsoredID: questKp.Address(),
		})

		// Build revoke sponsorship operation.
		questAddress := questKp.Address()
		ops = append(ops, &txnbuild.RevokeSponsorship{
			SourceAccount:   questKp.Address(),
			SponsorshipType: 1,
			Account:         &questAddress,
		})

		// Build end sponsoring future reserves operation.
		ops = append(ops, &txnbuild.EndSponsoringFutureReserves{
			SourceAccount: questKp.Address(),
		})

		// Build payment operation.
		ops = append(ops, &txnbuild.Payment{
			Destination:   generatedKp.Address(),
			Amount:        accountBalance,
			Asset:         txnbuild.NativeAsset{},
			SourceAccount: questKp.Address(),
		})
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, generatedKp, questKp)
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
