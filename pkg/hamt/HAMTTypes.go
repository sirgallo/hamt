package hamt


type HAMTNode struct {
	Key string
	Value interface{}
	IsLeafNode bool
	BitMap uint32
	Children []*HAMTNode
}

type HAMT struct {
	BitChunkSize int
	TotalChildren int
	Root *HAMTNode
}