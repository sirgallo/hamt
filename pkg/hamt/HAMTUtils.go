package hamt

import "fmt"
import "math/bits"
import "math"


func (hamt *HAMT) getSparseIndex(hash uint32, level int) int {
	return getIndex(hash, hamt.BitChunkSize, level)
}

func (hamt *HAMT) getPosition(bitMap uint32, hash uint32, level int) int {
	sparseIdx := getIndex(hash, hamt.BitChunkSize, level)
	mask := uint32((1 << (hamt.TotalChildren - sparseIdx)) - 1)
	isolatedBits := bitMap & mask
	
	return calculateHammingWeight(isolatedBits)
}

func calculateHammingWeight(bitmap uint32) int {
	return bits.OnesCount32(bitmap)
}

func getIndex(hash uint32, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	mask := uint32(slots - 1)
	shiftSize := slots - (chunkSize * (level + 1))

	return int(hash >> shiftSize & mask)
}

func setBit(bitmap uint32, totalChildren int, position int) uint32 {
	return bitmap ^ (1 <<  (totalChildren - position))
}

func isBitSet(bitmap uint32, totalChildren int, position int) bool {
	return (bitmap & (1 << (totalChildren - position))) != 0
}

func ExtendTable(orig []*HAMTNode, bitMap uint32, pos int, newNode *HAMTNode) []*HAMTNode {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode, tableSize)
	
	copy(newTable[:pos], orig[:pos])
	newTable[pos] = newNode
	copy(newTable[pos + 1:], orig[pos:])
	
	return newTable
}

func ShrinkTable(orig []*HAMTNode, bitMap uint32, pos int) []*HAMTNode {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode, tableSize)
	
	copy(newTable[:pos], orig[:pos])
	copy(newTable[pos:], orig[pos + 1:])

	return newTable
}

func (h *HAMT) PrintChildren() {
	h.printChildrenRecursive(h.Root, 0)
}

func (h *HAMT) printChildrenRecursive(node *HAMTNode, level int) {
	if node == nil { return }
	for i, child := range node.Children {
		if child != nil {
			fmt.Printf("Level: %d, Index: %d, Key: %s, Value: %v\n", level, i, child.Key, child.Value)
			h.printChildrenRecursive(child, level+1)
		}
	}
}