package entity

type HuffNode struct {
	Value int32
	Freq  int
	Left  *HuffNode
	Right *HuffNode
}
type TrieNode struct {
	Left  *TrieNode
	Right *TrieNode
	Value *int32
}
