# Change Trust

### Create a trustline between two accounts for a designated asset

<br>
So far in our Stellar Quest Learn journey, we’ve just worked with Stellar’s native token, XLM. In this quest, we’ll learn about what needs to happen to start handling custom assets. Custom assets aren’t created, they exist in the form of trust on other accounts. Like Santa, they only exist for those who believe (am I banned?). In Stellar, this belief is called a trustline.

Trustlines are an explicit opt-in for an account to hold a particular asset. To hold a specific asset, an account must establish a trustline with the issuing account using the `changeTrust` operation.

**In this quest, your challenge is to establish a trustline between the Quest Account and another account for a specific asset using the changeTrust operation.**

Get on over to the [Laboratory](https://laboratory.stellar.org/) and get started.

1. Navigate to the Build Transaction tab, fill out the necessary information, and select Change Trust for the operation type

**Asset**
Assets can exist in three forms: alphanumeric 4, alphanumeric 12, and liquidity pool shares. We’ll focus on the first two for now, liquidity pool shares will come into play later.

Alphanumeric 4: the asset code is less than or equal to four characters

Alphanumeric 12: the asset code is less than or equal to twelve characters

Asset codes are case sensitive, so `Pizza`, `PIZZA`, and `pizza` are all different assets.

2. Select Alphanumeric 4 or Alphanumeric 12
3. Input the Asset Code – it can be anything within the alphanumeric 4 or 12 character bounds! `SANTA`, `Astronauts`, or even `SANTAInSpace`
4. Input the Issuer Account ID, which can be any account’s public key

The Issuer Account ID is the public key of the account that is issuing the asset. The issuing account can’t hold a balance of its own asset. This field determines who you are trusting to send you this particular asset.

**Trust Limit**

How much do you _really_ trust Santa?

This field determines how much of the asset your account can hold. Leaving this field blank allows for the maximum amount.

To remove a trustline, get rid of any remaining amount of the asset by sending it back to the issuing account (burning it) or sending it to a different account. Then set the Trust Limit to 0.

**Source Account**

The source for this operation is not the Issuer Account ID. You’re not creating an asset, you’re trusting an asset. Therefore, the source of the operation is the account trusting the asset.

5. With all the necessary information filled out, go ahead and sign and submit the transaction.
6. If the transaction was successful, congratulations! Go ahead and hit the Verify button in the Stellar Quest Learn screen to claim your next badge!

Bonus quest! Combine the payment operation with the change trust operation to make a payment from the issuing account to the trusting account. This actually “mints” the asset, essentially bringing it into existence.
