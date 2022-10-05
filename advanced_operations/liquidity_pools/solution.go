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
	fmt.Printf("Please enter the quest account's secret key: ")
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

	// Build an account creation operation to create the random account first.
	createOp := txnbuild.CreateAccount{
		Destination: generatedKp.Address(),
		Amount:      "2000",
	}

	// Create the liquidity pool asset.
	asset := txnbuild.CreditAsset{
		Code:   "NOODLE",
		Issuer: questKp.Address(),
	}
	liquidityPool := txnbuild.LiquidityPoolShareChangeTrustAsset{
		LiquidityPoolParameters: txnbuild.LiquidityPoolParameters{
			AssetA: txnbuild.NativeAsset{},
			AssetB: asset,
			Fee:    30,
		},
	}

	// Get the liquidity pool ID.
	poolId, err := txnbuild.NewLiquidityPoolId(txnbuild.NativeAsset{}, asset)
	if err != nil {
		log.Fatal(err)
	}

	// Build a change trust operation for the liquidity.
	trustLiquidityOp := txnbuild.ChangeTrust{
		Line: liquidityPool,
	}

	// Build a liquidity pool deposit operation.
	depositOp := txnbuild.LiquidityPoolDeposit{
		LiquidityPoolID: poolId,
		MaxAmountA:      "100",
		MaxAmountB:      "100",
		MinPrice:        "1",
		MaxPrice:        "1",
	}

	// Build a change trust operation for the asset.
	trustAssetOp := txnbuild.ChangeTrust{
		Line:          asset.MustToChangeTrustAsset(),
		SourceAccount: generatedKp.Address(),
	}

	// Build a path payment strict receive operation.
	pathOp := txnbuild.PathPaymentStrictReceive{
		Destination: generatedKp.Address(),
		SendAsset:   txnbuild.NativeAsset{},
		SendMax:     "1000",
		DestAsset:   asset,
		DestAmount:  "1",
		Path:        []txnbuild.Asset{},
	}

	// Build a liquidity pool withdraw operation.
	withdrawOp := txnbuild.LiquidityPoolWithdraw{
		LiquidityPoolID: poolId,
		Amount:          "100",
		MinAmountA:      "0",
		MinAmountB:      "0",
	}

	// Construct the transaction.
	tx, err := txnbuild.NewTransaction(
		txnbuild.TransactionParams{
			SourceAccount:        &questAccount,
			IncrementSequenceNum: true,
			Operations:           []txnbuild.Operation{&createOp, &trustLiquidityOp, &depositOp, &trustAssetOp, &pathOp, &withdrawOp},
			BaseFee:              txnbuild.MinBaseFee,
			Timebounds:           txnbuild.NewTimeout(300),
		},
	)
	if err != nil {
		log.Fatal(err)
	}

	// Sign the transaction.
	tx, err = tx.Sign(network.TestNetworkPassphrase, questKp, generatedKp)
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
