# Set Options - Home Domain

### Set the home domain of an account

<br>

It’s time to dip our toes into the wild and wonderful world of the `setOptions` operation. There will be several quests that involve this operation but here we’ll just focus on the Home Domain field.

The Home Domain field points to the web domain where you post stellar.toml file, and your stellar.toml file provides a common place where the Internet can find information about your Stellar integration. There’s a lot of information you can put into the stellar.toml file, most common being `[[CURRENCIES]]` info for custom assets issued by your service.

The file is a bridge between domain names and the official Stellar accounts utilized by that service.

Anyone can look up your stellar.toml file, and it proves that you’re the owner of the HTTPS domain linked to a Stellar account and that you claim responsibility for the accounts and assets listed on it.

**In this quest, you will create, host, and link to a stellar.toml file from your Quest Keypair using the Set Options operation.**

If this seems intimidating, don’t fret! We’re going to provide a template demo that contains everything you need to configure a basic stellar.toml file. This will prove ownership over the domain you spin up.

Let’s get this going.

1. Navigate to the template demo: https://glitch.com/edit/#!/remix/sql-02-02
2. Click the stellar.toml tab in the left-hand navigation bar and input your Quest Keypair’s public key inside the quotes containing the `REPLACE_WITH_YOUR_QUEST_ACCOUNT_PUBLIC_KEY` text
3. Domain names must be less than 32 characters; to shorten the domain name, click the Settings tab in the left-hand navigation bar
4. Click Edit Project Details
5. Change the Project Name to what you’d like your domain to be and click Save
6. Click the Preview button at the bottom of the page and select Preview in a New Window
7. Copy the URL
8. Navigate to the Laboratory and build a transaction with Set Options as the operation type
9. Paste your domain into the Home Domain field
10. Remove the “http://” and the trailing “/” from the URL
11. Sign and submit the transaction to the testnet
12. If your transaction is successful, click Verify in the Stellar Quest Learn screen and claim your badge!

## Additional resources

Explore the spec for this stellar.toml Stellar Info File here: https://github.com/stellar/stellar-protocol/blob/master/ecosystem/sep-0001.md

Explore an example stellar toml file in the wild here: https://testanchor.stellar.org/.well-known/stellar.toml
