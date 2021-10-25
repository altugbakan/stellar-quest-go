# Create and submit a fee bump transaction
Fee channels are a common best practice in Stellar development. Their goal is to delegate fee payments away from user accounts for an improved UX. Protocol 13 saw a huge improvement to this flow with the introduction of fee bump transactions.

In this challenge your task is to create and execute a fee bump transaction which consumes the sequence number from your account but the transaction fee from some other account.

This is _actually_ how Stellar Quest delivers your prizes to you. A multi-operational transaction wrapped in a fee bump transaction. You pay the sequence number but the transaction fee is paid by Stellar Quest. How nice!
