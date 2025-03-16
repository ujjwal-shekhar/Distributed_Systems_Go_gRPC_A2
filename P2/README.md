# MapReduce Design Analysis

## Why forking workers from client makes sense

Forking workers from the client creates a clean separation with minimal overhead. Each mapper gets its own process space which prevents memory leaks bleeding across mappers, and it naturally parallelizes without needing complex thread management. The short sleeps after fork are are not very pretty but they get the job done -- they wait for the processes to fully setup and get ready to listen to client requests. An alternative would have been to use pipes but that would steer away from Distributed settings (unless the FS allows distributed syscalls).

## Streaming chunks via RPC is efficient

Instead of reading entire files upfront (memory intensive), streaming chunks from master to mappers means smaller memory footprint. The RPC approach gives built-in flow control - if a mapper gets overwhelmed, backpressure naturally develops in the RPC stream. This is better than naive file distribution where you might overload nodes. Plus hey, now I have submitted something related to streaming and thats a plus xD

Using hash(key) % num_reducers to determine destination files is simple and effective. I immediately emit the keyvals onto local storage, simple and a bit different from the paper. Ideally I would buffer these to local and only emit after either a timeout or buffer limit exceeds and do an RPC to tell master, but this also works by using the return of RPC as a signal to the master.

## Synchronization via wait groups

The wait group approach for coordinating mappers/reducers completion is lightweight compared to alternatives. The RPC return itself doubles as synchronization signal - elegant design that avoids extra message passing, explained above too.

Assuming intermediate results fit in memory simplifies implementation significantly. If this assumption breaks, an external sort would add considerable complexity. This is a reasonable tradeoff if your expected workloads have manageable intermediate result sizes. For production systems, might need fallback strategy when memory pressure exceeds threshold. I have not done an external sort here.

## Local file output

Having reducers write outputs locally aligns with typical distributed filesystem expectations. The filesystem handles the complexities of distributed storage/replication, allowing MapReduce to focus on computation rather than storage logistics.

This design prioritizes simplicity and robustness over theoretical efficiency optimizations. For many real-world workloads, this is the correct tradeoff since complexity is often a bigger enemy than suboptimal performance curves. Here, the reducers can directly read the files in the intermediate outputs of the mappers, this is alright since  a file system would abstract syscalls on itself.

