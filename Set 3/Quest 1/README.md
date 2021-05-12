# Make use of a sequence number bump operation in a transaction
In this inaugural challenge you must submit a transaction from your account, making use of the sequence number bump operation.

What good is the sequence number bump operation you ask? While it may not be a heavily utilized operation within an account's lifecycle it's an incredibly useful op when dealing with smart contracts, particularly around pre-signed transactions.

Imagine a scenario where you have two potential outcomes but only one of them should actually execute. Rather than having both transactions compete for the same sequence number you can control the outcome by bumping the sequence number to support whichever of the two scenarios you wish.

With functionality like this you can now block transaction submission both by time and by sequence. Control all the things!

Don't forget to check the [clue](https://horizon.stellar.org/transactions/073cde1fab7d28d3e322e5f61ac385c37859c3e82b17b6c92a9f0444420336cb).
