# BFT

I have followed from start to end the algorithm given by the original authors [Lamport et. al.](https://lamport.azurewebsites.net/pubs/byz.pdf)

In the original algorithm, there is no synchronisation in the rounds, you simply keep messaging and unless the value of a node changes, it is always RETREAT. To quote Lamport et. al.:
> A traitorous commander may decide not to send any order. Since the lieutenants must obey some order, they need some default order to obey in this case. We let RETREAT be this default order.

## Effectiveness

For real world distributed systems this algorithm is **insanely** inefficient. It needs $\mathcal{O}(N!)$ messages to conclude. And this horrible for large distributed system that are used nowadays. However there exist other BFT algorithms (class of which is literally named "practical" BFT).

## Correctness

We can clearly see that all the honest lieutinants converge at the same consensus after receiving the designated number of incoming messages. A formal proof is in the paper as well. I have logged them `<ID.out>` files.
