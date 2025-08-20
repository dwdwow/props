package props

type CircleSlice[T any] struct {
	circle    []T
	size      int
	realStart int
	realEnd   int
}

func NewCircleSlice[T any](size int) *CircleSlice[T] {
	return &CircleSlice[T]{circle: make([]T, size), size: size, realStart: -1, realEnd: -1}
}

func (c *CircleSlice[T]) expendRealEnd() {
	if c.realEnd == -1 {
		c.realStart = 0
		c.realEnd = 0
		return
	}
	c.realEnd = (c.realEnd + 1) % c.size
	if c.realStart == c.realEnd {
		c.realStart = (c.realStart + 1) % c.size
	}
}

func (c *CircleSlice[T]) Size() int {
	return c.size
}

func (c *CircleSlice[T]) Len() int {
	return c.len()
}

func (c *CircleSlice[T]) len() int {
	if c.realStart == -1 {
		return 0
	}
	return (c.realEnd-c.realStart+c.size)%c.size + 1
}

func (c *CircleSlice[T]) virtualIndexNotInRange(virtual int) bool {
	l := c.len()
	return virtual < -l || virtual >= l
}

func (c *CircleSlice[T]) Push(value T) {
	c.expendRealEnd()
	c.circle[c.realEnd] = value
}

func (c *CircleSlice[T]) Pop() (t T) {
	if c.realStart == -1 {
		panic("circle slice is empty")
	}
	t = c.circle[c.realStart]
	if c.realStart == c.realEnd {
		c.realStart = -1
		c.realEnd = -1
	} else {
		c.realStart = (c.realStart + 1) % c.size
	}
	return t
}

func (c *CircleSlice[T]) At(index int) (t T) {
	if c.virtualIndexNotInRange(index) {
		panic("index out of circle slice range")
	}
	idx := (c.realStart + index + c.size) % c.size
	return c.circle[idx]
}

func (c *CircleSlice[T]) Set(index int, value T) {
	if c.virtualIndexNotInRange(index) {
		panic("index out of circle slice range")
	}
	idx := (c.realStart + index + c.size) % c.size
	c.circle[idx] = value
}

func (c *CircleSlice[T]) ForEach(f func(index int, t T) (goOn bool)) {
	if c.realStart == -1 {
		return
	}
	idx := c.realStart
	for i := 0; i < c.len(); i++ {
		if !f(i, c.circle[idx]) {
			return
		}
		idx = (idx + 1) % c.size
	}
}

func (c *CircleSlice[T]) ForEachFromEnd(f func(index int, t T) (goOn bool)) {
	if c.realStart == -1 {
		return
	}
	idx := c.realEnd
	for i := c.len() - 1; i >= 0; i-- {
		if !f(i, c.circle[idx]) {
			return
		}
		idx = (idx - 1 + c.size) % c.size
	}
}

func (c *CircleSlice[T]) Filter(f func(t T) bool) []T {
	if c.realStart == -1 {
		return nil
	}
	var result []T
	idx := c.realStart
	for i := 0; i < c.len(); i++ {
		if f(c.circle[idx]) {
			result = append(result, c.circle[idx])
		}
		idx = (idx + 1) % c.size
	}
	return result
}

func (c *CircleSlice[T]) FiltOne(f func(T) bool) (t T, b bool) {
	if c.realStart == -1 {
		return t, false
	}
	idx := c.realStart
	for i := 0; i < c.len(); i++ {
		if f(c.circle[idx]) {
			return c.circle[idx], true
		}
		idx = (idx + 1) % c.size
	}
	return t, false
}

func (c *CircleSlice[T]) Remove(index int) (t T) {
	if c.virtualIndexNotInRange(index) {
		panic("index out of circle slice range")
	}
	idx := (c.realStart + index + c.size) % c.size
	t = c.circle[idx]
	if c.realStart == c.realEnd {
		c.realStart = -1
		c.realEnd = -1
		return t
	}
	nextIdx := (idx + 1) % c.size
	c.circle = append(c.circle[c.realStart:idx], c.circle[nextIdx:(c.realEnd+1)%c.size]...)
	c.realEnd = (c.realEnd - 1 + c.size) % c.size
	return t
}

func (c *CircleSlice[T]) Tidy() {
	if c.realStart == -1 {
		return
	}
	c.ForEach(func(index int, t T) bool {
		idx := (c.realStart + index) % c.size
		c.circle[idx] = *new(T)
		return true
	})
}

func (c *CircleSlice[T]) Clear() {
	c.realStart = -1
	c.realEnd = -1
	c.circle = make([]T, c.size)
}
