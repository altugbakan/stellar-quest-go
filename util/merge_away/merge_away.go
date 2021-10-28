package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the secret key from user input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
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

	// Fund the generated account.
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

	// Fetch the account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build manage data operations to delete all data.
	var ops []txnbuild.Operation
	for data := range questAccount.Data {
		ops = append(ops, &txnbuild.ManageData{
			Name: data,
		})
	}

	// Get the offers created by the account.
	offers, err := client.Offers(horizonclient.OfferRequest{
		ForAccount: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build delete offer operations to delete all offers.
	for _, offer := range offers.Embedded.Records {
		op, err := txnbuild.DeleteOfferOp(offer.ID)
		if err != nil {
			log.Fatal(err)
		}
		ops = append(ops, &op)
	}

	// Build change trust operations to delete all trustlines.
	for _, balance := range questAccount.Balances {
		if balance.Asset.Type != "native" {
			// Send the asset back to the issuer.
			asset := txnbuild.CreditAsset{
				Code:   balance.Asset.Code,
				Issuer: balance.Asset.Issuer,
			}
			amount, err := strconv.ParseFloat(balance.Balance, 32)
			if err != nil {
				log.Fatal(err)
			}
			if amount > 0 {
				ops = append(ops, &txnbuild.Payment{
					Destination: balance.Asset.Issuer,
					Amount:      balance.Balance,
					Asset:       asset,
				})
			}
			// Remove the trustline
			changeTrustAsset, err := asset.ToChangeTrustAsset()
			if err != nil {
				log.Fatal(err)
			}
			ops = append(ops, &txnbuild.ChangeTrust{
				Line:  changeTrustAsset,
				Limit: "0",
			})
		}
	}

	// Build an account merge operation.
	ops = append(ops, &txnbuild.AccountMerge{
		Destination: generatedKp.Address(),
	})

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
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
