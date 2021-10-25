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
	// Get the secret key, asset name and the issuer address
	// created on Quest 5 below from user input.
	var secret, assetName, issuerAddress string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)
	fmt.Printf("Please enter the name of the asset you created in Quest 5: ")
	fmt.Scanln(&assetName)
	fmt.Printf("Please enter the public key of the issuer address of the" +
		" asset you created on Quest 5: ")
	fmt.Scanln(&issuerAddress)

	// Get the keypair of the quest account from the secret key.
	questAccount, _ := keypair.Parse(secret)

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	// Create the asset.
	asset := txnbuild.CreditAsset{
		Code:   assetName,
		Issuer: issuerAddress,
	}

	// Build a sell offer operation.
	op := txnbuild.ManageSellOffer{
		Selling: asset,
		Buying:  txnbuild.NativeAsset{},
		Amount:  "0.1",
		Price:   "0.1",
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&op},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatalln(err)
	}

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full))
	if err != nil {
		log.Fatalln(err)
	}

	// Send the transaction to the network.
	status, err := client.SubmitTransaction(tx)
	if err != nil {
		log.Fatalln(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
