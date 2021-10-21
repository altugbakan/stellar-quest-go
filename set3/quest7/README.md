# Acquire and make use of a SEP-0010 JWT
Outside of the relatively simple and controlled world of Stellar operations, there's a whole universe of use cases and implementations. Even here though there is need for order and interoperability.

These needs are met by Stellar Ecosystem Proposals or SEPs. SEPs are the Stellar wilderness guidebooks ensuring everyone is following the same path and rules and is thus able to interoperate with each other.

***SEP-0010*** is an aythentication SEP outlining how to prove ownership of a Stellar account to a service. It is used in many other SEPs so it's an important foundational SEP to understand.

Today's task will be to acquire a ***SEP-0010*** JWT and then embed that JWT back into your account's manageData fields in identical fashion to how we embedded the NFT data in the previous quest.

Please note that this embed step **is not** part of ***SEP-0010*** and is **definitely not** something you'd do in practice. It is just included here as a method for reading back and making use of the generated ***SEP-0010*** as part of the verification step. It's also a good refresher for Quest 6.