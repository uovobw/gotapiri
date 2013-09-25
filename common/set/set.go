package set

import (
	"errors"
	"fmt"
)

type Set struct {
	data map[string]Hashable
	size int
}

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

func (s *Set) Contains(obj Hashable) (present bool) {
	_, present = s.data[obj.Hash()]
	return
}

func (s *Set) Get(id string) (obj Hashable, err error) {
	obj, ok := s.data[id]
	if ok {
		return obj, nil
	} else {
		return nil, errors.New("no object found")
	}
}

func (s *Set) Remove(id string) {
	delete(s.data, id)
}
