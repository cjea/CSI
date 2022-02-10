package skiplist

import (
	"fmt"
	"memtable/pkg/exercise"
)

var (
	ErrKeyNotFound = fmt.Errorf("key not found")
)

func (s *Skiplist) EnsureRootLevelParent(n *Node) *Node {
	for n.Prev != nil {
		n = n.Prev
	}
	if n.Parent != nil {
		return n.Parent
	}
	newRoot := MinNode()
	newRoot.Child = s.root
	s.root.Parent = newRoot
	s.root = newRoot
	return newRoot
}

func (s *Skiplist) Lift(n *Node) *Node {
	s.EnsureRootLevelParent(n)
	if n.Parent != nil {
		return n.Parent
	}

	tmp := n
	for ; tmp.Parent == nil; tmp = tmp.Prev {
	}
	parent := NewNode(n.Key.Raw, nil)
	n.Parent = parent
	parent.Child = n
	tmp.Parent.Append(parent)
	return tmp.Parent
}

type Skiplist struct {
	root *Node
}

func (s *Skiplist) Print() {
	start := s.Root()
	for ; start.Child != nil; start = start.Child {
	}
	for ; start != nil; start = start.Next {
		fmt.Printf("%v", start.Key.key)
		tmp := start
		for tmp.Parent != nil {
			fmt.Printf(" <- %v", tmp.Key.key)
			tmp = tmp.Parent
		}
		fmt.Printf("\n")
	}
}

func (s *Skiplist) Root() *Node {
	if s.root == nil {
		s.root = NewRootNode()
	}
	return s.root
}

func (s *Skiplist) ScanLevel(n *Node, key Key) (*Node, error) {
	if !lte(n.Key.key, key.key) {
		return nil, fmt.Errorf("invariant: level scan must start with key less than target key")
	}

	cur := n
	for ; lte(cur.Key.key, key.key); cur = cur.Next {
		if cur.IsLast() {
			return cur, nil
		}
	}
	return cur.Prev, nil
}

func (s *Skiplist) Scan(key Key) *Node {
	cur, err := s.ScanLevel(s.Root(), key)
	must(err)
	for !cur.IsLeaf() {
		cur, err = s.ScanLevel(cur.Child, key)
		must(err)
	}

	return cur
}

func (s *Skiplist) Get(key []byte) (value []byte, err error) {
	return s.get(NewKey(key))
}
func (s *Skiplist) get(key Key) (value []byte, err error) {
	n := s.Scan(key)
	if !n.Key.Eq(key) {
		return nil, ErrKeyNotFound
	}
	return n.Val, nil
}

func (s *Skiplist) Has(key []byte) (ret bool, err error) {
	return s.has(NewKey(key))
}
func (s *Skiplist) has(key Key) (ret bool, err error) {
	n := s.Scan(key)
	return n.Key.Eq(key), nil
}

func (s *Skiplist) Put(key, value []byte) error {
	return s.put(NewKey(key), value)
}

func (s *Skiplist) put(key Key, value []byte) error {
	n := s.Scan(key)
	if n.Key.Eq(key) {
		n.Val = value
	} else {
		n.Append(NewNode(key.Raw, value))
	}
	// TODO: traverse up levels and insert skip-nodes.
	//       Maybe have #Scan return a list access-path nodes from each level.
	//
	//    While bubble
	//      child = D, L = L + 1
	//      D' = (key, child=D); insert D' in Level[L++]
	//      child = D'
	//
	panic("Skiplist#Put is not implemented")
}

func (s *Skiplist) Delete(key []byte) error {
	return s.delete(NewKey(key))
}
func (s *Skiplist) delete(key Key) error {
	n := s.Scan(key)
	if n.Key.Eq(key) {
		n.Prev.Next = n.Next
	}
	return nil
}

func (s *Skiplist) RangeScan(start, limit []byte) (exercise.Iterator, error) {
	panic("Skiplist#RangeScan is not implemented")
}

func New() *Skiplist {
	ret := &Skiplist{}
	ret.Root()
	return ret
}

type Iterator struct{}

func (i Iterator) Next() bool {
	panic("Iterator#Next not implemented")
}
func (i Iterator) Error() error {
	panic("Iterator#Error not implemented")
}
func (i Iterator) Key() []byte {
	panic("Iterator#Key not implemented")
}
func (i Iterator) Value() []byte {
	panic("Iterator#Value not implemented")
}

func lte(b1, b2 []byte) bool {
	l1 := len(b1)
	l2 := len(b2)
	for p := 0; p < l1 && p < l2; p++ {
		if b1[p] > b2[p] {
			return false
		}
		if b1[p] < b2[p] {
			return true
		}
	}
	return l1 <= l2
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
