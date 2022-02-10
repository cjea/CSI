package skiplist

type Node struct {
	Prev, Next, Child, Parent *Node
	Key                       Key
	Val                       []byte
}

func (n *Node) IsLast() bool {
	return n.Next == nil
}

func (n *Node) IsLeaf() bool {
	return n.Child == nil
}

func (n *Node) Append(n2 *Node) *Node {
	tmp := n.Next

	n.Next = n2
	n2.Prev = n
	n2.Next = tmp
	if tmp != nil {
		tmp.Prev = n2
	}

	return n2
}

func NewNode(key Key, value []byte) *Node {
	return &Node{
		Prev:  nil,
		Next:  nil,
		Child: nil,
		Key:   key,
		Val:   value,
	}
}

func NewRootNode() *Node {
	return NewNode(Key{key: GLOBAL_MIN_KEY}, nil)
}

func MinNode() *Node {
	return NewNode(Key{key: GLOBAL_MIN_KEY}, nil)
}
