# Sponsorship

### Pay account subentry reserves with another account

<br>
Storing data on the Stellar blockchain sometimes requires you to also stake a small amount of XLM, which helps prevent spam and deter malicious behavior. This staked XLM is referred to as a base reserve, and each base reserve equals 0.5 XLM.

Typically, base reserves are staked automatically as you submit transactions to the network. Examples of data that require base reserves are:

**Subentries** - account data that requires one base reserve. Subentries can be trustlines, buy/sell offers on the decentralized exchange, additional signers, or manage data entries.

**Accounts** - every Stellar account requires two base reserves (1 XLM) to exist.

**Ledger entries** - ledger entries such as each claimant on a claimable balance require one base reserve.

While base reserve requirements exist for a good cause, they can raise the barrier to entry onto Stellar for traditional B2C users who need to learn how to acquire XLM to start an account. This can be tricky if users are unfamiliar with blockchain.

Enter sponsored reserves, the hero of this particular story. Sponsored reserves allow an account (the sponsoring account) to pay some or all base reserve requirements of another account (the sponsored account).

Sponsored reserves use three operations: begin sponsoring future reserves, end sponsoring future reserves, and revoke sponsorship.

**Begin sponsoring future reserves** - this operation initiates the sponsorship relationship and requires the sponsoring account’s signature.

**End sponsoring future reserves** - this operation allows the sponsored account to accept the sponsorship and ends the current sponsorship relationship. This requires the sponsored account’s signature, and the sponsored account must be the source account for this operation.

Note: Ongoing is-sponsoring-future-reserves-for relationships can’t exist on the network. To use sponsored reserves, you need to wrap a transaction’s operations in a delicious begin and end sponsoring future reserves sandwich. This means an account can sponsor a specific set of operations in a particular transaction, but is not giving blanket sponsorship for operations that take place outside of that transaction in the future.

**Revoke sponsorship** - allows the sponsoring account to transfer sponsorship of existing sponsored reserves. A revoke sponsorship operation needs to be between a set of begin and end sponsoring future reserves operations unless sponsorship is being transferred back to the originating account itself.

**In this quest, brave questers, you will set up your Quest Account to have a zero balance and to have its base reserves sponsored by another account.**

Let’s get into how this works.

**But wait!** Did you already hit that handy Fund button to fund your Quest Account? If so, we’ll do this quest a little differently. Navigate to part two in this instance.

Part 1: you didn’t fund your Quest Account with the Fund button (excellent patience!)

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Create a new account on the testnet- this will be your sponsoring account
3. Construct a transaction that funds and creates the Quest Account with the sponsoring account paying its base reserve requirements- this involves a few operations: begin sponsoring future reserves, create account, and end sponsoring future reserves

Hint: the sponsoring account will be the source account for the transaction, but pay attention to the source account for each operation

4. Sign and submit the transaction to the testnet
5. Verify the quest and claim your badge!

Part 2: you’re a bit overeager and funded your Quest Account with the Fund button

Since you already funded your Quest Account, your account is paying its own base reserves to exist. We need to transfer that sponsorship to another account and then get rid of all the XLM you got from friendbot.

1. You guessed it! Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Create a new account on the testnet- this will be your sponsoring account
3. Construct a transaction that transfers the sponsorship of the Quest Account’s base reserves and all available XLM to the sponsoring account- this involves a few operations: begin sponsoring future reserves, revoke sponsorship, end sponsoring future reserves, and a payment

Hint: the sponsoring account will be the source account for the transaction, but pay attention to the source account for each operation

4. Sign and submit the transaction to the testnet
5. Verify the quest and claim your badge!
