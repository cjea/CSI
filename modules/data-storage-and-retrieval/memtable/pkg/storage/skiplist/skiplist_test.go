package skiplist

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
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
	t.Run("lifting an element", func(t *testing.T) {
		l := New()
		a0 := NewNode([]byte{1}, nil)
		l.Root().Child.Append(a0)
		l.Lift(a0)
		if l.Root().Next != a0.Parent {
			t.Errorf("expected index to be lifted %p; got %+v", a0.Parent, l.Root())
		}
	})
}

func TestGet(t *testing.T) {
	l := New()
	for i := 0; i < 50; i++ {
		l.Put([]byte{byte(i)}, []byte(fmt.Sprintf("val_%d", i)))
	}
	table := []struct {
		key  []byte
		want string
		err  error
	}{
		{key: []byte{byte(0)}, want: "val_0", err: nil},
		{key: []byte{byte(20)}, want: "val_20", err: nil},
		{key: []byte{byte(42)}, want: "val_42", err: nil},
		{key: []byte{byte(51)}, want: "", err: ErrKeyNotFound},
	}
	for i, test := range table {
		name := fmt.Sprintf("test-%d", i)
		t.Run(name, func(t *testing.T) {
			val, err := l.Get(test.key)
			if err != test.err || string(val) != test.want {
				t.Errorf(
					"Get(%v) should equal %v (err: %v); got %v (err %v)",
					test.key, test.want, test.err, string(val), err,
				)
			}
		})
	}
}
func TestDelete(t *testing.T) {
	l := New()
	for i := 0; i < 50; i++ {
		l.Put([]byte{byte(i)}, []byte(fmt.Sprintf("val_%d", i)))
	}
	k := []byte{(byte(42))}
	l.Delete(k)
	if ok, _ := l.Has(k); ok {
		t.Errorf("Key %v should have been deleted but it was found", k)
	}
}

func TestRangeScan(t *testing.T) {
	l := New()
	for i := 0; i < 26; i++ {
		l.Put([]byte{byte(i + 97)}, []byte(fmt.Sprintf("val_%c", i+97)))
	}
	t.Run("half-inclusive range", func(t *testing.T) {
		it, err := l.RangeScan([]byte("ba"), []byte("f"))
		if err != nil {
			t.Fatal(err)
		}
		res := []string{}
		for it.Next() {
			res = append(res, string(it.Value()))
		}
		expected := []string{"val_c", "val_d", "val_e"}
		if len(res) != 3 {
			t.Errorf("Expected %d values; got %d (%#v)", len(expected), len(res), res)
		}
		if res[0] != expected[0] || res[1] != expected[1] || res[2] != expected[2] {
			t.Errorf("Expected %#v; got %#v", expected, res)
		}
	})

	t.Run("past the end of the list", func(t *testing.T) {
		it, err := l.RangeScan([]byte("z"), []byte("zz"))
		if err != nil {
			t.Fatal(err)
		}
		res := []string{}
		for it.Next() {
			res = append(res, string(it.Value()))
		}
		expected := []string{"val_z"}
		if len(res) != 1 {
			t.Errorf("Expected %d value; got %d (%#v)", len(expected), len(res), res)
		}
		if res[0] != expected[0] {
			t.Errorf("Expected %#v; got %#v", expected, res)
		}
	})

	t.Run("before the start of the list", func(t *testing.T) {
		it, err := l.RangeScan([]byte("01"), []byte("c"))
		if err != nil {
			t.Fatal(err)
		}
		res := []string{}
		for it.Next() {
			res = append(res, string(it.Value()))
		}
		expected := []string{"val_a", "val_b"}
		if len(res) != 2 {
			t.Errorf("Expected %d values; got %d (%#v)", len(expected), len(res), res)
		}
		if res[0] != expected[0] || res[1] != expected[1] {
			t.Errorf("Expected %#v; got %#v", expected, res)
		}
	})
}

func preGenRandomKeys(prefix string, amount, low, high int) []string {
	rand.Seed(time.Now().UnixNano())
	ret := make([]string, amount)
	for i := 0; i < amount; i++ {
		r := rand.Intn(high)
		ret[i] = prefix + fmt.Sprint(r)
	}
	return ret
}

func seed(n int) *Skiplist {
	l := New()
	for i := 0; i < n; i++ {
		str := fmt.Sprint(i)
		l.Put([]byte("key-"+str), []byte("val-"+str))
	}
	return l
}

var seededList4k = seed(1 << 12)
var seededList8k = seed(1 << 13)
var seededList16k = seed(1 << 14)
var seededList32k = seed(1 << 15)

func BenchmarkPut1k(b *testing.B) {
	l := seed(1 << 10)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Put([]byte("key-benchtest"), []byte{9, 10, 11, 12})
	}
}
func BenchmarkPut2k(b *testing.B) {
	l := seed(1 << 11)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Put([]byte("key-benchtest"), []byte{9, 10, 11, 12})
	}
}
func BenchmarkPut4k(b *testing.B) {
	l := seed(1 << 12)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Put([]byte("key-benchtest"), []byte{9, 10, 11, 12})
	}
}
func BenchmarkPut8k(b *testing.B) {
	l := seed(1 << 13)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Put([]byte("key-benchtest"), []byte{9, 10, 11, 12})
	}
}

func BenchmarkPut16k(b *testing.B) {
	l := seed(1 << 14)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		l.Put([]byte("key-benchtest"), []byte{9, 10, 11, 12})
	}
}

var keys = preGenRandomKeys("key-", 1<<10, 1, 1<<14)

func BenchmarkGet4k(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k := []byte(keys[i%len(keys)])
		b.StartTimer()
		_, _ = seededList4k.Get(k)
	}
}
func BenchmarkGet8k(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k := []byte(keys[i%len(keys)])
		b.StartTimer()
		_, _ = seededList8k.Get(k)
	}
}
func BenchmarkGet16k(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k := []byte(keys[i%len(keys)])
		b.StartTimer()
		_, _ = seededList16k.Get(k)
	}
}

func BenchmarkGet32k(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		k := []byte(keys[i%len(keys)])
		b.StartTimer()
		_, _ = seededList32k.Get(k)
	}
}
