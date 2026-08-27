package compression

import "aidoc/internal/entity"

type huffHeap []*entity.HuffNode

func (h huffHeap) Len() int            { return len(h) }
func (h huffHeap) Less(i, j int) bool  { return h[i].Freq < h[j].Freq }
func (h huffHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *huffHeap) Push(x interface{}) { *h = append(*h, x.(*entity.HuffNode)) }
func (h *huffHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
