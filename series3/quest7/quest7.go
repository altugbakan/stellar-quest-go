package main

import (
	"encoding/json"
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
	fmt.Printf("Please enter your secret key: ")
	fmt.Scanln(&secret)

	// Get the keypair of the quest account from the secret key.
	questAccount, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fund and create the quest account.
	resp, err := http.Get("https://friendbot.stellar.org/?addr=" + questAccount.Address())
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
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatal(err)
	}

	// Send a SEP-0010 request.
	resp, err = http.Get("https://testanchor.stellar.org/auth?account=" + questAccount.Address())
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Parse the response.
	var response map[string]string
	json.Unmarshal(body, &response)

	// Sign the response transaction.
	respTxGeneric, err := txnbuild.TransactionFromXDR(response["transaction"])
	if err != nil {
		log.Fatal(err)
	}
	respTx, _ := respTxGeneric.Transaction()
	signedTx, err := respTx.Sign(network.TestNetworkPassphrase, questAccount)
	if err != nil {
		log.Fatal(err)
	}
	signedStr, _ := signedTx.Base64()

	// Post the signed transaction to get the SEP-0010 token.
	body, _ = json.Marshal(map[string]string{"transaction": signedStr})
	resp, err = http.Post("https://testanchor.stellar.org/auth?account="+questAccount.Address(),
		"application/json", strings.NewReader(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()
	json.Unmarshal(body, &response)
	token := response["token"]
	fmt.Printf("Got token: %s\n", token)

	// Create a reader to split the encoding.
	reader := strings.NewReader(token)

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

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &sourceAccount,
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
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount)
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
