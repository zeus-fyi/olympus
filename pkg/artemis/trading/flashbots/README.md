Bundle ordering formula
The Flashbots builder uses a new algorithm designed to produce the most profitable block possible. This design introduces some important changes for searchers to be aware of:

```text
Instead of ranking and including bundles based off of effective gas price the algorithm now optimizes for overall block profit.
Top-of-block execution is no longer a guarantee.
Bundle ordering by effective gas price is no longer a guarantee.

Other transactions (e.g. from the mempool) may land between bundles (not between transactions in bundles, but between two different bundles).
For example:    

If you have a bundle comprised of transactions [B1, B2] and someone else has a bundle comprised of transactions [C1, C2]
and there are transactions in the mempool [t1, t2, ...], then the block may be built such that:

BLOCK_TXS = [..., B1, B2, t1, t2, C1, C2, ...].
```
