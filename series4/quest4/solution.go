package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Ask the user to send the first XDR.
	ans := "n"
	fmt.Printf("Did you send the XDR given? y/[n]: ")
	fmt.Scanln(&ans)

	// If they didn't, quit.
	if ans != "y" {
		fmt.Println("Please send the XDR given. If you get an error," +
			" please wait for up to 5 minutes before sending.")

		// Wait for user input to close.
		fmt.Println("Press \"Enter\" to exit.")
		fmt.Scanln()
		return
	}

	// Get the link from user input.
	var xdrLink string
	fmt.Printf("Please enter the link of the XDR that you have sent: ")
	fmt.Scanln(&xdrLink)

	// Fetch the quest XDR.
	resp, err := http.Get(xdrLink)
	if err != nil {
		log.Fatal(err)
	}

	// Get and print the response.
	if resp.Status == "200 OK" {
		fmt.Println("Successfully fetched the quest XDR.")
	} else {
		log.Fatal("Error fetching the quest XDR.")
	}

	// Get the quest XDR.
	questXdrBytes, _ := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	questXdr := string(questXdrBytes)
	resp.Body.Close()

	// Generate a random testnet account.
	generatedKp, err := keypair.Random()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The generated secret key is %v\n", generatedKp.Seed())
	fmt.Printf("The generated public key is %v\n", generatedKp.Address())

	// Fund and create the generated account.
	resp, err = http.Get("https://friendbot.stellar.org/?addr=" + generatedKp.Address())
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

	// Convert the XDR to transaction.
	xdrGenericTx, err := txnbuild.TransactionFromXDR(questXdr)
	if err != nil {
		log.Fatal(err)
	}

	// Get the merge destination.
	xdrTx, _ := xdrGenericTx.Transaction()
	lastOp := xdrTx.Operations()[len(xdrTx.Operations())-1]
	destination := lastOp.(*txnbuild.AccountMerge).Destination

	// Fetch the merge destination account from network.
	client := horizonclient.DefaultTestNetClient
	destinationAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: destination,
	})
	if err != nil {
		log.Fatal(err)
	}

	// Build an account merge operation.
	mergeOp := txnbuild.AccountMerge{
		Destination: generatedKp.Address(),
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &destinationAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&mergeOp},
			BaseFee:              txnbuild.MinBaseFee,
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
	fmt.Printf("\nTransaction XDR: %v\n\n", string(xdr))

	// Inform the user and wait for user input to exit.
	fmt.Printf("Account will be merged to %v, "+
		"which has the secret key %v.\nPress \"Enter\" to exit.\n",
		generatedKp.Address(), generatedKp.Seed())
	fmt.Scanln()
}
