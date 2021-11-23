package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// Get the secret key from user input.
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

	// Fetch the account from the network.
	client := horizonclient.DefaultTestNetClient
	generatedAccount, err := client.AccountDetail(horizonclient.AccountRequest{
		AccountID: questKp.Address(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// Create the asset
	asset := txnbuild.CreditAsset{
		Code:   "SQuid",
		Issuer: "GCCDCZZP7AGSP2F2VRKDBLMXSQPRGR5MRITY2HUWQ3KIS6M6B436IUUF",
	}

	// Fetch the claimable balances of the quest account from the network.
	claimableBalances, err := client.ClaimableBalances(
		horizonclient.ClaimableBalanceRequest{
			Claimant: questKp.Address(),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Get the SQuid claimable balance from the quest account.
	var balanceID, amount string
	for _, balance := range claimableBalances.Embedded.Records {
		if balance.Asset == asset.Code+":"+asset.Issuer {
			balanceID = balance.BalanceID
			amount = balance.Amount
			break
		}
	}
	if balanceID == "" {
		log.Fatal("no balances to claim")
	}

	// Build a change trust operation
	var ops []txnbuild.Operation
	changeTrustAsset, err := asset.ToChangeTrustAsset()
	if err != nil {
		log.Fatal(err)
	}
	ops = append(ops, &txnbuild.ChangeTrust{
		Line: changeTrustAsset,
	})

	// Build a claim claimable balance operation.
	ops = append(ops, &txnbuild.ClaimClaimableBalance{
		BalanceID:     balanceID,
		SourceAccount: questKp.Address(),
	})

	// Build a manage sell offer operation.
	ops = append(ops, &txnbuild.ManageSellOffer{
		Selling: asset,
		Buying:  txnbuild.NativeAsset{},
		Amount:  amount,
		Price:   "1",
	})

	// Build a change trust operation to remove trust
	ops = append(ops, &txnbuild.ChangeTrust{
		Line:  changeTrustAsset,
		Limit: "0",
	})

	// Build an account merge operation.
	ops = append(ops, &txnbuild.AccountMerge{
		Destination: generatedKp.Address(),
	})

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &generatedAccount,
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
