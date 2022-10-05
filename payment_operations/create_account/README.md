# Create Account

### Create an account on the Stellar network

<br>
If Stellar were a universe, and it is, accounts would be the planets, stars, moons, and asteroids within that fertile space.

Less abstractly, accounts are the central data structure in Stellar — they hold balances, sign transactions, and issue assets. Accounts can only exist with a valid keypair and the required minimum balance of lumens (XLM).

**In this quest, your challenge is to perform a `createAccount` operation using the Quest Keypair located in the box on the right-hand side of this screen.**

This quest has two components:

1. Create the Quest Account by funding the Quest Keypair with XLM from friendbot
2. Create a new account using the `createAccount` operation with the Quest Account as the source account

So let’s get started!

## Part 1: Create the Quest Account by funding the Quest Keypair with XLM from friendbot

In this first part, we will create the Quest Account on the Stellar test network by funding the Quest Keypair with lumens (XLM) from friendbot, a bot that funds accounts on the testnet with 10,000 fake XLM.

1. Go to the [Stellar Laboratory](https://laboratory.stellar.org/#account-creator?network=test) and select the Create Account tab in the main navigation bar
2. Copy and paste the Quest Keypair’s public key into the friendbot section
3. Click the Get Test Network Lumens button

Note: you will get a new Keypair for every Quest.

Incredible! You’ve begun interacting with the Stellar blockchain. No, seriously, you did! Let’s make sure the funding request went through, and the account exists on the network.

4. Navigate to [stellar.expert](https://stellar.expert/explorer/testnet)
5. Ensure you’re looking at the testnet, not the public network
6. Input the Quest Account’s public key (also called account ID) into the search bar at the top of the page and hit enter
7. You should see your account exists and is funded with 10,000 XLM

That was exciting! However, on the public network, we don’t have friendbot handing out free XLM to anyone that asks. So let’s look at things more realistically.

### Part 2: Create a new account using the createAccount operation with the Quest Account as the source account

In this second part, we will create a brand-new account by generating another keypair and funding it with XLM from the Quest Account. This will use the createAccount operation. Since this is your first transaction submission attempt on Stellar, we will walk through the process in detail. In later quests, we’ll give you less instruction to ensure you are doing the necessary work to learn and retain the information.

Let’s get into it!

1. Go to the [Stellar Laboratory](https://laboratory.stellar.org/#account-creator?network=test) and select the Create Account tab in the main navigation bar
2. Click the Generate Keypair button
3. Save the public and secret keys somewhere on your computer
4. Do not fund your account with friendbot, we will fund this account with the Quest Account
5. Navigate to the Build Transaction tab in the main navigation bar — don’t freak out! We’ll walk through the inputs below, so you have a solid understanding of what’s happening

#### Transaction Type

6. Set to Transaction

Fee Bumps are an advanced feature that we'll cover later.

#### Source Account

7. Input the public key from the Quest Account

The Source Account field serves at least three purposes, explained below. There are some exceptions where you can maneuver around these defaults, but you will generally find them to be true.

- Every transaction needs at least one signature to be valid. Transactions will always need the source account’s signature.
- Every transaction submitted to the network consumes fees, and that fee will be taken from the source account.
- Every transaction increases the source account’s sequence number by one.

#### Transaction Sequence Number

8. Select the Fetch Next Sequence Number for Account Starting with "G…" button

Every transaction increases the source account’s sequence number by one. A sequence number is like a page number and prevents the same transaction from being submitted to the network more than once. When you input the Source Account above, you should see a button appear that you can select to fetch the account’s next sequence number.

#### Base Fee

9. Set to 100 stroops

Stellar requires a small fee for all transactions. Currently, the network minimum is 100 stroops, but you can set your fee to anything above that.

Stroop: the smallest unit of a lumen, one ten-millionth of a lumen (.0000001 XLM)

The minimum fee required for a given transaction equals the number of operations in the transaction multiplied by the base fee for the given ledger. (transaction fee = # of operations \* base fee)

When you input a base fee price, you specify the maximum amount that you’re willing to pay per operation. That does not necessarily mean that you’ll pay that amount, you will only be charged the lowest amount needed for your transaction to make it to the ledger. If network traffic is light, and the number of submitted operations is below the network ledger limit, you will only pay the network minimum (currently 100 stroops). When network traffic exceeds the ledger limit, the network enters into surge pricing mode, and your fee becomes a max bid.

#### Memo

10. Leave blank

Input context for the transaction. This field is more informational than technical.

#### Time Bounds

11. Leave blank

Time bounds are optional but recommended, as they put a time limit on the transaction — so either the transaction makes it onto the ledger or it times out and fails, depending on your time parameters. You can also configure your transaction to not submit until a specified date and time.

#### Operation Type

12. Select Create Account

Transactions are made up of operations. On Stellar, transactions can contain up to 100 operations and are atomic. If one operation fails in a transaction, they all fail.

#### Destination

13. Input the second account’s public key

The destination account you’re sending XLM to.

#### Starting Balance

14. The Quest Account, which you're using as the source account for this operation, should have 10,000 XLM from friendbot, so select any amount less than this

How much XLM you’d like to send from the source account to this new destination account we’re creating.

#### Source Account

15. Leave blank

This operational source account field is optional because it assumes the source account for the transaction if left blank. You can specify a different source account for each operation if needed.

16. Scroll down and select the Sign in Transaction Signer button. This bundles the operation(s) into a transaction into an XDR and sends it along for signing.

Oh, geez, more forms and fields. Not to worry, let’s just talk it out.

#### Signing For

There are two primary official networks on Stellar, the public network and the test network. This field tells you what network you’re submitting to.

#### Transaction Envelope XDR

The XDR encoded transaction containing all our transaction information in a neat package. If you needed to store or forward the transaction elsewhere, this is likely what you would save to pass around.

#### Transaction Hash

The SHA256 hash of the XDR above. Fun fact: this is actually what gets signed during the signing process. If you ever decide to play around with hardware wallets and need to move around hashes for signing, this information will be helpful.

#### Source Account

The source account for the transaction. Crazy, I know.

#### Sequence Number

The sequence number that will be consumed when this transaction is submitted, never to be used again.

#### Transaction Fee (stroops)

The potential fee to be paid by the source account.

#### Number of Operations

The number of operations in the transaction.

#### Number of Existing Signatures

The number of signatures encoded within the XDR. A transaction can be signed and then passed on to other parties for further signing. For example, suppose a transaction with two or more operations has different source accounts, requiring more than one signature. You can add all the signatures at once or pass the partially signed XDR on for further signing before submission.

#### Signatures

There are four different options here, and they all represent different ways to accomplish the same thing: sign the transaction. Often, you will just use the Add Signer field(s), but if you have a hardware wallet, you may use the BIP Path field, and if you use a software wallet like Freighter or Albedo, you may also make use of those buttons. For now, though, we’ll just be using Add Signer.

#### Add Signer

17. Manually add the secret key of the source account to the Add Signer field. There are occasions where you’ll need more than one signature, which we’ll get into in a later quest.

Note: inputting your secret key into a field is uniquely accepted in the Laboratory. Copying and pasting your secret key is normally not a good idea as it’s easy to accidentally paste your secret key into the wrong field which could be devastating.

#### BIP Path

Weird name, bro. This is for hardware wallets like Ledger and Trezor. If you don’t know what those are, don’t worry about this field.

#### Freighter

Freighter is an SDF sanctioned non-custodial extension wallet. It’s pretty good, you should try it sometime.

#### Albedo

Another great wallet built by the same team that brings us our premier blockchain explorer, Stellar Expert, maintains this additional non-custodial signing option.

18. Click the Submit in Transaction Submitter button. The next screen will show the decoded XDR information. Double-check that everything here looks accurate before submitting.
19. Click the Submit Transaction button.
20. Hopefully, everything went well and you see a successful transaction! In the Stellar Quest screen, click the Verify button to the right to see if you passed, then collect your very first Stellar Quest Learn NFT badge.
21. If the quest failed, double-check your transaction to ensure everything is correct. If you have questions, head over to our Stellar Quest Discord to ask the community for help!
