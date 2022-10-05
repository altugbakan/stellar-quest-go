# Set Options - Weights, Thresholds, and Signers

### Set the master key’s weight, determine the operation threshold, and manage signers

<br>
Imagine you’re an adventurer. You and two companions have battled foes across land and sea to collect three keys, hoping they will unlock a vault door guarding a vast treasure. After much tribulation, your squad stands at the door, bewildered, on the brink of despair — the door is sealed with not three locks but five. You cast your eyes down in defeat. Wait, what is that? You see the number “2” etched into the key in your right palm. Curious, you implore your allies to look at their own keys. Ah hah! One has a “2” and the other a “1”. Your brilliant mind starts churning, and you realize that two of the keys fit two of the locks. Shaking, you insert the keys into their locks and turn. At the last spin of the final key, the vault door creaks open, revealing the priceless treasure you’ve journeyed so far to acquire. Ring Pops! Thousands upon thousands of sparkling, crystalline, tantalizing Ring Pops! What a journey. What a puzzle.

Alright, let’s circle it back. The keys in this riveting story work like signatures on the Stellar network. Signatures are needed to authorize transactions. You sign transactions with the master key, which is the private key that corresponds to the account’s public key. Some transactions require additional or alternate signatures, which is called multisignature (or just multisig).

**In this quest, we are going to set signature weights and operation thresholds, and add additional signers to the Quest Account using the Set Options operation. Then submit a successful multisignature transaction.**

Let’s get to it.

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Build a transaction with Set Options as the operation. Let’s talk about what we’re working with here:

**Master Weight**

The Master Weight field allows you to modify the master key’s signature weight. The weight is a numerical value that denotes how much power the signature has. The default master weight is 1.

3. Set the Master Weight to 1

**Thresholds (Low, Medium, High)**

Every operation falls into one of three predetermined threshold categories: low, medium, or high (check out operation thresholds in the [List of Operations](https://developers.stellar.org/docs/start/list-of-operations/) section of the docs), and you can set each threshold category’s weight between 0 and 255.

For a transaction to be successful, the combined weight of the signers must be greater than or equal to the threshold weight of each operation.

By default, all thresholds are set at 0. This doesn’t mean that any account can sign for operations, but ensures the default master key can.

**Note:** Always be careful setting thresholds and weights. If you set threshold values above 1 without adding more signers or increasing the master weight, you will lock yourself out of the account, and you will never be able to execute any operations ever again.

4. Set the Low, Medium, and High thresholds to 5

**Signer Type**

This field allows you to add or remove additional signers to accounts and configure their weights.

There are currently three signer types, but for this quest, we’re only interested in configuring additional Stellar keys or `Ed25519 Public Keys`.

5. Create two more Set Options operations, each adding one additional signer with weights of 2

With our thresholds all set at 5, our master weight at 1, and our additional signers’ weight at 2 each, we’d need all three signers to approve the operations for the transaction to be successful.

Master key (weight 1) + Signer key A (weight 2) + Signer key B (weight 2) = 5

6. Sign and submit the transaction
7. Submit another transaction to the network that is signed by all three signers
8. If your transaction was successful, congratulations! Go ahead and claim that badge.
