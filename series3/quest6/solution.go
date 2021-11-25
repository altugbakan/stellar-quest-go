package main

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

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

	// Note that getting the image from the provided link causes
	// a difference in image metadata due to Cloudflare polish.
	// Using the image uploaded to GitHub solves this issue.

	// Get the base64 encoding of the image.
	resp, err = http.Get("https://raw.githubusercontent.com/" +
		"altugbakan/stellar-quest-go/main/series3/quest6/NFT.png")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	imgBase64 := base64.StdEncoding.EncodeToString(bytes)

	// Create a reader to split the encoding.
	reader := strings.NewReader(imgBase64)

	// Split the encoding into multiple manage data operations.
	var ops []txnbuild.Operation
	index := 0
	name := make([]byte, 62)
	value := make([]byte, 64)
	for reader.Len() > 0 {
		// Trim slices to amount of bytes read.
		count, _ := reader.Read(name)
		name = name[:count]
		count, _ = reader.Read(value)
		value = value[:count]

		// Append the operation.
		ops = append(ops, &txnbuild.ManageData{
			Name:  fmt.Sprintf("%02d", index) + string(name),
			Value: []byte(string(value)),
		})
		index++
	}

	// Construct the transaction .
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
