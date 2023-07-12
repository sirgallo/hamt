# Hash Array Mapped Trie


## Installation

Import this repository to utilize the hash array mapped trie, implemented in `Go`.

in your `Go` project main directory (where the `go.mod` file is located)
```bash
go get github.com/sirgallo/hamt
go mod tidy
```


## Use

```go
package main

import "fmt"
import "github.com/sirgallo/hamt/pkg/hamt"

func main() {
  // initialize the hash array mapped trie
  // HAMT is generic, utilizing T comparable, so initialize with the appropriate type
  hamt := hamt.NewHAMT[string]()

  // insert a new key/val pair
  hamt.Insert("hello", "world")

  // retrieve the val at the key/val pair
  val1 := hamt.Retrieve("hello")
  fmt.Println("actual:", val1, "expected: world")

  // print the children
  hamt.PrintChildren()

  // delete a key/val pair
  hamt.Delete("hello", "world")
}
```


## Algorithm

Check out [HashArrayMappedTrie](./docs/HashArrayMappedTrie.md)


## Tests

[HAMT_test](./pkg/hamt/HAMT_test.go)

[HAMTUtils_test](./pkg/hamt/HAMTUtils_test.go)

To run tests:
```bash
go test -v ./pkg/hamt/tests
```