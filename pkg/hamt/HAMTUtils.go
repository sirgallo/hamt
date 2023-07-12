package hamt

import "fmt"
import "math/bits"
import "math"


func (hamt *HAMT[T]) getSparseIndex(hash uint32, level int) int {
	return GetIndex(hash, hamt.BitChunkSize, level)
}

func (hamt *HAMT[T]) getPosition(bitMap uint32, hash uint32, level int) int {
	sparseIdx := GetIndex(hash, hamt.BitChunkSize, level)
	mask := uint32((1 << sparseIdx) - 1)
	isolatedBits := bitMap & mask
	
	return calculateHammingWeight(isolatedBits)
}

func calculateHammingWeight(bitmap uint32) int {
	return bits.OnesCount32(bitmap)
}

func GetIndex(hash uint32, chunkSize int, level int) int {
	slots := int(math.Pow(float64(2), float64(chunkSize)))
	mask := uint32(slots - 1)
	shiftSize := slots - (chunkSize * (level + 1))

	return int(hash >> shiftSize & mask)
}

func SetBit(bitmap uint32, position int) uint32 {
	return bitmap ^ (1 <<  position)
}

func IsBitSet(bitmap uint32, position int) bool {
	return (bitmap & (1 << position)) != 0
}

func ExtendTable[T comparable](orig []*HAMTNode[T], bitMap uint32, pos int, newNode *HAMTNode[T]) []*HAMTNode[T] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode[T], tableSize)

	copy(newTable[:pos], orig[:pos])
	newTable[pos] = newNode
	copy(newTable[pos + 1:], orig[pos:])
	
	return newTable
}

func ShrinkTable[T comparable](orig []*HAMTNode[T], bitMap uint32, pos int) []*HAMTNode[T] {
	tableSize := calculateHammingWeight(bitMap)
	newTable := make([]*HAMTNode[T], tableSize)
	
	copy(newTable[:pos], orig[:pos])
	copy(newTable[pos:], orig[pos + 1:])

	return newTable
}

func (h *HAMT[T]) PrintChildren() {
	h.printChildrenRecursive(h.Root, 0)
}

func (h *HAMT[T]) printChildrenRecursive(node *HAMTNode[T], level int) {
	if node == nil { return }
	for i, child := range node.Children {
		if child != nil {
			fmt.Printf("Level: %d, Index: %d, Key: %s, Value: %v\n", level, i, child.Key, child.Value)
			h.printChildrenRecursive(child, level+1)
		}
	}
}