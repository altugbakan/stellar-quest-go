package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

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
	questAccount, _ := keypair.Parse(secret)

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

	// Create the asset
	asset := txnbuild.CreditAsset{
		Code:   "MULT",
		Issuer: "GDLD3SOLYJTBEAK5IU4LDS44UMBND262IXPJB3LDHXOZ3S2QQRD5FSMM",
	}

	// Build a change trust operation.
	changeTrustAsset, err := asset.ToChangeTrustAsset()
	op := txnbuild.ChangeTrust{
		Line: changeTrustAsset,
	}
	if err != nil {
		log.Fatal(err)
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
		log.Fatal(err)
	}

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full))
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
	signedTx, err := respTx.Sign(network.TestNetworkPassphrase, questAccount.(*keypair.Full))
	if err != nil {
		log.Fatal(err)
	}
	signedStr, _ := signedTx.Base64()

	// Post the signed transaction to get the SEP-0010 token.
	body, _ = json.Marshal(map[string]string{"transaction": signedStr})
	resp, err = http.Post("https://testanchor.stellar.org/auth?account="+questAccount.Address(),
		"application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Parse the response.
	json.Unmarshal(body, &response)
	token := response["token"]
	fmt.Printf("Got token: %s\n", token)

	// Fill the KYC information.
	kyc := map[string]string{
		"account":             questAccount.Address(),
		"first_name":          "Stellar",
		"last_name":           "Quest",
		"email_address":       "quest@stellar.org",
		"bank_number":         "07312014",
		"bank_account_number": "05282021",
	}

	// Send a PUT request to complete KYC.
	body, _ = json.Marshal(kyc)
	req, err := http.NewRequest("PUT", "https://testanchor.stellar.org/kyc/customer", bytes.NewBuffer(body))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Parse the response.
	json.Unmarshal(body, &response)
	id := response["id"]
	fmt.Printf("Submitted KYC. Account ID: %s\n", id)

	// Send the SEP-0006 request to deposit MULT to the quest account.
	req, err = http.NewRequest("GET", "https://testanchor.stellar.org/sep6/deposit?asset_code="+
		asset.Code+"&account="+questAccount.Address()+"&type="+"bank_account", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+token)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Parse the response.
	json.Unmarshal(body, &response)
	id = response["id"]
	fmt.Printf("Submitted deposit request. Request ID: %s\n", id)

	// Check the status of the request.
	req, err = http.NewRequest("GET", "https://testanchor.stellar.org/sep6/transaction?id="+id, nil)
	req.Header.Add("Authorization", "Bearer "+token)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Loop until the request is completed.
	for {
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		body, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var transactionWrapper map[string]map[string]string
		json.Unmarshal(body, &transactionWrapper)
		reqStatus := transactionWrapper["transaction"]["status"]
		fmt.Printf("Status of the request is: %s\n", reqStatus)
		if reqStatus == "completed" {
			break
		}
		time.Sleep(10 * time.Second)
	}

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
