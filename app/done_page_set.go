package app

import "sync"

type donePageSet struct {
	set *sync.Map
}

func newDonePageSet() donePageSet {
	return donePageSet{&sync.Map{}}
}

func (d donePageSet) Add(s string) bool {
	_, ok := d.set.LoadOrStore(s, nil)
	return ok
}
