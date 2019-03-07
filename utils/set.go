package utils

import "sort"

type StringMapSet struct {
	m map[string]interface{}
}

func NewStringMapSet() *StringMapSet {
	inst := new(StringMapSet)
	inst.m = make(map[string]interface{})
	return inst
}

func StringMapSetUnion(set1, set2 *StringMapSet, sets ...*StringMapSet) *StringMapSet {
	s := NewStringMapSet()
	s.Merge(set1, set2)
	for _, otherSet := range sets {
		s.Merge(otherSet)
	}
	return s
}

func (s *StringMapSet) Add(k string, v interface{}) {
	s.m[k] = v
}

func (s *StringMapSet) AddKey(k string, keys ...string) {
	s.Add(k, nil)
	for _, otherKey := range keys {
		s.Add(otherKey, nil)
	}
}

func (s *StringMapSet) List() []interface{} {
	l := make([]interface{}, 0)
	for _, v := range s.m {
		l = append(l, v)
	}
	return l
}

func (s *StringMapSet) Keys() []string {
	l := make([]string, 0)
	for k, _ := range s.m {
		l = append(l, k)
	}
	return l
}

func (s *StringMapSet) SortedList() []interface{} {
	keys := s.Keys()
	sort.Strings(keys)
	l := make([]interface{}, 0)
	for _, k := range keys {
		l = append(l, s.m[k])
	}
	return l
}

func (s *StringMapSet) Foreach(handler func(k string, v interface{}) bool) {
	for k, v := range s.m {
		if handler(k, v) {
			return
		}
	}
}

func (s *StringMapSet) Get(k string) (interface{}, bool) {
	v, ok := s.m[k]
	return v, ok
}

func (s *StringMapSet) Merge(set *StringMapSet, sets ...*StringMapSet) {
	if set != nil {
		set.Foreach(func(k string, v interface{}) bool {
			s.Add(k, v)
			return false
		})
	}
	for _, otherSet := range sets {
		if otherSet != nil {
			otherSet.Foreach(func(k string, v interface{}) bool {
				s.Add(k, v)
				return false
			})
		}
	}
}
