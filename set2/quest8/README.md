# Create and host a stellar.toml file for your account
Not all digital info should be stored on a blockchain. Some information needs to be mutable and derives no benefit from maintaining a blockchain paper trail. For these requirements we must look outside Stellar.

Blockchain database software like IPFS, Torrent or Filecoin can store stuff in a decentralized manner but are overkill when simply storing basic, mutable metadata for a Stellar account. For that we'll use SEP 1.

SEPs, or Stellar Ecosystem Proposals are ecosystem initiatives aimed at providing consensus arount common Stellar use cases. For SEP 1 that's probiding a common format for Stellar account metadata.

In today's challenge your task is to create, host and link to a stellar.toml file with an ***SQ02_EASTER_EGG*** field containing the text:

```text
Log into series 2 of Stellar Quest then visit quest.stellar.org/series2. Finally drag and drop your Stellar Quest series 2 badge PNG images onto the screen. Enjoy!
```
Note you won't be able to solve this challenge using only the laboratory. You'll need to host a toml file and for that you'll need a basic server.

## How to Host stellar.toml?
Please visit [my hosting repo](https://github.com/altugbakan/stellar-quest-go-hosting) to get help on hosting the file.
