package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"reflect"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the public key from user input.
	var public string
	fmt.Printf("Please enter the quest's public key: ")
	fmt.Scanln(&public)

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

	// Fetch the account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: public,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Get the account and signer quest account is sponsoring.
	accounts, err := client.Accounts(horizonclient.AccountsRequest{
		Sponsor: public,
	})
	if err != nil {
		log.Fatal(err)
	}
	sponsoredAccount := accounts.Embedded.Records[0].AccountID
	var sponsoredSigner string
	for _, signer := range accounts.Embedded.Records[0].Signers {
		if signer.Sponsor == public {
			sponsoredSigner = signer.Key
		}
	}
	if sponsoredSigner == "" {
		log.Fatal("no signer found")
	}

	// We'll cheat a bit here lol
	//---------------------------------CHEAT-START------------------------------------
	// Build a payment operation.
	paymentOp := txnbuild.Payment{
		Destination: sponsoredAccount,
		Amount:      "5",
		Asset:       txnbuild.NativeAsset{},
	}

	// Fetch the account from the network.
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: generatedKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Send the payment operation from the friendbot funded account,
	// as the quest account does not have enough XLM to revoke sponsorship.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&paymentOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the secret hash.
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
	//----------------------------------CHEAT-END-------------------------------------

	// Build a revoke sponsorship operation for the sponsored account.
	revokeAccOp := txnbuild.RevokeSponsorship{
		SponsorshipType: txnbuild.RevokeSponsorshipTypeAccount,
		Account:         &sponsoredAccount,
	}

	// Build a revoke sponsorship operation for the sponsored signer.
	revokeSignOp := txnbuild.RevokeSponsorship{
		SponsorshipType: txnbuild.RevokeSponsorshipTypeSigner,
		Signer: &txnbuild.SignerID{
			AccountID:     sponsoredAccount,
			SignerAddress: sponsoredSigner,
		},
	}

	// Build an account merge operation.
	mergeOp := txnbuild.AccountMerge{
		Destination: generatedKp.Address(),
	}

	// Get the operations for the account.
	operations, err := client.Operations(horizonclient.OperationRequest{
		ForAccount:    public,
		IncludeFailed: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Find the secret hash.
	var encodedHash string
	for _, operation := range operations.Embedded.Records {
		opRefVal := reflect.ValueOf(operation)
		if !operation.IsTransactionSuccessful() &&
			opRefVal.FieldByName("Name").String() ==
				"I AM THE SIGNATURE YOU'RE LOOKING FOR" {
			encodedHash = opRefVal.FieldByName("Value").String()
			break
		}
	}
	if encodedHash == "" {
		log.Fatal("no hash found")
	}

	// Decode the hash.
	hash, err := base64.StdEncoding.DecodeString(encodedHash)
	if err != nil {
		log.Fatal(err)
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&revokeAccOp, &revokeSignOp, &mergeOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction with the secret hash.
	tx, err = tx.SignHashX([]byte(hash))
	if err != nil {
		log.Fatal(err)
	}

	// Print the XDR.
	xdr, err := tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTransaction XDR: %v\n\n", string(xdr))

	// Inform the user and wait for user input to exit.
	fmt.Printf("Account will be merged to %v, "+
		"which has the secret key %v.\nPress \"Enter\" to exit.\n",
		generatedKp.Address(), generatedKp.Seed())
	fmt.Scanln()
}
