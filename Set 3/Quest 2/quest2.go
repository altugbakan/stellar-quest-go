package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Ask user for secret key input.
	var secret string
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the quest account from the secret key.
	questAccount, _ := keypair.Parse(secret)

	// Fund and create the quest account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + questAccount.Address())
	if err != nil {
		log.Fatal(err)
	}

	// Get and print the response from friendbot.
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Friendbot response:")
	fmt.Println(string(body))

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	// Build 100 dummy payment operations.
	var ops []txnbuild.Operation
	for i := 0; i < 100; i++ {
		ops = append(ops, &txnbuild.Payment{
			Destination: questAccount.Address(),
			Amount:      "0.01",
			Asset:       txnbuild.NativeAsset{},
		})
	}

	// Construct a transaction with 100 operations.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
			IncrementSequenceNum: true,
			Operations:           ops,
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(1000),
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

	// Wait for user input to close.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
