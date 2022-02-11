package skiplist

import (
	"fmt"
	"memtable/pkg/coinflip"
	"memtable/pkg/storage"
)

var (
	MAX_LEVEL      = 16
	ErrKeyNotFound = fmt.Errorf("key not found")
)

func (s *Skiplist) AddLevel() {
	newRoot := MinNode()
	old := s.Root()
	newRoot.Child = old
	old.Parent = newRoot
	s.root = newRoot
}

func (s *Skiplist) Lift(n *Node) *Node {
	if n.Parent != nil {
		return n.Parent
	}

	tmp := n
	for ; tmp != nil && tmp.Parent == nil; tmp = tmp.Prev {
	}
	if tmp == nil {
		return nil
	}
	parent := NewNode(n.Key.Raw, nil)
	n.Parent = parent
	parent.Child = n
	tmp.Parent.Append(parent)
	return parent
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

// ScanLevel returns the maximum node such that n <= node <= key.
func (s *Skiplist) ScanLevel(n *Node, key Key) (*Node, error) {
	if !n.Key.Lte(key) {
		return nil, fmt.Errorf(
			"invariant: level scan must start with key less than target key; node=%v, key=%v",
			n, key,
		)
	}

	cur := n
	for ; cur.Key.Lte(key); cur = cur.Next {
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
		newNode := NewNode(key.Raw, value)
		n.Append(newNode)
		lvl := RandomHeight()
		if lvl == 0 {
			return nil
		}

		for i := 0; i < lvl; i++ {
			newNode = s.Lift(newNode)
		}
	}
	return nil
}

func RandomHeight() int {
	lvl := 0
	for coinflip.Flip() && lvl < MAX_LEVEL {
		lvl++
	}
	return lvl
}

func (s *Skiplist) Delete(key []byte) error {
	return s.delete(NewKey(key))
}
func (s *Skiplist) delete(key Key) error {
	n := s.Scan(key)
	if n.Key.Eq(key) {
		n.Prev.Next = n.Next
		if !n.IsLast() {
			n.Next.Prev = n.Prev
		}
		for n.Parent != nil {
			n = n.Parent
			n.Prev.Next = n.Next
			if !n.IsLast() {
				n.Next.Prev = n.Prev
			}
		}
	}
	return nil
}

func (s *Skiplist) RangeScan(start, limit []byte) (storage.Iterator, error) {
	startKey := NewKey(start)
	endKey := NewKey(limit)
	if endKey.Lt(startKey) {
		return nil, fmt.Errorf(
			"range invalid: start of range must be less than or equal to limit (start=%v, limit=%v)",
			start, limit,
		)
	}
	n := s.Scan(startKey)
	for n != nil && n.Key.Lt(startKey) {
		n = n.Next
	}
	return &Iterator{
		Current: nil,
		onDeck:  n,
		end:     endKey,
		err:     nil,
	}, nil
}

func New() *Skiplist {
	ret := &Skiplist{
		root: MinNode(),
	}
	for i := 0; i < MAX_LEVEL; i++ {
		ret.AddLevel()
	}
	return ret
}

type Iterator struct {
	Current *Node
	onDeck  *Node
	end     Key
	err     error
}

func (i *Iterator) Next() bool {
	if i.Done() {
		return false
	}
	i.Current = i.onDeck
	i.onDeck = i.Current.Next
	return true
}

func (i *Iterator) Done() bool {
	return i.Error() != nil || i.onDeck == nil || !i.onDeck.Key.Lt(i.end)
}

func (i *Iterator) Error() error {
	return i.err
}

func (i *Iterator) Key() []byte {
	return i.Current.Key.key
}

func (i *Iterator) Value() []byte {
	return i.Current.Val
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
