# Clawbacks

### Burn an amount of a specific asset from a receiving account

<br>
Asset clawbacks allow an asset issuer to burn a specific amount of a clawback-enabled asset from a trustline or claimable balance, effectively destroying it and removing it from a recipient’s balance.

Clawbacks were designed to allow asset issuers to meet securities regulations, and many jurisdictions require asset issuers (or designated transfer agents) to have the ability to revoke assets under certain circumstances. This can happen in the event of a mistaken or fraudulent transaction or other regulatory action regarding a specific person or asset. Asset clawbacks make it easier to issue specific types of regulated assets, such as money market funds, bonds, and equities on Stellar.

Note that you cannot enable clawbacks on assets retrospectively. All clawback-enabled assets are clearly designated by account flags that are visible on the public ledger.

There are three primary steps to set up, issue, and claw back a Stellar asset.

- Set up the issuing account to be clawback-enabled
- Issue the asset
- Claw all or a portion of the asset back, burning it in the process

Let’s talk about the clawback operations involved in this process.

### Operations

#### Set Options

The issuer sets up their account to be clawback enabled. This account-level flag causes every subsequent trustline established from the account to be clawback enabled.

If an issuing account wants to set the authorization clawback enabled flag, it must also have the authorization revocable flag set. This allows an asset issuer to claw back balances locked up in offers by first revoking authorization from a trustline, which pulls all offers involving that trustline from the ledger. The issuer can then perform the clawback.

#### Clawback

The issuing account uses this operation to claw back some or all of an issued asset. Once an account holds a particular asset for which clawbacks have been enabled, the issuing account can claw it back, burning it. You need to specify the asset, the amount to claw back, and the account from which you’re clawing back the asset.

**In this quest, you will make a payment of a clawback-enabled asset with your Quest Account as the source account. Then issue a clawback operation to claw some or all of the asset back.**

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Create a new account on the testnet or use an account you’ve already created, we’ll call this Account 2
3. Build a transaction that sets the clawback-enabled flag on the Quest Account
4. Sign and submit the transaction to the testnet
5. Build another transaction that creates a trustline and issues a clawback-enabled asset from the Quest Account to Account 2
6. Sign and submit the transaction to the testnet
7. Build a third transaction that claws back a portion or all of the asset from Account 2
8. Sign and submit that transaction to the testnet
9. If all goes well, verify your quest and claim your badge!

Note that you can remove clawback capabilities on individual trustlines with the Set Trust Line Flags operation. You can only clear a flag, not set it, so clearing a clawback flag on a trustline is irreversible.

In addition, you have the ability to claw back claimable balances. You can only claw back the entire claimable balance, not a portion.
