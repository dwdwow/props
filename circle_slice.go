package props

type CircleSlice[T any] struct {
	circle         []T
	size           int
	realFirstIndex int
	realLastIndex  int
}

func NewCircleSlice[T any](size int) *CircleSlice[T] {
	return &CircleSlice[T]{circle: make([]T, size), size: size, realFirstIndex: -1, realLastIndex: -1}
}

func (c *CircleSlice[T]) expendRealLastIndex() {
	if c.realLastIndex == -1 {
		c.realFirstIndex = 0
		c.realLastIndex = 0
		return
	}
	c.realLastIndex = (c.realLastIndex + 1) % c.size
	if c.realFirstIndex == c.realLastIndex {
		c.realFirstIndex = (c.realFirstIndex + 1) % c.size
	}
}

func (c *CircleSlice[T]) Size() int {
	return c.size
}

func (c *CircleSlice[T]) Len() int {
	return c.len()
}

func (c *CircleSlice[T]) len() int {
	if c.realFirstIndex == -1 {
		return 0
	}
	return (c.realLastIndex-c.realFirstIndex+c.size)%c.size + 1
}

func (c *CircleSlice[T]) virtualIndexNotInRange(virtual int) bool {
	l := c.len()
	return virtual < -l || virtual >= l
}

func (c *CircleSlice[T]) Push(value T) {
	c.expendRealLastIndex()
	c.circle[c.realLastIndex] = value
}

func (c *CircleSlice[T]) Pop() (t T, ok bool) {
	if c.realFirstIndex == -1 {
		return t, false
	}
	t = c.circle[c.realFirstIndex]
	if c.realFirstIndex == c.realLastIndex {
		c.realFirstIndex = -1
		c.realLastIndex = -1
	} else {
		c.realFirstIndex = (c.realFirstIndex + 1) % c.size
	}
	return t, true
}

func (c *CircleSlice[T]) At(index int) (t T, ok bool) {
	if c.virtualIndexNotInRange(index) {
		return t, false
	}
	idx := (c.realFirstIndex + index + c.size) % c.size
	return c.circle[idx], true
}

func (c *CircleSlice[T]) Set(index int, value T) (ok bool) {
	if c.virtualIndexNotInRange(index) {
		return false
	}
	idx := (c.realFirstIndex + index + c.size) % c.size
	c.circle[idx] = value
	return true
}

func (c *CircleSlice[T]) ForEach(f func(index int, t T)) {
	if c.realFirstIndex == -1 {
		return
	}
	idx := c.realFirstIndex
	for i := 0; i < c.len(); i++ {
		f(i, c.circle[idx])
		idx = (idx + 1) % c.size
	}
}

func (c *CircleSlice[T]) Filter(f func(t T) bool) []T {
	if c.realFirstIndex == -1 {
		return nil
	}
	var result []T
	idx := c.realFirstIndex
	for i := 0; i < c.len(); i++ {
		if f(c.circle[idx]) {
			result = append(result, c.circle[idx])
		}
		idx = (idx + 1) % c.size
	}
	return result
}

func (c *CircleSlice[T]) FiltOne(f func(T) bool) (t T, b bool) {
	if c.realFirstIndex == -1 {
		return t, false
	}
	idx := c.realFirstIndex
	for i := 0; i < c.len(); i++ {
		if f(c.circle[idx]) {
			return c.circle[idx], true
		}
		idx = (idx + 1) % c.size
	}
	return t, false
}

func (c *CircleSlice[T]) Remove(index int) (t T, ok bool) {
	if c.virtualIndexNotInRange(index) {
		return
	}
	idx := (c.realFirstIndex + index + c.size) % c.size
	t = c.circle[idx]
	if c.realFirstIndex == c.realLastIndex {
		c.realFirstIndex = -1
		c.realLastIndex = -1
		return t, true
	}
	nextIdx := (idx + 1) % c.size
	c.circle = append(c.circle[c.realFirstIndex:idx], c.circle[nextIdx:(c.realLastIndex+1)%c.size]...)
	c.realLastIndex = (c.realLastIndex - 1 + c.size) % c.size
	return t, true
}

func (c *CircleSlice[T]) Tidy() {
	if c.realFirstIndex == -1 {
		return
	}
	c.ForEach(func(index int, t T) {
		idx := (c.realFirstIndex + index) % c.size
		c.circle[idx] = *new(T)
	})
}

func (c *CircleSlice[T]) Clear() {
	c.realFirstIndex = -1
	c.realLastIndex = -1
	c.circle = make([]T, c.size)
}
