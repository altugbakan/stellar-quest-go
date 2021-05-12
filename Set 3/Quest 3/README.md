# Submit a hash signed transaction
So you've heard about Stellar's multisig, but did you know that you can sign with far more than just a simple Ed25519 secret key? That's right! There's also sha256 hashes and pre-authorized transaction hashes.

In this challenge your task is to add a very specific and special sha256 hash signer to your account and then to submit a second transaction using that signer to remove itself as a signer from the account. A sort of one time use key if you will.

Don't forget to check the [clue](https://horizon.stellar.org/transactions/3b00ce719e8c4f4d2218944fd60a78d2da83356f55c38ab733e0a46e386e25df/operations).
