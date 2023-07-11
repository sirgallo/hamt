package hamt

import "math"

import "github.com/sirgallo/hamt/pkg/utils"


func NewHAMT() *HAMT {
	bitChunkSize := 5
	totalChildren := int(math.Pow(float64(2), float64(bitChunkSize)))
	
	return &HAMT{
		BitChunkSize: bitChunkSize,
		TotalChildren: totalChildren,
		Root: &HAMTNode{
			BitMap: 0,
			Children: []*HAMTNode{},
		},
	}
}

func NewLeafNode(key string, value interface{}) *HAMTNode {
	return &HAMTNode{ 
		Key: key, Value: value, 
		IsLeafNode: true, BitMap: 0, 
		Children: []*HAMTNode{},
	}
}

func NewInternalNode() *HAMTNode {
	return &HAMTNode{
		IsLeafNode: false, BitMap: 0,
		Children: []*HAMTNode{},
	}
}

func (hamt *HAMT) Insert(key string, value interface{}) {
	hamt.insertRecursive(hamt.Root, key, value, 0)
}

func (hamt *HAMT) insertRecursive(node *HAMTNode, key string, value interface{}, level int) {
	hash := utils.FnvHash(key)
	index := hamt.getSparseIndex(hash, level)

	if ! isBitSet(node.BitMap, hamt.TotalChildren, index) {
		newLeaf := NewLeafNode(key, value)
		node.BitMap = setBit(node.BitMap, hamt.TotalChildren, index)
		pos := hamt.getPosition(node.BitMap, hash, level)
		node.Children = ExtendTable(node.Children, node.BitMap, pos, newLeaf)
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]

		if childNode.IsLeafNode {
			if key == childNode.Key {
				node.Value = value
			} else {
				newInternalNode := NewInternalNode()
				node.Children[pos] = newInternalNode
				
				hamt.insertRecursive(node.Children[pos], childNode.Key, childNode.Value, level + 1)
				hamt.insertRecursive(node.Children[pos], key, value, level + 1)
			}
		} else { hamt.insertRecursive(node.Children[pos], key, value, level + 1) }
	}
}

func (hamt *HAMT) Retrieve(key string) interface{} {
	hash := utils.FnvHash(key)
	return hamt.retrieveRecursive(hamt.Root, key, hash, 0)
}

func (hamt *HAMT) retrieveRecursive(node *HAMTNode, key string, hash uint32, level int) interface{} {
	index := hamt.getSparseIndex(hash, level)
	
	if ! isBitSet(node.BitMap, hamt.TotalChildren, index) { 
		return nil 
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]

		if childNode.IsLeafNode && key == childNode.Key {
			return childNode.Value
		} else { return hamt.retrieveRecursive(node.Children[pos], key, hash, level + 1) }
	}
}

func (hamt *HAMT) Delete(key string) bool {
	hash := utils.FnvHash(key)
	return hamt.deleteRecursive(hamt.Root, key, hash, 0)
}

func (hamt *HAMT) deleteRecursive(node *HAMTNode, key string, hash uint32, level int) bool {
	index := hamt.getSparseIndex(hash, level)
	
	if ! isBitSet(node.BitMap, hamt.TotalChildren, index) { 
		return false 
	} else {
		pos := hamt.getPosition(node.BitMap, hash, level)
		childNode := node.Children[pos]
		
		if childNode.IsLeafNode {
			if  key == childNode.Key {
				node.BitMap = setBit(node.BitMap, hamt.TotalChildren, index)
				node.Children = ShrinkTable(node.Children, node.BitMap, pos)
				
				return true
			}
			
			return false
		} else { 
			hamt.deleteRecursive(node.Children[pos], key, hash, level + 1) 
			popCount := calculateHammingWeight(node.Children[pos].BitMap)

			if popCount > 0 {
				node.BitMap = setBit(node.BitMap, hamt.TotalChildren, index)
				node.Children = ShrinkTable(node.Children, node.BitMap, pos)
			} 

			return true
		}
	}
}