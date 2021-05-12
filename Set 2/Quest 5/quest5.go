package main

import (
	"fmt"
	"log"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/network"
	"github.com/stellar/go/txnbuild"
)

func main() {
	// If you used the Quest 4 code to solve the question, you must wait
	// 3.5 days or create a new claimable balance to solve the question.

	// Enter your quest account secret key below.
	secret := "SAGXUH5I7IMSDT6RLCF7HSP4UISLYTF6FVAYLTRYX6KSNLJQN266JHPK"
	// Supply a custom claimable balance ID if you need it.
	balanceID := ""
	// ..........................................

	// Get the keypair of the quest account from the secret key.
	questAccount, _ := keypair.Parse(secret)

	// Fetch the quest account from the network.
	client := horizonclient.DefaultTestNetClient
	ar := horizonclient.AccountRequest{AccountID: questAccount.Address()}
	sourceAccount, err := client.AccountDetail(ar)
	if err != nil {
		log.Fatalln(err)
	}

	// If not supplied, fetch the claimable balance ID.
	if balanceID == "" {
		balances, err := client.ClaimableBalances(
			horizonclient.ClaimableBalanceRequest{
				Claimant: questAccount.Address(),
			},
		)
		if err != nil {
			log.Fatalln(err)
		}
		balanceID = balances.Embedded.Records[0].BalanceID
	}

	// Build a claim claimable balance operation.
	op := txnbuild.ClaimClaimableBalance{
		BalanceID:     balanceID,
		SourceAccount: questAccount.Address(),
	}

	// Construct the transaction with both operations.
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
	fmt.Println(status)
}