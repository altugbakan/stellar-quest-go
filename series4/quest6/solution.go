package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
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

	// Fetch the claimable balances of the account from the network.
	// Fetch it manually as the Stellar Go SDK does not have a limit parameter!
	resp, err = http.Get("https://horizon-testnet.stellar.org/claimable_balances/" +
		"?claimant=" + public + "&limit=200")
	if err != nil {
		log.Fatal(err)
	}

	// Get the claimable balances.
	claimableBalancesByte, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var jsonResponse map[string]map[string][]map[string]string // ugly but works!
	json.Unmarshal(claimableBalancesByte, &jsonResponse)
	resp.Body.Close()

	// Get the claimable balance IDs and the total amount.
	var totalClaimAmount float64
	var balanceIDs []string
	for _, balance := range jsonResponse["_embedded"]["records"] {
		if balance["asset"] == "native" {
			balanceIDs = append(balanceIDs, balance["id"])
			claimAmount, err := strconv.ParseFloat(balance["amount"], 64)
			if err != nil {
				log.Fatal(err)
			}
			totalClaimAmount += claimAmount
		}
	}
	if len(balanceIDs) == 0 {
		log.Fatal("no balances found")
	}

	// Build claim claimable balance operations.
	var ops []txnbuild.Operation
	for _, balanceID := range balanceIDs {
		if err != nil {
			log.Fatal(err)
		}
		ops = append(ops, &txnbuild.ClaimClaimableBalance{
			BalanceID: balanceID,
		})
	}

	// Get the quest account's balance
	balance, err := questAccount.GetNativeBalance()
	if err != nil {
		log.Fatal(err)
	}
	balanceFloat, err := strconv.ParseFloat(balance, 64)
	if err != nil {
		log.Fatal(err)
	}

	// Calculate the maximum payable balance.
	maxPayableBalance := totalClaimAmount + balanceFloat - 1 - (0.5 * float64(len(questAccount.Signers)))

	// Build a payment operation, merge account is
	// not initially available for this quest.
	ops = append(ops, &txnbuild.Payment{
		Destination: generatedKp.Address(),
		Amount:      fmt.Sprintf("%.2f", maxPayableBalance),
		Asset:       txnbuild.NativeAsset{},
	})

	// Construct the transaction.
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

	// Print the XDR.
	xdr, err := tx.MarshalText()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("\nTransaction XDR: %v\n\n", string(xdr))

	// Inform the user and wait for user input to exit.
	fmt.Printf("Account's balance will be sent to %v, "+
		"which has the secret key %v.\nPress \"Enter\" to exit.\n",
		generatedKp.Address(), generatedKp.Seed())
	fmt.Scanln()
}
