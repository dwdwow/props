package props

type CircleSlice[T any] struct {
	slice []T
	size  int
	index int
}

func NewCirculaSlice[T any](size int) *CircleSlice[T] {
	return &CircleSlice[T]{slice: make([]T, size), size: size, index: -1}
}

func (c *CircleSlice[T]) calIndex(offset int) (i int, ok bool) {
	if offset <= -c.size || offset >= c.size {
		return 0, false
	}
	return (c.index + offset + c.size) % c.size, true
}

func (c *CircleSlice[T]) nextIndex() {
	c.index = (c.index + 1) % c.size
}

func (c *CircleSlice[T]) Size() int {
	return c.size
}

func (c *CircleSlice[T]) LastOneIndex() int {
	return c.index
}

func (c *CircleSlice[T]) Push(value T) int {
	c.nextIndex()
	c.slice[c.index] = value
	return c.index
}

func (c *CircleSlice[T]) Get(offsetAgainstLastOne int) (t T, ok bool) {
	i, ok := c.calIndex(offsetAgainstLastOne)
	if !ok {
		return t, false
	}
	return c.slice[i], true
}
