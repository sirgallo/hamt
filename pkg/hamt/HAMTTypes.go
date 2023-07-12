package hamt


type HAMTNode [T comparable] struct {
	Key string
	Value T
	IsLeafNode bool
	BitMap uint32
	Children []*HAMTNode[T]
}

type HAMT [T comparable] struct {
	BitChunkSize int
	TotalChildren int
	Root *HAMTNode[T]
}