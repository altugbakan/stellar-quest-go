# Manage Offers

### Create an offer to buy or sell a specific asset for another

<br>
Now that we know how to trust and issue assets other than the native XLM token, we have what we need to begin utilizing Stellar’s native decentralized exchange or DEX. Manage offer operations allow you to offer to buy or sell a specific amount of an asset at a specific exchange rate for a different asset. For example, sell 14 of asset A for 64 of asset B.

Stellar has three operations that manage these exchange offers:

- Manage buy offer
- Manage sell offer
- Create passive sell offer

**In this quest, your challenge is to open a buy or sell offer on the Quest Account using the `manageBuyOffer`, `manageSellOffer` or `createPassiveSellOffer` operation.**

Let’s walk through the three manage offer operations in the [Laboratory](https://laboratory.stellar.org/)

1. Navigate to the Build Transaction tab, fill out the necessary fields, and select Manage Buy Offer, Manage Sell Offer or Create Passive Sell Offer for the Operation Type. Let’s talk about what each of these fields means.

### Manage Buy Offer

Every offer is technically both a buy and sell offer. Selling 100 XLM for 10 USD is identical to buying 10 USD for 100 XLM. The difference is primarily syntactical to make it easier to reason about the creation of offers.

#### Selling

The asset you’re offering to give in the exchange

#### Buying

The asset you’re seeking to receive in the exchange

#### Amount you are buying

The amount of the Buying asset you’re required to receive for this offer to be taken

#### Price of 1 unit of buying in terms of selling

This one’s a little more complicated. Divide the amount you’re selling by the amount you’re buying. For example, if you want to buy 10 USD for 100 XLM, your price would be 100/10=10.

The reason this is tricky is that the result (10) isn’t necessarily the amount of either the selling or buying asset, it’s the price point for the counter asset of the offer. So, for the buy offer in the case above, multiply that back by the amount you’re buying to get the amount you’ll be selling: 10 XLM per USD \* 10 USD=100 XLM.

Understanding this calculation will make it easier to reason about buy and sell offers when the amount of the counter asset is larger than the primary asset.

#### Offer ID

Set to 0 to create a new offer.

To update or delete an existing offer, input the offer ID here. Find offer IDs by querying the various offer endpoints available via the Horizon API.

#### Source Account

The account that gives the selling asset and receives the buying asset in this offer.

### Manage Sell Offer

This operation is technically identical to Manage Buy Offer, with the primary and counter assets swapped

#### Selling

The asset you’re offering to give in the exchange

#### Buying

The asset you’re seeking to receive in the exchange

#### Amount you are selling

The amount of the Selling asset you’re required to give for this offer to be taken

#### Price of 1 unit of selling in terms of buying

Like in the above description, if you create an identical offer of receiving 10 USD in exchange for 100 XLM, flip the denominator and numerator (10/100). Rather than the price of 10, the price would be 0.1. And since the selling amount is 100, 0.1 USD per XLM \* 100 XLM=10 USD.

As you can see, the calculation is simple once you understand the numerator and denominator difference between the buy and sell offers.

#### Offer ID

Set to 0 to create a new offer.

To update or delete an existing offer, input the offer ID here. Find offer IDs by querying the various offer endpoints available via the Horizon API.

#### Source Account

The account that gives the selling asset and receives the buying asset in this offer

### Create Passive Sell Offer

This last operation creates an offer to sell one asset for another without taking a reverse offer of equal price. This allows you to maintain an order book that both buys and sells an asset equally without the offers actually executing. This is pretty common with stablecoins, where you want to maintain a 1:1 ratio for an asset pair where you act as the intermediary liquidity provider for both sides of the market.

There is no offer ID for this operation because it behaves exactly like a regular buy or sell offer once it’s been created. You can manage it via the regular buy or sell offer operations.

---

Note: offers may or may not be taken immediately upon submission to the network. If a relevant counter offer is already “on the books,” the offer will execute immediately. Otherwise, the offer will sit idle on the network until you either remove the offer, or a counterparty takes it. This will become important when building atomic transactions and should be considered when building more complex logic chains.

2. Once you fill in the fields for your Manage Buy Offer or Manage Sell Offer operation, sign and submit the transaction
3. If the transaction is successful, congratulations! Hit the Verify button on the Stellar Quest Learn site to claim your new badge!
