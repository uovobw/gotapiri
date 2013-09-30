// Package set tries to implement python-style size-limited set
// semantics in Go. Caveat: the object to be added
// to the set must implement the common.Hashable interface
package set

import (
	"errors"
	"fmt"
)

// type Set uses the data map to contain the common.Hashable objects
// and registers the maximum number of items it can hold in the size variable
type Set struct {
	data map[string]Hashable
	size int
}

// Function New returns a new set. Size must be > 0
func New(size int) (s *Set) {
	s = &Set{}
	s.data = make(map[string]Hashable)
	if size < 0 {
		panic("Cannot create set with negative size")
	} else {
		s.size = size
	}
	return s
}

// Function Add adds an hashable object to the set, if the item is already
// present in the set, the original set is not modified and an error is returned.
// If the addition of this item makes the set larger than its intended size,
// the item is not added and an error is returned instead
func (s *Set) Add(obj Hashable) (result bool, err error) {
	if len(s.data)+1 > s.size && s.size != 0 {
		return false, errors.New(fmt.Sprintf("map size exceeded [%d]", s.size))
	}
	_, present := s.data[obj.Hash()]
	if present {
		return false, nil
	}
	s.data[obj.Hash()] = obj
	return true, nil
}

// Function Contains returns true if a common.Hashable object is contained
// in the set, without returning the object
func (s *Set) Contains(obj Hashable) (present bool) {
	_, present = s.data[obj.Hash()]
	return
}

// Function Get returns, without removing it from the set, a given
// common.Hashable object from the set, identified by its Hash
func (s *Set) Get(id string) (obj Hashable, err error) {
	obj, ok := s.data[id]
	if ok {
		return obj, nil
	} else {
		return nil, errors.New("no object found")
	}
}

// Function Remove deletes the item, without returning it, identified by
// its Hash from the set. Calling Remove on an item that is NOT in the set is
// a no-op
func (s *Set) Remove(id string) {
	delete(s.data, id)
}

// Function Iterator returns a channel of Hashable objects that is
// closed once the last item of the set has been sent through
func (s *Set) Iterator() (ret chan Hashable) {
	ret = make(chan Hashable)
	go func(s *Set, ret chan Hashable) {
		for _, obj := range s.data {
			ret <- obj
		}
		close(ret)
	}(s, ret)
	return ret
}
