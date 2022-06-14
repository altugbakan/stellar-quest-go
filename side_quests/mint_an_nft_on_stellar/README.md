# Side Quest - Mint an NFT on Stellar
### Publish an NFT on the Stellar blockchain
<br>

Welcome, quester! Time to explore the world of NFTs on Stellar. Yes, it’s true! You can mint NFTs on the network! And it’s actually not too complicated, as NFTs behave pretty much like any other asset.

Let’s dive in.

**In this quest, we’re going to set up an issuing account to host the metadata for an NFT and then mint the NFT by sending it to the Quest Account (minter account).**

We’ll walk through this quest using the Laboratory and with code using the Stellar SDK. Expand each step to view the code portion. Knowledge of NodeJS is required for the SDK sections.

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test) and create a new account on the testnet, this account will be our issuing account- remember to save the public and secret keys somewhere on your computer

<details>
<summary>See it in code</summary>
<br>

```javascript
// Include the StellarSDK.
const StellarSdk = require('stellar-sdk');
const fetch = require('node-fetch');

// Generate two Keypairs: one for issuing the NFT, and one for receiving it.
const issuerKeypair = StellarSdk.Keypair.random();
const receiverKeypair = StellarSdk.Keypair.fromSecret('SECRET_KEY_FOR_YOUR_QUEST_ACCOUNT');

// Optional: Log the keypair details if you want to save the information for later.
console.log(`Issuer Public Key: ${issuerKeypair.publicKey()}`);
console.log(`Issuer Secret Key: ${issuerKeypair.secret()}`);
console.log(`Receiver Public Key: ${receiverKeypair.publicKey()}`);
console.log(`Receiver Secret Key: ${receiverKeypair.secret()}`);
```
Fund your accounts with Friendbot

```javascript
// Fund both accounts using Friendbot. We're performing the fetch operation, and ensuring the response comes back "OK".
await Promise.all([ issuerKeypair, receiverKeypair ].map(async (kp) => {
  // Set up the Friendbot URL endpoints
  const friendbotUrl = `https://friendbot.stellar.org?addr=${kp.publicKey()}`;
  let response = await fetch(friendbotUrl);

  // // Optional: Looking at the responses from fetch
  // let json = await response.json();
  // console.log(json);

  // Check that the response is OK, and give a confirmation message.
  if (response.ok) {
    console.log(`Account ${kp.publicKey()} successfully funded.`);
  } else {
    console.log(`Something went wrong funding account: ${kp.publicKey}.`);
  }
}))
```
</details>
<br>

Next, we’re going to get your NFT hosted on IPFS. Note that you don’t need to use IPFS, you can use any other storage mechanic, but we encourage using IPFS because it’s decentralized.

2. Navigate to [NFT.storage](https://nft.storage/), a free service where you can host your files on IPFS
3. Upload whatever image you’d like to use for your NFT
4. Head on over to [SEP-0039 Interoperability Recommendations for NFTs](https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0039.md)- this documentation dictates the best practices for NFTs on Stellar
5. Scroll down to the “Describing Your NFT” section and copy and paste the JSON file into your preferred text editor
6. Input your NFT information in the following fields:
- Name: the name of your NFT, go ahead and make something up
- Description: describe your NFT
- URL: copy and paste the URL from your image uploaded on NFT.storage
- Issuer: input the issuing account’s public key
- Code: the Stellar asset code for your NFT (remember: this cannot be any longer than 12 characters and must be made from characters in the set [a-z][A-Z][0-9])
7. Save the completed JSON file onto your computer
8. Upload the JSON file to NFT.storage

<details>
<summary>See it in code</summary>
<br>

```javascript
const { NFTStorage, Blob } = require('nft.storage');
const fs = require('fs');

// Create the Asset so we can issue it on the network.
const nftAsset = new StellarSdk.Asset('vvNFT', issuerKeypair.publicKey());

// Store the Image and metadata using nft.storage
const NFT_STORAGE_TOKEN = 'your_api_key'; // Get this from https://nft.storage/manage
const IMAGE_PATH = '/path/to/your/image.jpg';
const client = new NFTStorage({ token: NFT_STORAGE_TOKEN });

const imageCID = await client.storeBlob(new Blob([fs.readFileSync(IMAGE_PATH)]));
console.log(`imageCID: ${imageCID}`);

const metadata = {
  name: "Very Valuable NFT",
  description: "This is the most valuable NFT available on any blockchain. Ever.",
  url: `ipfs://${imageCID}`,
  issuer: nftAsset.getIssuer(),
  code: nftAsset.getCode()
};
const metadataCID = await client.storeBlob(new Blob([JSON.stringify(metadata)]));
console.log(`metadataCID: ${metadataCID}`);
```
</details>
<br>

13. Add a Change Trust operation with your Stellar Asset Code and Issuer Account filled in- ensure the Quest Account is the Source Account for this operation

<details>
<summary>See it in code</summary>
<br>

```javascript
// Perform a `changeTrust` operation to create a trustline for the receiver account.
transaction = transaction
  .addOperation(StellarSdk.Operation.changeTrust({
    asset: nftAsset,
    limit: '0.0000001',
    source: receiverKeypair.publicKey(),
  }));
```
</details>
<br>

14. Add a Payment operation with the Quest Account (minter account) used as the Destination and the Asset Code and Issuing Account filled in accordingly

15. Mint just one of the NFT by inputting one stroop as the Amount (0.0000001) - a stroop is the smallest unit of an asset on Stellar, one ten-millionth of the asset

Note that you can mint more than one of the NFT by increasing the stroops. For example, to issue 100 of the NFT, you’d change the amount to 0.0000100.

<details>
<summary>See it in code</summary>
<br>

```javascript
// Add a `payment` operation to send the NFT to the receiving account.
transaction = transaction
.addOperation(StellarSdk.Operation.payment({
  destination: receiverKeypair.publicKey(),
  asset: nftAsset,
  amount: "0.0000001",
  source: issuerKeypair.publicKey(),
}));
```
</details>
<br>

16. Sign and submit the transaction! Congratulations, you’ve minted an NFT!

<details>
<summary>See it in code</summary>
<br>

```javascript
// setTimeout is required for a transaction, and it also must be built.
transaction = transaction
  .setTimeout(30)
  .build();

// Sign the transaction with the necessary keypairs.
transaction.sign(issuerKeypair);
transaction.sign(receiverKeypair);

try {
  await server.submitTransaction(transaction);
  console.log('The asset has been issued to the receiver');
} catch (error) {
  console.log(`${error}. More details: \n${error.response.data}`);
}
```
</details>
<br>

17. Now you can verify the quest and claim your badge

### Additional Notes

You can control certain aspects of your NFT with the Set Options operation. We won’t use these in this quest, but it’s important to know for future NFT fun.

#### Set Flags
Configure features like royalty payments and control access to your NFT by setting flags on the issuing account.

[Learn more in this SQ Learn quest](https://quest.stellar.org/learn/series/2/quest/5)

#### Master Weight
Lock the issuing account by setting the master weight to 0. This will ensure the following:

- The issuing account will never be able to issue additional units of the NFT
- Data entries on the issuing account can never be modified
- The Account Merge operation cannot be used, so the issuing account cannot be deleted
- Ownership of the NFT can never be revoked

Locking the account is irreversible, so take care when doing so.

[Learn more in this SQ Learn quest](https://quest.stellar.org/learn/series/2/quest/4)

<br>

### Full Code
<details>
<summary>See the full code implementation</summary>
<br>

```javascript
(async () => {
  // Include the StellarSDK and some other utilities.
  const StellarSdk = require('stellar-sdk');
  const fetch = require('node-fetch');
  const { NFTStorage, Blob } = require('nft.storage');
  const fs = require('fs');

  // Generate two Keypairs: one for issuing the NFT, and one for receiving it.
  const issuerKeypair = StellarSdk.Keypair.random();
  const receiverKeypair = StellarSdk.Keypair.fromSecret('SECRET_KEY_FOR_YOUR_QUEST_ACCOUNT');

  // Optional: Log the keypair details if you want to save the information for later.
  console.log(`Issuer Public Key: ${issuerKeypair.publicKey()}`);
  console.log(`Issuer Secret Key: ${issuerKeypair.secret()}`);
  console.log(`Receiver Public Key: ${receiverKeypair.publicKey()}`);
  console.log(`Receiver Secret Key: ${receiverKeypair.secret()}`);

  // Fund both accounts using Friendbot. We're performing the fetch operation, and ensuring the response comes back "OK".
  await Promise.all([ issuerKeypair, receiverKeypair ].map(async (kp) => {
    // Set up the Friendbot URL endpoints
    const friendbotUrl = `https://friendbot.stellar.org?addr=${kp.publicKey()}`;
    let response = await fetch(friendbotUrl);

    // // Optional: Looking at the responses from fetch
    // let json = await response.json();
    // console.log(json);

    // Check that the response is OK, and give a confirmation message.
    if (response.ok) {
      console.log(`Account ${kp.publicKey()} successfully funded.`);
    } else {
      console.log(`Something went wrong funding account: ${kp.publicKey}.`);
    }
  }))

  // Create the Asset so we can issue it on the network.
  const nftAsset = new StellarSdk.Asset('vvNFT', issuerKeypair.publicKey());

  // Store the Image and metadata using nft.storage
  const NFT_STORAGE_TOKEN = 'your_api_key'; // Get this from https://nft.storage/manage
  const IMAGE_PATH = '/path/to/your/image.jpg';
  const client = new NFTStorage({ token: NFT_STORAGE_TOKEN });

  const imageCID = await client.storeBlob(new Blob([fs.readFileSync(IMAGE_PATH)]));
  console.log(`imageCID: ${imageCID}`);

  const metadata = {
    name: "Very Valuable NFT",
    description: "This is the most valuable NFT available on any blockchain. Ever.",
    url: `ipfs://${imageCID}`,
    issuer: nftAsset.getIssuer(),
    code: nftAsset.getCode()
  };
  const metadataCID = await client.storeBlob(new Blob([JSON.stringify(metadata)]));
  console.log(`metadataCID: ${metadataCID}`);


  // Connect to the testnet with the StellarSdk.
  const server = new StellarSdk.Server('https://horizon-testnet.stellar.org');
  const account = await server.loadAccount(issuerKeypair.publicKey());

  // Build a transaction that mints the NFT.
  let transaction = new StellarSdk.TransactionBuilder(
    account, {
      fee: StellarSdk.BASE_FEE,
      networkPassphrase: StellarSdk.Networks.TESTNET
    })
    // Add the NFT metadata to the issuer account using a `manageData` operation.
    .addOperation(StellarSdk.Operation.manageData({
      name: 'ipfshash',
      value: metadataCID,
      source: issuerKeypair.publicKey(),
    }))
    // Perform a `changeTrust` operation to create a trustline for the receiver account.
    .addOperation(StellarSdk.Operation.changeTrust({
      asset: nftAsset,
      limit: '0.0000001',
      source: receiverKeypair.publicKey(),
    }))
    // Add a `payment` operation to send the NFT to the receiving account.
    .addOperation(StellarSdk.Operation.payment({
      destination: receiverKeypair.publicKey(),
      asset: nftAsset,
      amount: "0.0000001",
      source: issuerKeypair.publicKey(),
    }))
    // setTimeout is required for a transaction, and it also must be built.
    .setTimeout(30)
    .build();

  // Sign the transaction with the necessary keypairs.
  transaction.sign(issuerKeypair);
  transaction.sign(receiverKeypair);

  try {
    await server.submitTransaction(transaction);
    console.log('The asset has been issued to the receiver');
  } catch (error) {
    console.log(`${error}. More details: \n${error.response.data}`);
  }
})();
```
</details>
<br>