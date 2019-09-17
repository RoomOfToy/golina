package container

type void struct{}

// Int Set
type IntSet map[int]void

func NewIntSet(capacity int) *IntSet {
	if capacity < 0 {
		capacity = 0
	}
	s := make(IntSet, capacity)
	return &s
}

func NewIntSetFromIntArray(array []int) *IntSet {
	s := make(IntSet, len(array))
	for _, i := range array {
		s[i] = void{}
	}
	return &s
}

func IntSetEqual(s, s1 *IntSet) bool {
	if s.Length() != s1.Length() || s.Difference(s1).Length() != 0 {
		return false
	}
	return true
}

func (s *IntSet) Get(x int) bool {
	_, ok := (*s)[x]
	return ok
}

func (s *IntSet) Add(x int) {
	(*s)[x] = void{}
}

func (s *IntSet) Delete(x int) {
	delete(*s, x)
}

func (s *IntSet) Length() int {
	return len(*s)
}

func (s *IntSet) Union(s1 *IntSet) *IntSet {
	ns := make(IntSet, s.Length()+s1.Length())
	for k := range *s {
		ns[k] = void{}
	}
	for k := range *s1 {
		ns[k] = void{}
	}
	return &ns
}

func (s *IntSet) Intersection(s1 *IntSet) *IntSet {
	tmp0, tmp1 := s, s1
	if s.Length() > s1.Length() {
		tmp0 = s1
		tmp1 = s
	}
	ns := make(IntSet, tmp0.Length())
	for k := range *tmp0 {
		if _, ok := (*tmp1)[k]; ok {
			ns[k] = void{}
		}
	}
	return &ns
}

func (s *IntSet) Difference(s1 *IntSet) *IntSet {
	ns := s.Union(s1)
	for k := range *(s.Intersection(s1)) {
		ns.Delete(k)
	}
	return ns
}

// Float Set
type FloatSet map[float64]void

func NewFloatSet(capacity int) *FloatSet {
	if capacity < 0 {
		capacity = 0
	}
	s := make(FloatSet, capacity)
	return &s
}

func NewFloatSetFromFloatArray(array []float64) *FloatSet {
	s := make(FloatSet, len(array))
	for _, i := range array {
		s[i] = void{}
	}
	return &s
}

func FloatSetEqual(s, s1 *FloatSet) bool {
	if s.Length() != s1.Length() || s.Difference(s1).Length() != 0 {
		return false
	}
	return true
}

func (s *FloatSet) Get(x float64) bool {
	_, ok := (*s)[x]
	return ok
}

func (s *FloatSet) Add(x float64) {
	(*s)[x] = void{}
}

func (s *FloatSet) Delete(x float64) {
	delete(*s, x)
}

func (s *FloatSet) Length() int {
	return len(*s)
}

func (s *FloatSet) Union(s1 *FloatSet) *FloatSet {
	ns := make(FloatSet, s.Length()+s1.Length())
	for k := range *s {
		ns[k] = void{}
	}
	for k := range *s1 {
		ns[k] = void{}
	}
	return &ns
}

func (s *FloatSet) Intersection(s1 *FloatSet) *FloatSet {
	tmp0, tmp1 := s, s1
	if s.Length() > s1.Length() {
		tmp0 = s1
		tmp1 = s
	}
	ns := make(FloatSet, tmp0.Length())
	for k := range *tmp0 {
		if _, ok := (*tmp1)[k]; ok {
			ns[k] = void{}
		}
	}
	return &ns
}

func (s *FloatSet) Difference(s1 *FloatSet) *FloatSet {
	ns := s.Union(s1)
	for k := range *(s.Intersection(s1)) {
		ns.Delete(k)
	}
	return ns
}
