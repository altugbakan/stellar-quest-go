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
	// Get the secret key and emoji from user input.
	var secret string
	fmt.Printf("Please enter the quest's secret key: ")
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

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	questAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build the pre-authorized set options operation.
	preAuthOp := txnbuild.SetOptions{
		MasterWeight: txnbuild.NewThreshold(1),
	}

	// Construct the pre-authorized transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&preAuthOp},
			BaseFee:              10000000,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Print the XDR.
	xdr, err := tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTransaction 1 XDR: %v\n\n", string(xdr))

	// Get the quest account's balance
	balance, err := questAccount.GetNativeBalance()
	if err != nil {
		log.Fatal(err)
	}
	balanceFloat, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the maximum payable balance
	// Note that we subtract the pre-authorized transaction.
	maxPayableBalance := balanceFloat - 1 - (0.5 * float64(len(questAccount.Signers)-1))

	// Build a payment operation, merge account is
	// not initially available for this quest.
	paymentOp := txnbuild.Payment{
		Destination: generatedKp.Address(),
		Amount:      fmt.Sprintf("%.2f", maxPayableBalance),
		Asset:       txnbuild.NativeAsset{},
	}

	// Construct the transaction.
	tx, err = txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&paymentOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewInfiniteTimeout(),
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

	// Print the XDR.
	xdr, err = tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Transaction 2 XDR: %v\n\n", string(xdr))

	// Inform the user and wait for user input to exit.
	fmt.Printf("Account's balance will be sent to %v, "+
		"which has the secret key %v.\nPress \"Enter\" to exit.\n",
		generatedKp.Address(), generatedKp.Seed())
	fmt.Scanln()
}
