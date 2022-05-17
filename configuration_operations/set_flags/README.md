# Set Flags
### Control access to an asset by setting flags on its issuing account
<br>
So far, our Stellar Quest Learn adventure has been a breeze, right? Well, this quest is a gale-force wind. It’s your first pro-level concept. Deep breaths, I believe in you. Your mom believes in you. Your cat believes in you. Just kidding, your cat is a jerk.

Let’s get to it.

Flags control access to an asset and are set by the issuing account of that asset. You can set flags at any time — but they are not applied retroactively to existing trustlines.

**In this quest, you will issue an asset to the Quest Keypair with the Authorization Required flag set using the Set Options and Set Trust Line Flags operations. This transaction requires several operations to succeed.**

1. Nothing new here — navigate to the [Laboratory](https://laboratory.stellar.org/#account-creator?network=test)
2. Either create a new account or use one you’ve already created- this will be your issuing account
3. Navigate to the Build Transaction tab
4. Configure a transaction that issues an authorization required asset from the issuing account to the Quest Account

Let’s talk about the flag operations involved in this transaction:

### **Set Options**
There are four possible flags you can set using the Set Options operation, and you can set them concurrently. These flags are set at the account level.

**Authorization Required**

An issuer must approve an account before that account can hold its asset

**Authorization Revocable**

An issuer can revoke an existing trustline’s authorization, preventing that account from transferring or trading the asset and closing open orders

**Authorization Immutable**

None of the other flags (Authorization Required, Authorization Revocable or Clawback Enabled) can be set and the issuing account can’t be merged

**Clawback Enabled**

An issuer can claw back any portion of an account's balance of the asset

### **Set Trust Line Flags**
This operation allows the issuing account to configure various flags for individual trustlines. Note that the issuing account must always be the source account for this operation.

Let’s go through the various components here:

**Asset**

The asset we’re setting flags for

**Trustor**

The account we’re configuring the trustline flag(s) for

**Set Flags**

There are two flags you can set in this operation:

- Authorized — signifies complete authorization, allowing an account to transact freely with the asset to make and receive payments and place orders
- Authorized to maintain liabilities — limits authorization, allowing an account to maintain current orders but preventing other operations

**Clear Flags**

There are three flags you can remove with this operation:

- Authorized
- Authorized to maintain liabilities
- Clawback enabled

5. Once you have built the transaction, sign it with the secret keys from both the issuing account and Quest Account
6. Submit the transaction to the testnet
7. If your transaction is a success, go ahead and hit the Verify button on the Stellar Quest Learn page to ensure you completed the quest correctly
8. If so, congratulations! You’ve completed a pretty complex transaction on the Stellar network. Claim that badge, you deserve it!
**Note:** It wasn't required to complete this quest, but you can use a single transaction to grant authorization, send the asset, and remove authorization, essentially creating an authorization sandwich. That way, you can allow a user to hold an asset, but require them to get permission from you, the issuer, before they transfer or sell it. Pretty neat!

## Example scenarios
If you’re having trouble, these example scenarios may help you understand how flags work a bit better.

### **Scenario A - Auth Required Asset**
As an asset issuer, you want to maintain control over your asset, which you do by setting the Authorization Required and Authorization Revocable flags on your account. Set the Authorization Required flag to ensure every potential holder obtains explicit permission from your account via the setTrustLineFlags operation before they can receive your asset. Set the Authorization Revocable flag to be able to clear flags, limit permissions, and clear orders for asset holders via the setTrustLineFlags operation.

This flow can be used for some pretty fancy DeFi situations. For example, you issue locked assets and then only open permissions in certain transaction scenarios before relocking the asset at the end of the transaction. This is useful in royalty payment requirements for NFT assets where part of the sale of an NFT includes additional payments to the original NFT author.

### **Scenario B - Clawback Enabled Asset**
As an asset issuer, you may wish to have the ability to clawback an asset from specific holders. Do this by setting the Authorization Revocable and Clawback Enabled flags. In this scenario, any holder of an asset issued by this account can have its balance revoked and sent back to the issuing account, burning the asset.

This is used in the regulated asset space where issuers need to maintain compliance with the assets they issue. It’s also a great trick for enabling some fun and interesting DeFi scenarios where assets can be entertainingly sent and removed under specific contract criteria.