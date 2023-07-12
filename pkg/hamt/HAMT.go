package hamt

import "math"

import "github.com/sirgallo/hamt/pkg/utils"


func NewHAMT[T comparable]() *HAMT[T] {
	bitChunkSize := 5
	totalChildren := int(math.Pow(float64(2), float64(bitChunkSize)))
	
	return &HAMT[T]{
		BitChunkSize: bitChunkSize,
		TotalChildren: totalChildren,
		Root: &HAMTNode[T]{
			BitMap: 0,
			Children: []*HAMTNode[T]{},
		},
	}
}

func NewLeafNode[T comparable](key string, value T) *HAMTNode[T] {
	return &HAMTNode[T]{ 
		Key: key, 
		Value: value, 
		IsLeafNode: true,
	}
}

func NewInternalNode[T comparable]() *HAMTNode[T] {
	return &HAMTNode[T]{
		IsLeafNode: false, 
		BitMap: 0,
		Children: []*HAMTNode[T]{},
	}
}

func (hamt *HAMT[T]) Insert(key string, value T) {
	hamt.insertRecursive(hamt.Root, key, value, 0)
}

func (hamt *HAMT[T]) insertRecursive(node *HAMTNode[T], key string, value T, level int) {
	hash := utils.FnvHash(key)
	index := hamt.getSparseIndex(hash, level)

	if ! IsBitSet(node.BitMap, index) {
		newLeaf := NewLeafNode(key, value)
		node.BitMap = SetBit(node.BitMap, index)
		pos := hamt.getPosition(node.BitMap, hash, level)
		node.Children = ExtendTable[T](node.Children, node.BitMap, pos, newLeaf)
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				node.Value = value
			} else {
				newInternalNode := NewInternalNode[T]()
				node.Children[pos] = newInternalNode
				
				hamt.insertRecursive(node.Children[pos], childNode.Key, childNode.Value, level + 1)
				hamt.insertRecursive(node.Children[pos], key, value, level + 1)
			}
		} else { hamt.insertRecursive(node.Children[pos], key, value, level + 1) }
	}
}

func (hamt *HAMT[T]) Retrieve(key string) T {
	hash := utils.FnvHash(key)
	return hamt.retrieveRecursive(hamt.Root, key, hash, 0)
}

func (hamt *HAMT[T]) retrieveRecursive(node *HAMTNode[T], key string, hash uint32, level int) T {
	index := hamt.getSparseIndex(hash, level)
	
	if ! IsBitSet(node.BitMap, index) { 
		return utils.GetZero[T]() 
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			return childNode.Value
		} else { return hamt.retrieveRecursive(node.Children[pos], key, hash, level + 1) }
	}
}

func (hamt *HAMT[T]) Delete(key string) bool {
	hash := utils.FnvHash(key)
	return hamt.deleteRecursive(hamt.Root, key, hash, 0)
}

func (hamt *HAMT[T]) deleteRecursive(node *HAMTNode[T], key string, hash uint32, level int) bool {
	index := hamt.getSparseIndex(hash, level)
	
	if ! IsBitSet(node.BitMap, index) { 
		return false 
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]
		
		if childNode.IsLeafNode {
			if  key == childNode.Key {
				node.BitMap = SetBit(node.BitMap, index)
				node.Children = ShrinkTable[T](node.Children, node.BitMap, pos)
				
				return true
			}
			
			return false
		} else { 
			hamt.deleteRecursive(node.Children[pos], key, hash, level + 1) 
			popCount := calculateHammingWeight(node.Children[pos].BitMap)

			if popCount == 0 { // if empty internal node, remove from the mapped array
				node.BitMap = SetBit(node.BitMap, index)
				node.Children = ShrinkTable[T](node.Children, node.BitMap, pos)
			} 

			return true
		}
	}
}