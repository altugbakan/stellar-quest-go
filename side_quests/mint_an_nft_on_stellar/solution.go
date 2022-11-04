package main

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

type metadata struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Url         string `json:"url"`
	Issuer      string `json:"issuer"`
	Code        string `json:"code"`
}

type nftStorageResponse struct {
	Value struct {
		Cid string `json:"cid"`
	} `json:"value"`
}

//go:embed pug.png
var image embed.FS

func main() {
	// Get the NFT.Storage API key from user input
	var apiKey string
	fmt.Printf("Please enter your NFT.Storage API key: ")
	fmt.Scanln(&apiKey)

	// Get the secret key from user input.
	var secret string
	fmt.Printf("Please enter the quest account's secret key: ")
	fmt.Scanln(&secret)

	// Store the image using NFT.Storage
	img, err := image.ReadFile("pug.png")
	if err != nil {
		log.Fatal(err)
	}
	req, err := http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewReader(img))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var imageResponse nftStorageResponse
	err = json.Unmarshal(body, &imageResponse)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Get and print the response from NFT.Storage.
	if resp.Status == "200 OK" {
		fmt.Println("Successfully uploaded image.")
	} else {
		fmt.Println("Error uploading image.")
	}

	// Generate a random testnet account.
	generatedKp, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The generated secret key is %v\n", generatedKp.Seed())
	fmt.Printf("The generated public key is %v\n", generatedKp.Address())

	// Create the NFT asset
	nftAsset := txnbuild.CreditAsset{
		Code:   "pugNFT",
		Issuer: generatedKp.Address(),
	}

	// Create the NFT metadata
	nftMetadata := &metadata{
		Name:        "Pug",
		Description: "Cutest pug in the world!",
		Url:         "ipfs://" + imageResponse.Value.Cid,
		Issuer:      nftAsset.Issuer,
		Code:        nftAsset.Code,
	}
	nftMetadataJson, err := json.Marshal(nftMetadata)
	if err != nil {
		log.Fatal(err)
	}

	// Store the metadata using NFT.Storage
	req, err = http.NewRequest("POST", "https://api.nft.storage/upload", bytes.NewBuffer(nftMetadataJson))
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+apiKey)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var metadataResponse nftStorageResponse
	err = json.Unmarshal(body, &metadataResponse)
	if err != nil {
		log.Fatal(err)
	}
	resp.Body.Close()

	// Get and print the response from NFT.Storage.
	if resp.Status == "200 OK" {
		fmt.Println("Successfully uploaded metadata.")
	} else {
		fmt.Println("Error uploading metadata.")
	}

	// Get the keypair of the quest account from the secret key.
	questKp, err := keypair.ParseFull(secret)
	if err != nil {
		log.Fatal(err)
	}

	// Fund and create the quest account.
	resp, err = http.Get("https://friendbot.stellar.org/?addr=" + questKp.Address())
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

	// Build an account creation operation to create the random account first.
	createOp := txnbuild.CreateAccount{
		Destination: generatedKp.Address(),
		Amount:      "2",
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&createOp},
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

	// Fetch the generated account from the network.
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: generatedKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build a manage data operation to add the NFT metadata
	manageOp := txnbuild.ManageData{
		Name:  "ipfshash",
		Value: []byte(metadataResponse.Value.Cid),
	}

	// Build a change trust operation to accept the NFT
	changeTrustAsset, err := nftAsset.ToChangeTrustAsset()
	if err != nil {
		log.Fatal(err)
	}
	trustOp := txnbuild.ChangeTrust{
		Line:          changeTrustAsset,
		Limit:         "0.0000001",
		SourceAccount: questKp.Address(),
	}

	// Build a payment operation to mint the NFT
	paymentOp := txnbuild.Payment{
		Destination: questKp.Address(),
		Amount:      "0.0000001",
		Asset:       nftAsset,
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&manageOp, &trustOp, &paymentOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction using both keys.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp, generatedKp)
	if err != nil {
		log.Fatal(err)
	}

	// Send the transaction to the network.
	status, err = client.SubmitTransaction(tx)
	if err != nil {
		log.Fatal(err)
	}

	// Print the response.
	fmt.Printf("Successfully submitted transaction!\nTransaction ID: %v\n", status.ID)

	// Wait for user input to exit.
	fmt.Println("Press \"Enter\" to exit.")
	fmt.Scanln()
}
