package main

type Set map[interface{}]struct{}

func NewSet() Set { return make(Set) }

func (s Set) Contains(v interface{}) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Include(v interface{}) { s[v] = struct{}{} }
func (s Set) Exclude(v interface{}) { delete(s, v) }
