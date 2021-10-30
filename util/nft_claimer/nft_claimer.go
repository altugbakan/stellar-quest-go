package main

import (
	"fmt"
	"log"
	"strings"

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
	kp, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the account from the network.
	client := horizonclient.DefaultPublicNetClient
	account, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: kp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Fetch the claimable balances of the account from the network.
	claimableBalances, err := client.ClaimableBalances(horizonclient.ClaimableBalanceRequest{
		Claimant: kp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get the claimable balance IDs from the Stellar Quest NFT issuers.
	var balanceIDs, balanceAmounts, issuers, assetNames []string
	for _, balance := range claimableBalances.Embedded.Records {
		if strings.HasPrefix(balance.Asset, "SQ") {
			balanceIDs = append(balanceIDs, balance.BalanceID)
			balanceAmounts = append(balanceAmounts, balance.Amount)
			assetNames = append(assetNames, balance.Asset[:strings.Index(balance.Asset, ":")])
			issuers = append(issuers, balance.Asset[strings.Index(balance.Asset, ":")+1:])
		}
	}
	if len(balanceIDs) == 0 {
		fmt.Println("No NFTs found.\nPress \"Enter\" to exit.")
		fmt.Scanln()
		return
	}

	// Build change trust and claim claimable balance operations.
	var ops []txnbuild.Operation
	for i := range balanceIDs {
		asset, err := txnbuild.CreditAsset{
			Code:   assetNames[i],
			Issuer: issuers[i],
		}.ToChangeTrustAsset()
		if err != nil {
			log.Fatal(err)
		}
		ops = append(ops, &txnbuild.ChangeTrust{
			Line:  asset,
			Limit: balanceAmounts[i],
		})
		ops = append(ops, &txnbuild.ClaimClaimableBalance{
			BalanceID: balanceIDs[i],
		})
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &account,
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
	tx, err = tx.Sign(network.PublicNetworkPassphrase, kp)
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
	fmt.Print("The claimed NFTs are: ")
	for i, asset := range assetNames {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%v", asset)
	}
	fmt.Println("\nPress \"Enter\" to exit.")
	fmt.Scanln()
}
