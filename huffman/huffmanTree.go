package huffman

// Node - Корень дерева
type Node struct {
	Val       byte
	Frequency int
	Left      *Node
	Right     *Node
}

// CreateCompressionMap - Создание мапы кодов символов, используя дерево Хаффмана
func (N *Node) CreateCompressionMap(codeBefore, direction string, encodingDict map[byte]string) {
	if N.Left == nil && N.Right == nil {
		encodingDict[N.Val] = codeBefore + direction
	} else {
		N.Left.CreateCompressionMap(codeBefore+direction, "0", encodingDict)
		N.Right.CreateCompressionMap(codeBefore+direction, "1", encodingDict)
	}
}

// Tree - Дерево
type Tree struct {
	Root     Node
	Priority int
	Index    int
}
