package lorm

type iPK interface {
	~int | ~int64 | ~string
}

type Set[PK iPK] map[PK]struct{}

func (s Set[PK]) Add(pk PK) { s[pk] = struct{}{} }

func (s Set[PK]) Contains(pk PK) bool { _, ok := s[pk]; return ok }

func (s Set[PK]) AsList() []PK {
	pks := make([]PK, 0, len(s))
	for pk := range s {
		pks = append(pks, pk)
	}
	return pks
}
