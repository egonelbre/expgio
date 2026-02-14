package main

type Set map[any]struct{}

func NewSet() Set { return make(Set) }

func (s Set) Len() int    { return len(s) }
func (s Set) Empty() bool { return s.Len() == 0 }

func (s Set) Contains(v any) bool {
	_, ok := s[v]
	return ok
}

func (s Set) Toggle(v any) {
	if s.Contains(v) {
		s.Exclude(v)
	} else {
		s.Include(v)
	}
}

func (s Set) Include(v any) bool {
	if s.Contains(v) {
		return false
	}
	s[v] = struct{}{}
	return true
}

func (s Set) Exclude(v any) bool {
	if s.Contains(v) {
		delete(s, v)
		return true
	}
	return false
}

func (s Set) Set(v any) {
	s.Clear()
	s.Include(v)
}

func (s Set) Clear() {
	for k := range s {
		delete(s, k)
	}
}
