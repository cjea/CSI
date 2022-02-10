package skiplist

import (
	"fmt"
	"testing"
)

func TestLte(t *testing.T) {
	table := []struct {
		b1, b2 []byte
		want   bool
	}{
		{b1: []byte{}, b2: []byte{}, want: true},
		{b1: []byte{}, b2: []byte{0}, want: true},
		{b1: []byte{0}, b2: []byte{}, want: false},
		{b1: []byte{0}, b2: []byte{0}, want: true},
		{b1: []byte{1}, b2: []byte{0}, want: false},
		{b1: []byte{1, 0}, b2: []byte{1}, want: false},
		{b1: []byte{1, 2, 3}, b2: []byte{1, 2, 3}, want: true},
	}
	for i, test := range table {
		name := fmt.Sprintf("test-%d", i)
		t.Run(name, func(t *testing.T) {
			if actual := lte(test.b1, test.b2); actual != test.want {
				t.Errorf("lte(%v, %v) = %v; got %v", test.b1, test.b2, test.want, actual)
			}
		})
	}
}

func TestKeyEq(t *testing.T) {
	table := []struct {
		k1, k2 Key
		want   bool
	}{
		{k1: Key{key: []byte{}}, k2: Key{key: []byte{}}, want: true},
		{k1: Key{key: []byte{0}}, k2: Key{key: []byte{}}, want: false},
		{k1: Key{key: []byte{0}}, k2: Key{key: []byte{0}}, want: true},
		{k1: Key{key: []byte{1, 2}}, k2: Key{key: []byte{1, 2}}, want: true},
		{k1: Key{key: []byte{1, 2, 3}}, k2: Key{key: []byte{1, 2}}, want: false},
	}
	for i, test := range table {
		name := fmt.Sprintf("test-%d", i)
		t.Run(name, func(t *testing.T) {
			if actual := test.k1.Eq(test.k2); actual != test.want {
				t.Errorf("%v.Eq(%v) = %v; got %v", test.k1, test.k2, test.want, actual)
			}
		})
	}
}

func TestNewKey(t *testing.T) {
	table := []struct {
		raw  []byte
		want []byte
	}{
		{raw: []byte{}, want: []byte{1}},
		{raw: []byte{0}, want: []byte{1, 0}},
		{raw: []byte{2, 3, 4}, want: []byte{1, 2, 3, 4}},
	}
	for i, test := range table {
		name := fmt.Sprintf("test-%d", i)
		t.Run(name, func(t *testing.T) {
			actual := NewKey(test.raw)
			for i := range test.want {
				if actual.key[i] != test.want[i] {
					t.Errorf("NewKey(%v) should have bytes %v; got %v", test.raw, test.want, actual.key)
				}
			}
		})
	}
}

func TestNodeAppend(t *testing.T) {
	t.Run("appending a node in the middle of a list", func(t *testing.T) {
		n2 := &Node{Next: nil}
		n1 := &Node{Next: n2}
		n2.Prev = n1
		app := NewNode(nil, nil)

		n1.Append(app)
		if n1.Next != app {
			t.Errorf("expected n1.Next to be app; got %v\n", n1.Next)
		}
		if app.Prev != n1 {
			t.Errorf("expected app.Prev to be n1; got %v\n", app.Prev)
		}
		if app.Next != n2 {
			t.Errorf("expected app.Next to be n2; got %v\n", app.Next)
		}
		if n2.Prev != app {
			t.Errorf("expected n2.Prev to be app; got %v\n", n2.Prev)
		}
	})
	t.Run("appending a node to the end of a list", func(t *testing.T) {
		n2 := &Node{Next: nil}
		n1 := &Node{Next: n2}
		n2.Prev = n1
		app := NewNode(nil, nil)

		n2.Append(app)
		if n1.Next != n2 {
			t.Errorf("expected n1.Next to be n2; got %v\n", n1.Next)
		}
		if app.Prev != n2 {
			t.Errorf("expected app.Prev to be n2; got %v\n", app.Prev)
		}
		if app.Next != nil {
			t.Errorf("expected app.Next to be nil; got %v\n", app.Next)
		}
		if n2.Prev != n1 {
			t.Errorf("expected n2.Prev to be n1; got %v\n", n2.Prev)
		}
	})
}

func TestScanLevel(t *testing.T) {
	t.Run("scanning empty level", func(t *testing.T) {
		l := New()
		n, err := l.ScanLevel(l.Root(), NewKey([]byte{1, 2, 3}))
		if err != nil {
			t.Fatal(err)
		}
		if n.Key.key[0] != 0 {
			t.Errorf("expected to find first node of (empty) level; got %+v", n)
		}
	})

	t.Run("scanning level with strictly lower values than target", func(t *testing.T) {
		term := NewNode([]byte{2}, nil)
		l := New()
		l.Root().Append(NewNode([]byte{2}, nil)).Append(term)
		n, err := l.ScanLevel(l.Root(), NewKey([]byte{3}))
		if err != nil {
			t.Fatal(err)
		}
		if n != term {
			t.Errorf("expected to find terminal node %+v; got %+v", term, n)
		}
	})

	t.Run("scanning level with higher and lower values than target", func(t *testing.T) {
		mid := NewNode([]byte{1}, nil)
		term := NewNode([]byte{3}, nil)
		l := New()
		l.Root().Append(mid).Append(term)
		n, err := l.ScanLevel(l.Root(), NewKey([]byte{2}))

		if err != nil {
			t.Fatal(err)
		}
		if n != mid {
			t.Errorf("expected to find middle node %+v; got %+v", mid, n)
		}

	})
}

func TestLift(t *testing.T) {
	t.Run("lifting root", func(t *testing.T) {
		l := New()
		root1 := l.Root()
		l.Lift(root1)
		root2 := l.Root()
		l.Lift(root2)
		if l.Root().Child != root2 {
			t.Errorf("expected root child to be %+v; got %+v", root2, l.Root().Child)
		}
		if l.Root().Child.Child != root1 {
			t.Errorf("expected root grandchild to be %+v; got %+v", root1, l.Root().Child.Child)
		}
		if l.Root() != root2.Parent {
			t.Errorf("expected root2 to have parent %+v; got %+v", l.Root(), root2.Parent)
		}
		if l.Root() != root1.Parent.Parent {
			t.Errorf("expected root1 to have grandparent %+v; got %+v", l.Root(), root1.Parent.Parent)
		}
	})

	t.Run("lifting non-root element", func(t *testing.T) {
		l := New()
		a0 := NewNode([]byte{1}, nil)
		originalRoot := l.Root()
		l.Root().Append(a0)
		l.Lift(a0)
		newRoot := l.Root()
		if originalRoot.Parent != newRoot {
			t.Errorf("expected original root to have parent %+v; got %+v", newRoot, originalRoot.Parent)
		}
		if newRoot.Child != originalRoot {
			t.Errorf("expected new root to have child %+v; got %+v", originalRoot, newRoot.Child)
		}
	})
}
