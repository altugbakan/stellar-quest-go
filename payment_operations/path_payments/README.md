# Path Payments
### Send or receive an asset that is different from the received or sent asset
<br>
With the concept of trustlines under our belt, let’s turn our attention to a new type of payment, the powerful path payment.

In a classic payment scenario, the asset sent is the same as the asset received.

In a path payment, the asset received is different than the asset sent. For example, you can send XLM and have the recipient receive USDC. How does this work, you ask? Well, rather than the operation transferring assets directly from one account to another, path payments cross through the DEX and/or liquidity pools before arriving at their final destination. For the path payment to succeed, there has to be a DEX offer or liquidity pool exchange path in existence. It can sometimes take several hops of conversion to succeed.

For example:

Account A sells XLM → [buy XLM / sell ETH → buy ETH / sell BTC → buy BTC / sell USDC] → Account B receives USDC

Path payments can fail if there are no viable exchange paths. However, this is still a convenient operation for currency conversions.

**In this quest, your challenge is to successfully send a path payment from the Quest Account to another account on the Stellar test network using the `pathPaymentStrictSend` or `pathPaymentStrictReceive` operation.**

Note: you may need to set up some offers using the Manage Buy Offer or Manage Sell Offer operations to execute a path payment successfully

Navigate to the Laboratory, and let’s get started.
1. Fill out the necessary fields, and select Path Payment Strict Send or Path Payment Strict Receive as the Operation Type. Let’s walk through what we see.