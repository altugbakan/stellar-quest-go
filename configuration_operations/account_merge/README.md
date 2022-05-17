# Accout Merge
### Delete an account by transferring its XLM balance to another account
<br>
Before we progress too far into our Stellar operational journey, we should talk about how to delete accounts.

**In this quest, you’ll delete an account by transferring all of its native XLM balance to a destination account using the Account Merge operation.**

Every Stellar account must maintain a minimum balance of XLM, which is calculated using base reserves (each base reserve is 0.5 XLM). A Stellar account must maintain two base reserves (1 XLM) to exist on the ledger — so transferring this XLM balance removes the account from the ledger.

Let’s get started.

1. Ensure the Quest Keypair is funded
2. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
3. Go to the Build Transaction tab in the main navigation bar
4. Fill out the necessary information with the Quest Keypair public key as the source account
5. Select Account Merge as the operation
6. Input an existing account’s public key as the destination

Note that it’s not just the XLM balance that needs to be removed from an account to merge it. You also need to remove all subentries: trustlines need to be emptied and removed, additional signers need to be dropped, open offers need to be closed, and data entries must be deleted.

These are irrelevant for this quest as our Quest Keypair doesn’t have any additional subentries, but it’s important to note for future use.

7. Sign and submit the transaction.
8. If your transaction was successful, click the Verify button on the Stellar Quest Learn page, claim your badge, and let’s move on to the next quest!