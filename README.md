# LFU Cache Implementation

A Go implementation of a Least Frequently Used (LFU) cache with constant-time operations.

## Features

- O(1) time complexity for all core operations (Get, Put)
- Generic implementation supporting any comparable key and value types
- Configurable cache capacity
- Frequency tracking for cache entries
- Efficient traversal of all cache entries ordered by frequency
- Strict memory limits (O(capacity))

## Interface

```go
type Cache[K comparable, V any] interface {
    Get(key K) (V, error)
    Put(key K, value V)
    All() iter.Seq2[K, V]
    Size() int
    Capacity() int
    GetKeyFrequency(key K) (int, error)
}
```
## Implementation details
- Uses doubly-linked lists for O(1) operations
- Maintains frequency buckets for efficient eviction
- Thread-unsafe (concurrent access requires external synchronization)

## Usage example

```go
// Create cache with default capacity (5)
cache := lfu.New[string, int]()

// Or with custom capacity
cache := lfu.New[string, int](100)

// Basic operations
cache.Put("a", 1)
value, err := cache.Get("a")
freq, err := cache.GetKeyFrequency("a")

// Iterate through cache
for k, v := range cache.All() {
    fmt.Println(k, v)
}
```