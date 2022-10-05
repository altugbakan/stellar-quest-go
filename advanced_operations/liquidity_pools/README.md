# Liquidity Pools

### Deposit or withdraw assets of a liquidity pool in exchange for pool shares

<br>
Every asset issued on Stellar starts with zero liquidity. “Liquidity” is just a fancy finance word that means how much of an asset in a market is available to be bought or sold. To create and sustain a market, people must be willing to put capital into the network and use that capital to facilitate asset conversion. Originally on Stellar, there was only one way to do that: by creating orders on the order books. However, the desire to improve liquidity and make cross-border payments faster, cheaper, and more user-friendly led to the proposal to add automated market maker (AMM) functionality to Stellar.

With the implementation of Protocol 18 and CAP-0038 (Core Advancement Proposals, or CAPS, are ideas/suggestions that guide the development and direction of the Stellar network), numerous new features and operations were added to the Stellar network that enable on-chain, constant product AMMs via the liquidity pool (LP) feature.

This quest, my friends, is a capstone quest. A real headbanger. So buckle up! Pour that coffee. And let your friends know you won’t be seeing them for a while.

**This quest has three objectives:**

- **Open a new LP with the Quest Account by making an initial deposit**
- **Use the LP with a path payment operation from another account**
- **Finally, close your position in the LP with the Quest Account and receive any fee rewards**

Let’s start by getting into what liquidity pools are and how they work.

LPs exist as deposits of two different assets. Accounts can deposit these assets into the LP and, in return, receive pool shares representing their ownership of that asset. When an account withdraws assets from the LP, it will receive a portion of the fees generated as transactions move through it.

To participate in an LP, you must have a trustline established with three different assets: both the reserve assets (unless one of them is XLM) and the liquidity pool share itself.

Let’s get to it.

Our first task is to establish a trustline to the new LP, deposit the two reserve assets, and, in turn, receive the appropriate pool shares.

1. Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Construct a transaction with Change Trust as the operation and the Quest Account as the source account- let’s talk about the fields here:

**Asset**

Select liquidity pool shares.

**Asset A**

This is the first reserve asset in the liquidity pool. For this quest, we’ll use the native token.

**Asset B**

The second reserve asset in the liquidity pool. For this quest, we’ll use our own self-issued NOODLE token (the Quest Account will be the asset issuer). ‘Cause what’s a pool without some noodles?

**Fee**

AMMs charge fees on every trade- and fees (for now) are statically set at 30 bps (**b**asis **p**oint**s**), which means a fee of 0.3% is required whenever the pool is being used during a trade.

3. Sign and submit this transaction to the testnet

Now that we’ve established our LP share trustline, it’s time to finish creating the LP by depositing some reserve assets.

4. Construct a new transaction with the Quest Account as the source account and Liquidity Pool Deposit as the operation, and let’s talk about what in the heck is going on here:

**Liquidity Pool ID**

The LP trustline you just created carries this entry. Look up the Quest Account in the Explore Endpoints tab to find this ID. This LP ID is derived by hashing various bits of information about the LP. This means there can only be one single LP on the network for any given asset pair. For you, right now, it also means you don’t have to specify the particular Asset A or Asset B in your deposit operation.

**Max Amount A, Max Amount B, Min Price, and Max Price**

There’s a lot going on here, and I’ll spare you the Pythagorean theorem refresher and say that these fields become much more important when you want to ensure the outcome of your deposit falls within specific acceptable tolerances.

LPs are meant to be active, dynamic, and automated (hence the A in AMMs). Participants must be cautious when depositing and withdrawing from LPs to try to protect against accidentally incurring significant losses. The point of LPs is to have assets moving into and out of the pool, and when you lose Asset A, you receive an equal valued amount in Asset B plus fees, resulting in a net gain of Asset B. If you’re not careful on price movements or are trading risky assets, you could end up with less of the asset you want and more of one that’s no longer as highly valued as when you initially deposited. Here be dragons, and dragons can be good when carefully trained but can also burn stuff down really quickly. Make sure you read our disclaimer at the footer of this page.

In this quest, you will be the only depositor in our own liquidity pool, and this is all on testnet, so things are pretty safe. The real world is not as safe or simple as this. Make sure you read our disclaimer at the footer of this page

**Max Amount A:** enter 100 of Asset A (XLM)

**Max Amount B:** enter 100 of Asset B (NOODLE) - note that you are issuing this new asset from your Quest Account by depositing it into the LP

**Min Price:** switch to Fraction and set to 1 / 1

**Max Price:** same as Min Price here: Fraction set to 1 / 1. This ensures the deposit will be exactly 100 of each asset

Normally, when depositing into an existing LP, you’ll want to carefully consider slippage in your max and min price values to ensure you enter at a price you’re happy with. For example, are you entering a 1:1 LP, a 1:2 LP, or something more variable like USD:EUR, but the tolerance can be somewhere between 1:1 and 0.99:1.01?

Now that you’ve deposited into the LP notice that your LP shares in the trustline have increased to 100. Why 100? It’s all about that triangle math, but the equation is actually pretty simple. It’s just the square root of the amount of Asset A you deposited multiplied by the amount of Asset B you deposited. So, in JS, this would be: `Math.sqrt(100 * 100) = 100`. If, however, we had deposited 50 XLM and 100 NOODLE, it would be: `Math.sqrt(50 * 100) = 70.7106781`. It’s hip to be square but rad to be a triangle.

5. Sign and submit the transaction to the testnet

Alright, let’s now use our newly created liquidity pool with the path payment operation.

6. Create a new account on the testnet or use one you’ve already created, we’ll call this Account 2
7. Establish a trustline between Account 2 and the NOODLE asset; since we are using the LP only for a path payment here, we won’t need a trustline for the pool shares
8. Build a new transaction with Path Payment Strict Receive as the operation and Account 2 as the source account and destination (yes, you can send yourself payments on Stellar)

Here we want to ensure we acquire exactly one NOODLE spending up to 1,000 XLM.

Let’s think about how much we’ll actually spend on the XLM side for this transaction. With both 100 XLM and NOODLE, the price is 1:1, so if we take one NOODLE, we should pay one XLM, right? Wrong. Well, yeah, right, but don’t forget the fee: 0.3%. So, `1 * 1.003 = 1.003`. Cool.

9. Sign and submit the transaction to the testnet

Now let’s see if the withdrawal amount is what we had expected.

10. Go to the Laboratory’s Explore Endpoints tab > Effects > Effects for Transaction and input the transaction’s hash; then look at account_debited. 1.0131405. Hmm. That’s not what we expected. BACK TO THE TRIANGLE MATH!

Turns out that LP deposit and withdrawal math is not linear, it’s curved. But not unpredictably so. The math is still pretty simple, you just have to keep in mind that for every withdrawal, you must refill the opposite side of the pool with a sufficient amount to keep the pool’s price balanced. The more you take from one side, the more you must fill on the opposite side exponentially.

Again I’ll spare you the math lesson and just give you the logic to pick through on your own time. This formula provides the logic for why our XLM amount withdrawn was what it was and may help in the future when you’re reasoning about how much it may cost to withdraw assets from an LP.

```
# All work is done in stroops (7 decimal places) to ensure correct rounding
# This also lets us work with integers throughout the process

stroop = 10,000,000

# This is the current state of our LP, we’ve deposited 100 each of XLM and NOODLE

currentXLM = 100 * stroop = 1,000,000,000
currentNOODLE = 100 * stroop = 1,000,000,000

# All LPs on Stellar, are “constant product” LPs. This is designed to keep the
# total value of each asset in relative equilibrium. The value `constantProduct`
# will change with each operation that affects the LP (deposit, withdraw, or
# asset swaps), but the calculations are always based on this formula. This will
# be important in the calculations to come.

constantProduct = currentXLM * currentNoodle = 10,000 XLM (1 * 10^18 in stroops)

# We are buying 1 NOODLE in this path payment which decreases the amount in the LP to 99 NOODLE

receiveNOODLE = 1 * stroop = 10,000,000

newNOODLE = currentNOODLE - receiveNOODLE = 99 * stroop = 990,000,000

# We calculate how much XLM we will need in the pool to maintain the overall value of the LP

newXLM = constantProduct / newNOODLE = 101.0101010 * stroop = 1,010,101,010

# Subtract the XLM that is already in the pool, and we know how much we need to send (before fees)

sendXLM = newXLM - currentXLM = 1.0101010 * stroop = 10,101,010

# The fee calculation can be a bit tricky. We need a 0.3% fee, but perhaps not
# in the way you are expecting. We are NOT finding 0.3% of the `sendXLM` amount.
# Rather, we are using a reverse percentage to find the total amount we must
# send so that our `sendXLM` amount is 99.7% of our total payment. The difference
# is subtle and not always intuitive, so be careful if you're calculating these.
# We also use `ceil()` to simplify rounding, and avoid payments with 0 fee.

feePercent = 0.003
paymentPercent = 1 - feePercent = 0.997
totalSendXLM = ceil(sendXLM / paymentPercent) = ceil(10,131,404,314) / stroop = 1.0131405 XLM
```

A more all-encompassing version of this calculation might look something like this:

$totalSendXLM = { ceil(({ constantProduct \over currentNOODLE - receiveNOODLE } - currentXLM) x ({ 1 \over 0.997 })) \over stroop}$

Alright, sick. We’ve successfully stood up and used an LP. Now it’s time to withdraw our shares and receive our portion of the LP’s deposits and fee rewards.

11. Construct a new transaction with Liquidity Pool Withdraw as the operation and the Quest Account as the source account

Here we’ll denote the amount of LP shares we’d like to burn and the minimum amount of Asset A and Asset B we need to receive for the transaction to execute. The Min Amounts are value guardrails to ensure you don’t accidentally submit an operation that withdraws amid a price downturn, liquidating your LP stake at an unacceptable rate.

For this quest, as the only holder in the pool, we’ll set the Min Amounts to 0 and the Amount to 100, indicating that we want entirely out of the LP and are willing to incur a zero-sum asset payback result.

12. Sign and submit the transaction to the testnet

13. Inspect your Quest Account’s balance to ensure the LP trustline is empty and the XLM balance has increased by that 1.0131405. Nice! Note that you don’t receive any of the NOODLE asset in the withdrawal operation since your Quest Account cannot hold the asset it has issued itself.

14. We’re done! Let’s call that a success. Go ahead and nab that badge!

**Further Reading**

- Stellar Documentation: [Liquidity on Stellar: SDEX and Liquidity Pools](https://developers.stellar.org/docs/encyclopedia/liquidity-on-stellar-sdex-liquidity-pools#liquidity-pools)
- SDF Blog: [Introducing Automated Market Makers on Stellar](https://stellar.org/blog/introducing-automated-market-makers-on-stellar) by Justin Rice
- Stellar Developers Blog: [Liquidity, Liquidity, Liquidity…](https://stellar.org/developers-blog/liquidity-liquidity-liquidity) by George Kudrayvtsev
- Stellar Core: [Core Advancement Proposal 0038 (CAP-0038)](https://github.com/stellar/stellar-protocol/blob/master/core/cap-0038.md)
  - [How Path Payments Utilize Liquidity Pools](https://github.com/stellar/stellar-protocol/blob/master/core/cap-0038.md#pathpaymentstrictsendop-and-pathpaymentstrictreceiveop)
- Horizon Servers: [Liquidity Pool API](https://developers.stellar.org/api/resources/liquiditypools/list/)

<sub><sup>**Legal Disclaimer:** Stellar Quest tasks are educational in nature and provide a structured sandbox environment for learning how to use Stellar software. All Stellar Quest tasks are performed on a testnet using test assets. Nothing in these Stellar Quest instructions should be construed as financial, legal or investment advice. Separately, if you choose to interact with AMM functionality and liquidity pools using real assets on Stellar mainnet then you should ensure you understand the technology, the assets and that you are aware of the risks involved in such operations. Remember, the value of crypto assets can be extremely volatile and unpredictable, which can result in significant losses in a short time including possibly a loss of total value.</sup></sub>
