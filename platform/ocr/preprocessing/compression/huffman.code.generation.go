package compression

import "aidoc/internal/entity"

func generateCodes(node *entity.HuffNode, prefix string, codes map[int32]string) {
	if node.Left == nil && node.Right == nil {
		codes[node.Value] = prefix
		return
	}
	if node.Left != nil {
		generateCodes(node.Left, prefix+"0", codes)
	}
	if node.Right != nil {
		generateCodes(node.Right, prefix+"1", codes)
	}
}
