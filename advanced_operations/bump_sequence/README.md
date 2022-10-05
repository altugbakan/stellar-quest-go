# Bump Sequence

### Increase the sequence number of an account to a new given sequence number

<br>
Every Stellar account maintains a sequence number that increments by one whenever it submits a transaction to the network. These sequence numbers are used to identify and verify transactions with the same source account and prevent double-spending. Valid transactions have a sequence number that is the source account’s current sequence number plus one.

Note: transactions can consume sequence numbers even if they subsequently fail to validate and don’t make it to the ledger. For example, if a transaction is submitted for a 100 XLM payment and your account only has 99 XLM, the transaction will use the sequence number but will still ultimately fail.

Alright. So let’s say you want to skip a few sequence numbers. Well, you’ve come to the right place! There’s a fun little Stellar operation called Bump Sequence that allows you to skip ahead to a desired future sequence number.

How is this useful? Good question. Let’s dig in.

Again, for a transaction to be valid, its sequence number must be the source account’s current sequence number plus one. With this in mind, you can configure multiple possible scenarios by building a series of transactions with different future sequence numbers, then use the bump sequence operation to choose which scenario you’d like to be valid.

For example, let’s say you’re betting on turtle racing. Your account’s current sequence number is 0, and you want to build two possible future scenarios: (1) Speedy wins, and you need to pay your friend Buster - set this scenario to sequence number 5, and (2) Snappy wins, and you need to pay your friend Wade - set this scenario to sequence number 10. You can then choose the valid scenario by bumping the account’s sequence number to either 4 or 9, depending on who wins the race. And then maybe reconsider gambling, as you’re not very good at it.

Note that by choosing sequence number 4, there is a possible future where 10 is still submittable, and by choosing 9, you make 5 impossible ever to execute.

Combine this with pre-authorized transactions, and you can accomplish some pretty interesting, though rudimentary, smart contract scenarios.

**In this quest, you will create and submit a bump sequence transaction that increases the sequence number of the Quest Account, then submit a second transaction consuming the new bumped-to number.**

Warning! Be careful when setting a new sequence number - accounts have a maximum sequence number and, if reached, will prevent your account from ever being a valid source account again. The maximum sequence number is the largest possible 64-bit integer: (2^63) - 1 = 9,223,372,036,854,775,807

Let’s get started.

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Construct a bump sequence transaction with the Quest Account as the source account
3. Sign and submit this bump sequence transaction to the testnet
4. Construct a new transaction starting from the bumped-to sequence number with the Quest Account as the source account. You can use any operation you’d like for this second transaction.
5. Sign and submit _that_ transaction to the testnet
6. Verify the quest and claim your badge!
