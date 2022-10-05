# Claimable Balances

### Send an amount of an asset from a source account to be accepted by a receiving account

<br>
Payments and offers are two wonderful ways to transfer value on the Stellar network, but there is yet another powerful method: claimable balances.

Claimable balances aim to solve the age-old problem of needing to open trustlines before making non-XLM payments. To do so, claimable balances split a payment into two parts:

Part 1: the sending account creates a payment using the Create Claimable Balance operation, essentially staking it on-chain

Part 2: the destination account(s), or claimant(s), accepts the claimable balance using the Claim Claimable Balance operation

### Claimable Balance Operations

#### Create Claimable Balance

The first operation we’ll use is Create Claimable Balance, which allows you to lock up an asset for one or more claimants to accept once a specific set of predicates (or conditions) is satisfied.

#### Asset

The asset you’d like to lock up. Nothing fancy here, just ensure the source account owns the desired asset and holds enough of a balance. Keep in mind, once you create the claimable balance, you’ll be unable to access those funds ever again unless you’ve set yourself up as one of the claimants.

#### Amount

The amount of the asset you want to lock up in the claimable balance.

#### Claimants

One of the fanciest pairs of pants that this feature holds in its armoire is its ability to list up to ten claimants or recipients as the destination for the claimable balance. In addition to the actual destination address(s), you can also configure a conditional predicate under which that address can accept the claimable balance.

A claimable balance can only be claimed once. If you have ten claimants listed, whoever claims the claimable balance first will receive the entire balance, and the claimable balance will be removed from the network.

#### Destination

The destination address that has permission to claim this claimable balance, as long as the predicate conditions are satisfied.

#### Predicate

The conditions under which the destination address may successfully claim the claimable balance. Currently, there are only two predicates to choose from.

**Unconditional**

The claimant’s address can claim the claimable balance at any time.

**Conditional**

This gives us a few time-based conditions to work with. Think carefully while configuring these predicates, as they can be a bit tricky to reason through.

Relative time - a relative span of time for when the claimable balance can be claimed. This value represents the number of seconds since the close time of the ledger containing the transaction that created the claimable balance (or, more simply, the time that the balance was created).

Absolute time - a deadline for when the claimable balance can be claimed. This value represents a specific UNIX epoch timestamp (the number of milliseconds since 00:00:00 UTC on 1 January 1970).

AND - multiple conditions must be met for the claimant to claim the claimable balance. For example: X minutes must have passed since the creation of the claimable balance, AND you need to claim it before Y date.

OR - the claimant can claim the claimable balance under two circumstances. For example: claim the claimable balance after X minutes have passed OR after Y date.

NOT - the claimant cannot claim the claimable balance in a certain circumstance. For example: the claimant cannot claim the claimable balance until X amount of time has passed.

### Claim Claimable Balance

The second operation is Claim Claimable Balance, where a claimant attempts to (you guessed it) claim the claimable balance. This operation only has one unique parameter. Easy right? Yes! Well, hold on, it’s kinda easy. The tricky part is determining the claimable balance ID. It’s not returned with the transaction submission response from the creation operation, so you’ll need to use the claimable balances endpoint to search for the claimable balance ID using some search criteria, such as the source account or claimant address.

#### Claimable Balance ID

Once you’ve found the claimable balance ID, you can use it in the Claim Claimable Balance operation to attempt to claim it.

**Alright, questers! It’s time for your task. In this quest, you will create a claimable balance with your Quest Account with at least one claimant where the predicate conditions are such that the claimant is blocked from claiming the balance for exactly five minutes after successfully creating the claimable balance.**

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Create a new account on the testnet or use one you’ve already created, we’ll call this Account 2
3. With your Quest Account as the source account, create a claimable balance with Account 2 as a claimant and a predicate stating they cannot claim the claimable balance until five minutes have passed. Sign and submit this transaction to the testnet.
4. Once your five minutes are up, construct a new transaction where you claim the claimable balance with Account 2. Sign and submit this transaction to the testnet, as well.

Hint: look up the claimable balance ID by using the Explore Endpoints section of the Laboratory

5. Verify the quest and get your badge!
