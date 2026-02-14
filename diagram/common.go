package main

type Set map[any]struct{}

func NewSet() Set { return make(Set) }

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
func (s Set) Include(v any) { s[v] = struct{}{} }
func (s Set) Exclude(v any) { delete(s, v) }
