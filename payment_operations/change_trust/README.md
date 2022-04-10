# Change Trust
### Create a trustline between two accounts for a designated asset
<br>
So far in our Stellar Quest Learn journey, we’ve just worked with Stellar’s native token, XLM. In this quest, we’ll learn about what needs to happen to start handling custom assets. Custom assets aren’t created, they exist in the form of trust on other accounts. Like Santa, they only exist for those who believe (am I banned?). In Stellar, this belief is called a trustline.

Trustlines are an explicit opt-in for an account to hold a particular asset. To hold a specific asset, an account must establish a trustline with the issuing account using the `changeTrust` operation.

**In this quest, your challenge is to establish a trustline between the Quest Account and another account for a specific asset using the changeTrust operation.**

Get on over to the [Laboratory](https://laboratory.stellar.org/) and get started.

1. Navigate to the Build Transaction tab, fill out the necessary information, and select Change Trust for the operation type