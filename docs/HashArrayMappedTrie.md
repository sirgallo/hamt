# Hash Array Mapped Trie


## Background

Purely out of curiousity, I stumbled across the idea of `Hash Array Mapped Tries` in my search to create my own implementations of thread safe data structures in `Go` utilizing non-blocking operations, specifically, atomic operations like [Compare-and-Swap](https://en.wikipedia.org/wiki/Compare-and-swap). Selecting a data structure for maps was a challenge until I stumbled up the [CTrie](https://en.wikipedia.org/wiki/Ctrie), which is a thread safe, non-blocking implementation of a HAMT. Before implementing my own concurrent trie, I wanted to understand the fundamentals of HAMTs and how to implement them myself. This implementation doesn't aim to be the most optimized, but it does aim to be clear and readable, since there does not seem to be a lot of documentation regarding this data structure beyond the original whitepaper, written by Phil Bagwell in 2000 [1](https://lampwww.epfl.ch/papers/idealhashtrees.pdf).


## Overview

A `Hash Array Mapped Trie` is a memory efficient data structure that can be used to implement maps (associative arrays) and sets. HAMTs, when implemented with path copying and garbage collection, become persistent as well, which means that any function that utilizes them becomes pure.


### Why use a HAMT?

HAMTs can be useful to implement maps and sets in situations where you want memory efficiency. They also are dynamically sized, unlike other map implementations where the map needs to be resized.


## Design


### Data Structure

At it's core, a `Hash Array Mapped Trie` is an extended version of a classical [Trie](https://en.wikipedia.org/wiki/Trie). Tries, or Prefix Trees, are data structures that are composed of nodes, where all children of a node in the tree share a common prefix with parent node.

a diagram of a trie:
```
                    root --> usually null
                 /           \
                d              c
              /   \             \
            e       o            a
           /       / \          /  \
          n       t    g       t     r
```

Searching a trie is done [depth-first](https://en.wikipedia.org/wiki/Depth-first_search), where search time complexity is O(m), m is the length of the string being search. A node is checked to see if it contains the key at the index in the string it is searching for, creating a path as it traverses down the branches.

Hash Array Mapped Tries take this concept a step further. 

a `Hash Array Mapped Trie` Node (pseudocode):
```
HAMTNode {
  key KeyType
  value ValueType
  isLeafNode bool
  indexBitMap [k]bits
  children [k]*ptr
}
```

The basic idea of the HAMT is that the key is hashed. At each level within the HAMT, we take the t most significant bits, where 2^t represents the table size, and we index into the table using the index created by modified portion of the hash. The index mapping is explained further below.

Time Complexity for operations Search, Insert, and Delete on an order 32 HAMT is O(log32n), which is close to hash table time complexity.


### Hashing

Keys in the trie are first hashed before an operation occurs on them, using the [Fowler–Noll–Vo Hash Function](https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function). This creates a uint32 value which is then used to index the key within the trie structure.


### Array Mapping

Using the hash from the hashed key, we can determine:

1. The index in the sparse index
2. The index in the actual dense array where the node is stored


#### Sparse Index

Each node contains a sparse index for the mapped nodes in a uint32 bit map. 

To calculate the index:
```go
func getIndex(hash uint32, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	mask := uint32(slots - 1)
	shiftSize := slots - (chunkSize * (level + 1))

	return int(hash >> shiftSize & mask)
}
```

Using bitwise operators, we can get the index at a particular level in the trie by shifting the hash over the chunk size t, and then apply a mask to the shifted hash to return an index mapped in the sparse index. Non-zero values in the sparse index represent indexes where nodes are populated.


#### Dense Index

To limit table size and create dynamically sized tables to limit memory usage (instead of fixed size child node arrays), we can take the calculated sparse index for the key and, for all non-zero bits to the right of it, caclulate the population count ([Hamming Weight](https://en.wikipedia.org/wiki/Hamming_weight))

In go, we can utilize the `math/bits` package to calculate the hamming weight efficiently:
```go
func calculateHammingWeight(bitmap uint32) int {
	return bits.OnesCount32(bitmap)
}
```

calculating hamming weight naively:
```
hammingWeight(uint32 bits): 
  weight = 0
  for bits != 0:
    if bits & 1 == 1:
      weight++
    bits >>= 1
  
  return weight
```

to calculate position:
```go
func (hamt *HAMT) getPosition(bitMap uint32, hash uint32, level int) int {
	sparseIdx := getIndex(hash, hamt.BitChunkSize, level)
	mask := uint32((1 << (hamt.TotalChildren - sparseIdx)) - 1)
	isolatedBits := bitMap & mask
	
	return calculateHammingWeight(isolatedBits)
}
```

`isolatedBits` is all of the non-zero bits right of the index, which can be calculated by is applying a mask to the bitMap at that particular node. The mask is calculated from all of from the start of the sparse index right.


### Table Resizing


#### Extend Table

When a position in the new table is calculated for an inserted element, the original table needs to be resized, and a new row at that particular location will be added, maintaining the sorted nature from the sparse index. This is done using go array slices, and copying elements from the original to the new table.

```go
func ExtendTable(orig []*HAMTNode, bitMap uint32, pos int, newNode *HAMTNode) []*HAMTNode {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode, tableSize)
	
	copy(newTable[:pos], orig[:pos])
	newTable[pos] = newNode
	copy(newTable[pos + 1:], orig[pos:])
	
	return newTable
}
```


#### Shrink Table

Similarly to extending, shrinking a table will remove a row at a particular index and then copy elements from the original table over to the new table.

```go
func ShrinkTable(orig []*HAMTNode, bitMap uint32, pos int) []*HAMTNode {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode, tableSize)
	
	copy(newTable[:pos], orig[:pos])
	copy(newTable[pos:], orig[pos + 1:])

	return newTable
}
```


#### Insert

Pseudo-code:
```
  1.) hash the incoming key
  2.) calculate the index of the key
    I.) if the bit is not set 
      1.) create a new leaf node with the key and value
      2.) set the bit in the bitmap to 1
      3.) calculate the position in the dense array where the node will reside based on the current size
      4.) extend the table and add the leaf node
    II.) if the bit is set
      1.) caculate the position of the child node
        a.) if the node is a leaf node, and the key is the same as the incoming, update the value
        b.) if it is a leaf node, and the keys are not the same, create a new internal node (which acts as a branch), and then recursively add the new leaf node and the existing leaf node at this key to the internal node
        c.) if the node is an internal node, recursively move into that branch and repeat 2
```


### Retrieve

Pseudo-code:
```
  1.) hash the incoming key
  2.) calculate the index of the key
    I.) if the bit is not set, return null since the value does not exist
    II.) if the bit is set
      a.) if the node at the index in the dense array is a leaf node, and the keys match, return the value
      b.) otherwise, recurse down a level and repeat 2
```


### Delete

Pseudo-code:
```
  1.) hash the incoming key
  2.) calculate the index of the key
    I.) if the bit is not set, return false, since we are not deleting anything
    II.) if the bit is set
      a.) calculate the index in the children, and if it is a leaf node and the key is equal to incoming, shrink the table and remove the element, and set the bitmap value at the index to 0
      b.) otherwise, recurse down a level since we are at an internal node
```


## Sources

[1] [Ideal Hash Trees, Phil Bagwell](https://lampwww.epfl.ch/papers/idealhashtrees.pdf)

[2] [libhamt](https://github.com/mkirchner/hamt)

[3] [Hash Array Mapped Tries](https://worace.works/2016/05/24/hash-array-mapped-tries/)

[4] [Hash array mapped trie - wiki](https://en.wikipedia.org/wiki/Hash_array_mapped_trie)