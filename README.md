`embarcadero` is a tool for conducting trustless, decentralized SF<->SC trades
on the Sia blockchain.

Bidders create partially-signed transactions that commit to an SF input and SC
output (or vice versa). These transactions can then be shared directly with a
counterparty, or posted publicly by including them in the arbitrary data of a
Sia transaction.

The partially-signed bid transactions are not valid on their own. They only
become valid when a counterparty adds their own input and output to the
transaction. No escrow service is necessary.

For example, Alice creates a bid with one input worth 5 SF, and one output worth
10MS. The output is addressed to her own wallet.

Bob sees the bid and thinks it's a good deal. He takes Alice's transaction and
adds an input worth 10MS and an output worth 5 SF, addressed to his own wallet.
Once he signs and broadcasts the transaction, Alice's wallet will lose 5 SF and
gain 10MS, and Bob's wallet will gain 5 SF and lose 10MS.

This repo includes a server, `embd`, that watches the blockchain (via a local
`siad` node) for bids contained in arbitrary data. The server can then be queried
to list e.g. all unfilled bids, or all completed trades.

Traders interact with the server via a client program, `embc`. This program can
also be used to conduct OTC trades, in which bids are sent directly to the
counterparty instead of being published on-chain. Note that such bids can also
be shared in public venues, e.g. a Discord server, allowing anyone in the venue
to fill the bid.

Basic usage:

```
# Place a bid to sell SF:
$ embc place 5SF 10MS
Bid created successfully.
Your bid has been submitted for inclusion in the next block.
When the bid appears on-chain, it will be listed in the 'bids' command.
Your bid ID is:
    19e0d14d67936f944d446eeb5cff3acfbd62da3282c256218cc8575415daf851

# Place a bid to buy SF:
embc place 15MS 3SF
Bid created successfully.
Your bid has been submitted for inclusion in the next block.
When the bid appears on-chain, it will be listed in the 'bids' command.
Your bid ID is:
    4ea343afcf400eacd2359c46116f7f8b767818d75ecfc2ee02cb794cac2fb1af

# List open bids:
$ embc bids
ID        Height  Bid    Ask
19e0d14d  257426  5 SF   10 MS
4ea343af  257182  15 MS  3 SF

# Fill a bid
$ embc fill 19e0d14d
Bid details:
Counterparty wants to trade their 5 SF for your 10 MS.
Accept? [y/n]: y
Bid filled successfully.
Tranasction ID: db4d1d047094b3351304e3ba18a2bea73e32ea301b5080d922e54e6461608cb1

# Create OTC bid using Skynet:
$ embc --skynet place 5SF 10MS
Bid created successfully!
Share this link with your desired counterparty:
    sia://AABTXD9PxgjkyJVfiJLcD-K4OreIEGNR1-48SMxMme2B5g

# Fill OTC bid using Skynet:
$ embc --skynet fill sia://AABTXD9PxgjkyJVfiJLcD-K4OreIEGNR1-48SMxMme2B5g
Bid details:
Counterparty wants to trade their 5 SF for your 10 MS.
Accept? [y/n]: y
Bid filled successfully.
Tranasction ID: 1854aff6687390bec126e13b24cf61123bda58f442ca7bddedcd753c53f0af7f
```
