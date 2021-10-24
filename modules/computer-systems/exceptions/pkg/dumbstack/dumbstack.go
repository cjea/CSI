package dumbstack

type DumbStack struct{ items []string }

func (ds *DumbStack) Push(str string) {
	ds.items = append(ds.items, str)
}
func (ds *DumbStack) Pop() (string, bool) {
	nItems := len(ds.items)
	if nItems < 1 {
		return "", false
	}
	ret := ds.items[nItems-1]
	ds.items = ds.items[0:nItems]

	return ret, true
}

func New() *DumbStack {
	return &DumbStack{}
}
