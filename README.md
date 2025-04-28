# go-bloom-filter
[![Go Version](https://img.shields.io/badge/Go-1.22%2B-blue)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)  
A simple Bloom filter library implementation using Go!

## Installation

```bash
go get github.com/wongzc/go-bloom-filter/bloomfilter
```

## Theory
- A **Bloom Filter** is, simply put, an array combined with multiple hash functions to probabilistically check if an item is present.
- **Hash functions**: generate multiple hash values for a given input.
- **Array**: stores bits at the positions indicated by hash values.

**Checking membership**:
- Compute the hash values.
- If **all** corresponding bits are set: the item **might be present** (false positives possible).
- If **any** bit is not set: the item is **definitely not present**.

**Typical use cases**:
- Cases where most checks are expected to return `false`.
- Detecting malicious URLs.
- Checking if a user has read an article before.

**Advantages**:
- Extremely **fast** and **memory-efficient** (only a single array needed to store millions of entries!)



## Features of This Implementation

- **Custom hash function injection** — users provide their own hash functions.
- **Memory-efficient** — stores bits in a `[]byte` slice (saving 1/8th the memory vs `[]bool`).
- **Double hashing** — derives `k` independent hashes using 2 base hash functions and index offsetting.
- **Dynamic sizing** — automatically calculates array size and number of hash functions based on:
  - Desired false positive rate
  - Expected number of items
- **Thread-safe** — safe for concurrent `Set` and `Get` via mutex locking.
- **Background batching** — new items are queued into a channel and processed in batches to reduce lock contention and improve performance.
- **Performance monitoring utilities** — with functions to calculate:
  - Theoretical false positive rate
  - Bit saturation rate
  - Bit distribution variance
  - Randomized bit heatmap visualization
- **Drop counting** — tracks how many `Set` operations were dropped due to full channel buffer.


## Formulas Used

### Bit Array Size

$$
m = -\frac{n \cdot \ln p}{(\ln 2)^2}
$$

- `n`: number of items to store
- `p`: desired false positive rate

---

### Number of Hash Functions

$$
k = \frac{m}{n} \cdot \ln 2
$$

- `m`: bit array size
- `n`: number of items

---

### False Positive Rate (FPR)

$$
\text{FPR} \approx \left(1 - e^{-\frac{kn}{m}}\right)^k
$$

- `m`: bit array size
- `k`: number of hash functions
- `n`: number of inserted items

---

## Example Usage

```go
package main

import (
    "fmt"
    "github.com/wongzc/go-bloom-filter/bloomfilter"
    "github.com/cespare/xxhash/v2"
    "time"
)

func hash1(data string) uint32 {
    return uint32(xxhash.Sum64String(data))
}

func hash2(data string) uint32 {
    h := xxhash.New()
    h.Write([]byte(data + "salt"))
    return uint32(h.Sum64())
}

func main() {
    bf := bloomfilter.New(10000, 0.01, hash1, hash2)

    bf.Set("hello")

    time.Sleep(150*time.Millisecond) // give background setter some time

    fmt.Println(bf.Get("hello")) // true
    fmt.Println(bf.Get("world")) // probably false

    bf.Close()
}
```

## Benchmarks

Performance results on:

- CPU: Intel Core i5-14400F
- OS: Windows
- Go version: 1.22+
- Dataset: 10 million random strings

| Operation    | Ops/sec    | Time per Op | Memory per Op | Allocations per Op |
|:-------------|:-----------|:------------|:--------------|:-------------------|
| Insert       | ~20 million | 53.9 ns     | 32 B          | 2 allocs           |
| Lookup       | ~10 million | 133.3 ns    | 56 B          | 3 allocs           |

**Additional Info**:
- Achieved bit saturation rate: **50.12%**.
- Bit distribution variance: **1.998**.
- Measured false positive rate: **~0.1044%** (close to theoretical 0.10%).

**Benchmark command**:

```bash
go test -bench=. -benchmem ./bloomfilter
```

**Example benchmark output**
```bash
BenchmarkInsert-16                      20212159                53.93 ns/op           32 B/op          2 allocs/op
BenchmarkLookup-16                      10276664               133.3 ns/op            56 B/op          3 allocs/op
BenchmarkMemoryUsage-16                 19393624                55.59 ns/op           32 B/op          2 allocs/op
```

**Summary**
- Insert operations complete in ~50ns.
- Lookup operations complete in ~130ns.
- Minimal memory overhead and allocations.
- False positive rate matches theoretical prediction very closely.

## Others

- Counting Bloom Filter
    - Allow Deletion from Bloom Filter
    - **More Memory Usage**: Use Counter to store, so memory is 4~16x higher than bloom filter
    - **Counter Overflow**: Risk of counter overflow if too many insertion at the same index
    - **Slow Performance**: Instead of flipping bits, need to read -> increase/decrease counter -> write to counter
    - **False Negative**: If mistakenly delete something that never insert, to avoid this, need to: 
        1. Check Couting Bloom Filter, proceed if it show "might exist"
        2. Proceed if cache show exist
        3. Delete and decrement Bloom Counter

- Cuckoo Filter
    - Store short hash (fingerprint) instead of bit
    - Allow deletion of item
    - Better lookup performance
    - Use lesser space in most case