package compression

import (
	"aidoc/internal/entity"
	"container/heap"
)

func buildHuffmanTree(freq map[int32]int) *entity.HuffNode {
	if len(freq) == 0 {
		return nil
	}
	h := &huffHeap{}
	heap.Init(h)
	for val, f := range freq {
		heap.Push(h, &entity.HuffNode{Value: val, Freq: f})
	}
	for h.Len() > 1 {
		left := heap.Pop(h).(*entity.HuffNode)
		right := heap.Pop(h).(*entity.HuffNode)
		parent:=&entity.HuffNode{
			Freq: left.Freq + right.Freq,
			Left: left,
			Right: right,
		}
		heap.Push(h,parent)
	}
	return heap.Pop(h).(*entity.HuffNode)
}
