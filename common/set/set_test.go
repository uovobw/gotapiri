package set

import (
	"crypto/sha256"
	//"fmt"
	"io"
	"testing"
)

type TestObj struct {
	Name string
	Id   int
}

func (t TestObj) Hash() (hash string) {
	h := sha256.New()
	io.WriteString(h, t.Name+string(t.Id))
	return string(h.Sum(nil))
}

func TestNewSet(t *testing.T) {
	s := New(0)
	if s.size != 0 {
		t.Fatalf("no size set")
	}
	ss := New(1)
	if ss.size != 1 {
		t.Fatalf("did not set size")
	}
}

func TestSetLimit(t *testing.T) {
	s := New(2)
	ok, err := s.Add(TestObj{"one", 1})
	ok, err = s.Add(TestObj{"two", 2})
	if !ok && err == nil {
		t.Fatalf("Added two elements but failed")
	}
	ok, err = s.Add(TestObj{"three", 3})
	if ok || (err == nil) {
		t.Fatalf("Added more elements than allowed")
	}
}

func TestAdd(t *testing.T) {
	s := New(0)
	o1 := TestObj{"one", 1}
	o2 := TestObj{"two", 2}
	hash1 := o1.Hash()
	hash2 := o2.Hash()
	s.Add(o1)
	s.Add(o2)
	_, found := s.data[hash1]
	_, found = s.data[hash2]
	if !found {
		t.Fatalf("added but not retrieved")
	}
}

func TestContains(t *testing.T) {
	s := New(0)
	o1 := TestObj{"one", 1}
	o2 := TestObj{"two", 2}
	s.Add(o1)
	s.Add(o2)
	if !s.Contains(o1) || !s.Contains(o2) {
		t.Fatalf("Added but not contained")
	}
}

func TestGet(t *testing.T) {
	s := New(0)
	o1 := TestObj{"one", 1}
	o2 := TestObj{"two", 2}
	s.Add(o1)
	s.Add(o2)
	obj1, err := s.Get(o1.Hash())
	obj2, err := s.Get(o2.Hash())
	if (obj1 == nil) ||
		(obj2 == nil) ||
		(err != nil) {
		t.Fatalf("added but not Get-ed back")
	}
}

func TestRemove(t *testing.T) {
	s := New(0)
	o1 := TestObj{"one", 1}
	o2 := TestObj{"two", 2}
	s.Add(o1)
	s.Add(o2)
	s.Remove(o1.Hash())
	if ne := s.data[o1.Hash()]; ne != nil {
		t.Fatalf("not removed")
	}
	s.Remove(o2.Hash())
	if ne := s.data[o2.Hash()]; ne != nil {
		t.Fatalf("not removed")
	}
}

func TestIterator(t *testing.T) {
	s := New(0)
	o1 := TestObj{"one", 1}
	o2 := TestObj{"two", 2}
	o3 := TestObj{"three", 3}
	s.Add(o1)
	s.Add(o2)
	s.Add(o3)
	for each := range s.Iterator() {
		if each != o1 && each != o2 && each != o3 {
			t.Fatalf("Did not get the whole set")
		}
	}
}
