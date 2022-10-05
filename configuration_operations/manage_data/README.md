# Manage Data

### Add a data entry to an account

<br>
Accounts are the backbone of Stellar applications and have a part in any action you take on the network. One class of actions, and it’s an important one, involves setting, modifying, or deleting data entries within an account’s metadata. Data entries are a key/value pair and attach arbitrary, application-specific data to an account on the Stellar network. Each data entry increases the account’s minimum balance by one base reserve.

**In this quest, you will add a data entry to the Quest Account using the Manage Data operation.**

1. You guessed it, navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Build a transaction with Manage Data as the operation — let’s talk about what’s here:
   **Entry Name**

The string value name of your key/value pair; must be less than or equal to 64 characters

**Entry Value**

A 64-byte binary value; leave this blank to delete the data entry

In Horizon the Entry Value is base64 encoded. To get to the actual value you’ll need to convert these values from base64.

**Source Account**

The account you want to add, remove, or edit data entries for

3. Once you’ve input data for the Entry Name and Entry Value, sign and submit the transaction to the testnet
4. If successful, click that Verify button on the Stellar Quest Learn page and claim your spiffy new badge
5. Remember, you can always view activity on an account by visiting [Stellar Expert](https://stellar.expert/explorer/testnet) and searching for your account’s public key
