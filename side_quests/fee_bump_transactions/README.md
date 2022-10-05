# Side Quest - Fee-Bump Transactions

### Pay the transaction fee for another account

<br>

Welcome, quester! We are thrilled that you’ve decided to venture off the beaten path to discover this fantastically fun side quest. It’s time to explore a whole new type of transaction, the fee-bump transaction.

**In this quest, you will perform a fee-bump transaction on the Stellar testnet with your Quest Account as the fee-bump source account.**

Let’s start with a little background on fees. Stellar requires a small fee for all transactions on the network. This prevents ledger spam and helps to prioritize transactions during surge pricing. The minimum fee for a given transaction equals the number of operations multiplied by the base fee for the ledger.

Base fee \* # of operations = transaction fee

Fees on Stellar are dynamic, so the base fee for a ledger can vary. Under normal conditions, the base fee is 100 stroops per operation. When the number of operations submitted to a ledger exceeds network capacity, however, the network enters [surge pricing mode](https://developers.stellar.org/docs/glossary/fees/#surge-pricing), and the fee you specify for a transaction serves as the highest bid you're willing to pay to make the ledger.

A stroop, by the way, is the smallest unit of a lumen (0.0000001 XLM).

A fee-bump transaction allows an account to pay the transaction fees for an existing transaction without having to re-sign it or manage sequence numbers. Fee-bumps are useful in a few scenarios:

- You’re building a service where you want to cover user fees
- You want to increase the fee on an existing transaction so it has a better chance of making it to the ledger during surge pricing
- You want to adjust the fee on preauthorized transactions so they can make it to the ledger if minimum network fees have increased

We’ll walk through this quest using the Laboratory and with code using the Stellar SDK. Expand each step to view the code portion. Knowledge of NodeJS is required for the SDK sections.

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Either create two new accounts on the testnet or use two you’ve already created- these will be your Sender Account and Destination Account

<details>
<summary>See it in code</summary>
<br>
Generate two new accounts: senderAccount being the account that sends the payment; destinationAccount being the account that receives the payment

```javascript
// include the StellarSDK
const StellarSdk = require("stellar-sdk");
const fetch = require("node-fetch");

// Generate two Keypairs: a sender, and a destination.
const senderKeypair = StellarSdk.Keypair.random();
const destinationKeypair = StellarSdk.Keypair.random();

// Optional: Log the keypair details if you want to save the information for later.
console.log(`Sender Public Key: ${senderKeypair.publicKey()}`);
console.log(`Sender Secret Key: ${senderKeypair.secret()}`);
console.log(`Desintation Public Key: ${destinationKeypair.publicKey()}`);
console.log(`Destination Secret Key: ${destinationKeypair.secret()}`);
```

Fund the accounts using Friendbot

```javascript
await Promise.all(
  [senderKeypair, destinationKeypair].map(async (kp) => {
    // Set up the Friendbot URL endpoints.
    const friendbotUrl = `https://friendbot.stellar.org?addr=${kp.publicKey()}`;
    let response = await fetch(friendbotUrl);

    // // Optional Looking at the responses from fetch.
    // let json = await response.json()
    // console.log(json)

    // Check that the response is OK, and give a confirmation message.
    if (response.ok) {
      console.log(`Account ${kp.publicKey()} successfully funded.`);
    } else {
      console.log(`Something went wrong funding account: ${kp.publicKey()}`);
    }
  })
);
```

</details>
<br>
3. Create the inner transaction by navigating to the Build Transaction tab and using a Payment operation to send 100 XLM from the Sender Account to the Destination Account

<details>
<summary>See it in code</summary>
<br>
Create the inner transaction where the senderAccount transfers 100 XLM to the destinationAccount

```javascript
// Connect to the testnet with the StellarSdk.
const server = new StellarSdk.Server("https://horizon-testnet.stellar.org");
const senderAccount = await server.loadAccount(senderKeypair.publicKey());

// Build the inner transaction. This will be the transaction where the payment is actually made.
let innerTransaction = new StellarSdk.TransactionBuilder(senderAccount, {
  fee: StellarSdk.BASE_FEE,
  networkPassphrase: StellarSdk.Networks.TESTNET,
})
  .addOperation(
    StellarSdk.Operation.payment({
      destination: destinationKeypair.publicKey(),
      asset: StellarSdk.Asset.native(),
      amount: "100",
      source: senderKeypair.publicKey(),
    })
  )
  .setTimeout(30)
  .build();
```

</details>
<br>
4. Sign the inner transaction, but don’t submit it to the network just yet

<details>
<summary>See it in code</summary>
<br>
Sign the inner transaction, but don’t submit it to the network just yet

```javascript
// Sign the inner transaction using the sender keypair. But, we will not be directly submitting this inner transaction on its own.
innerTransaction.sign(senderKeypair);
console.log("Inner transaction has been signed.");
```

</details>
<br>

5. Click the XDR located in the box above the three bottom purple buttons and copy the string
6. In a new Laboratory window, create a new transaction with Fee-Bump as the Transaction Type and the Quest Keypair’s public key as the Source Account
7. Paste the inner transaction XDR in the appropriate box

Note: now that you’ve done it the hard way, there’s a button in the Laboratory that makes life a little easier. You can click the “Wrap with Fee Bump” button after signing the transaction to navigate directly to the Fee-Bump Transaction

8. Sign and submit the transaction to the network

<details>
<summary>See it in code</summary>
<br>
Create, sign, and submit the separate fee-bump transaction, paying the fee for the inner transaction

```javascript
// Build the fee-bump transaction.  We will use your Quest Account as the "channel account."
// It will be this account that pays the transaction fee and the sequence number.
const questKeypair = StellarSdk.Keypair.fromSecret(
  "SECRET_KEY_FOR_YOUR_QUEST_ACCOUNT"
);
let feeBumpTransaction =
  new StellarSdk.TransactionBuilder.buildFeeBumpTransaction(
    questKeypair,
    StellarSdk.BASE_FEE,
    innerTransaction,
    StellarSdk.Networks.TESTNET
  );

// Sign the fee-bump transaction using the channel account keypair.
feeBumpTransaction.sign(questKeypair);
console.log("Fee-bump transaction has been signed.");

// Finally, submit the fee-bump transaction to the testnet.
try {
  let response = await server.submitTransaction(feeBumpTransaction);
  console.log(
    `Fee-bump transaction was successfully submitted.\nFee-bump transaction hash: ${response.fee_bump_transaction.hash}\nInner transaction hash: ${response.inner_transaction.hash}`
  );
} catch (error) {
  console.log(
    `${error}. More details:\n${JSON.stringify(error.response.data)}`
  );
}
```

</details>
<br>

9. Let’s see if everything worked- go to [stellar.expert](https://stellar.expert/explorer/testnet) and search for the Sending Account’s public key
10. Click on the latest operation
11. You should see the Fee Source Account is now your Quest Account. That means your Sender Account paid no fees for this transaction to make it onto the ledger
12. Great work- click Verify and claim your badge!

<br>

### Full Code

<details>
<summary>See the full code implementation</summary>
<br>

```javascript
(async () => {
  // include the StellarSDK
  const StellarSdk = require("stellar-sdk");
  const fetch = require("node-fetch");

  // Generate two Keypairs: a sender, and a destination.
  const senderKeypair = StellarSdk.Keypair.random();
  const destinationKeypair = StellarSdk.Keypair.random();

  // Optional: Log the keypair details if you want to save the information for later.
  console.log(`Sender Public Key: ${senderKeypair.publicKey()}`);
  console.log(`Sender Secret Key: ${senderKeypair.secret()}`);
  console.log(`Desintation Public Key: ${destinationKeypair.publicKey()}`);
  console.log(`Destination Secret Key: ${destinationKeypair.secret()}`);

  await Promise.all(
    [senderKeypair, destinationKeypair].map(async (kp) => {
      // Set up the Friendbot URL endpoints.
      const friendbotUrl = `https://friendbot.stellar.org?addr=${kp.publicKey()}`;
      let response = await fetch(friendbotUrl);

      // // Optional Looking at the responses from fetch.
      // let json = await response.json()
      // console.log(json)

      // Check that the response is OK, and give a confirmation message.
      if (response.ok) {
        console.log(`Account ${kp.publicKey()} successfully funded.`);
      } else {
        console.log(`Something went wrong funding account: ${kp.publicKey()}`);
      }
    })
  );

  // Connect to the testnet with the StellarSdk.
  const server = new StellarSdk.Server("https://horizon-testnet.stellar.org");
  const senderAccount = await server.loadAccount(senderKeypair.publicKey());

  // Build the inner transaction. This will be the transaction where the payment is actually made.
  let innerTransaction = new StellarSdk.TransactionBuilder(senderAccount, {
    fee: StellarSdk.BASE_FEE,
    networkPassphrase: StellarSdk.Networks.TESTNET,
  })
    .addOperation(
      StellarSdk.Operation.payment({
        destination: destinationKeypair.publicKey(),
        asset: StellarSdk.Asset.native(),
        amount: "100",
        source: senderKeypair.publicKey(),
      })
    )
    .setTimeout(30)
    .build();

  // Sign the inner transaction using the sender keypair. But, we will not be directly submitting this inner transaction on its own.
  innerTransaction.sign(senderKeypair);
  console.log("Inner transaction has been signed.");

  // Build the fee-bump transaction.  We will use your Quest Account as the "channel account."
  // It will be this account that pays the transaction fee and the sequence number.
  const questKeypair = StellarSdk.Keypair.fromSecret(
    "SECRET_KEY_FOR_YOUR_QUEST_ACCOUNT"
  );
  let feeBumpTransaction =
    new StellarSdk.TransactionBuilder.buildFeeBumpTransaction(
      questKeypair,
      StellarSdk.BASE_FEE,
      innerTransaction,
      StellarSdk.Networks.TESTNET
    );

  // Sign the fee-bump transaction using the channel account keypair.
  feeBumpTransaction.sign(questKeypair);
  console.log("Fee-bump transaction has been signed.");

  // Finally, submit the fee-bump transaction to the testnet.
  try {
    let response = await server.submitTransaction(feeBumpTransaction);
    console.log(
      `Fee-bump transaction was successfully submitted.\nFee-bump transaction hash: ${response.fee_bump_transaction.hash}\nInner transaction hash: ${response.inner_transaction.hash}`
    );
  } catch (error) {
    console.log(
      `${error}. More details:\n${JSON.stringify(error.response.data)}`
    );
  }
})();
```

</details>
<br>
