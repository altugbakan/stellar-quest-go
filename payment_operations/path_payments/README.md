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

Navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test), and let’s get started.

1. Fill out the necessary fields, and select Path Payment Strict Send or Path Payment Strict Receive as the Operation Type. Let’s walk through what we see.

### Path Payment Strict Send
Since path payments cross the order book, it’s important to set safety parameters around the amount sent and the amount received. For denoting a specific send amount, use this strict send operation.

#### Destination
The destination account’s public key

#### Sending Asset
The asset you’re sending

#### Send Amount
The strict amount of the Sending Asset

#### Intermediate Path
This is optional and specifies the assets you want the payment to path through. You can find usable paths via the [/paths/strict-send](https://horizon.stellar.org/paths/strict-send) endpoint.

#### Destination Asset
The asset the destination account will receive

#### Minimum Destination Amount
Since path payments travel through live open offers, there is no static standard “rate” for sending and receiving. To counter this variable, you can guard the value conversion with this field to ensure the Destination will at least receive this amount of the Destination Asset.

#### Source Account
The source account for the Sending Asset.

### Path Payment Strict Receive
Statically sets the receive amount to ensure that the destination account receives a specified amount of the Destination Asset.

#### Destination
The destination account’s public key

#### Sending Asset
The asset you’re sending

#### Maximum Send Amount
As in the Minimum Destination Account above, this field specifies the maximum amount of the Sending Asset that the Source Account is willing to send.

#### Intermediate Path
This is optional and specifies the assets you want the payment to path through. You can find usable paths via the [/paths/strict-receive](https://horizon.stellar.org/paths/strict-receive) endpoint.

#### Destination Asset
The asset the destination account will receive

#### Destination Amount
The static amount of the Destination Asset that the Destination will receive

#### Source Account
The source account for the Sending Asset

2. Once you complete the fields for your Path Payment Strict Send or Path Payment Strict Receive operation, sign and submit the transaction
3. If the transaction is successful, great! Hit the Verify button on the Stellar Quest Learn page and claim your final badge in the Payment Operations series!