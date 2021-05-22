# Mint a Stellar Quest style NFT

This challenge will be difficult so prepare yourself. Take breaks, deep breaths, walk away for a while if you need to. Don't stress it, the rewards will be worth the effort.

In reality minting an NFT is actually fairly straightforward, and while there are several variant methods we've used over the past few series we're going to explore the simplest.

What you'll need to do is take the [PNG image](https://api.stellar.quest/badge/GCEE5H3RI2MFP4UQ4NHFKLGTIHILWA775AM7KTLU5HUBSLOBJN7M4RSL?network=public&v=1) provided, snag its base64 encoded string and bake that into your account's manageData fields utilizing both the key and value fields.

So just methodically slice off bits of the base64 string and pack those into both the key and value slots of consecutive manageData operations until you run out of string.

An important gotcha caveat is to prefix the key with **2** characters of indexing data in order to ensure accurate reassembly of the base64 string later.
