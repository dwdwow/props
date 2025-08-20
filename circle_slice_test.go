package props

import "testing"

func TestCircleSlice(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cs := NewCircleSlice[int](3)

		// Test initial state
		if cs.Len() != 0 {
			t.Errorf("Expected length 0, got %d", cs.Len())
		}
		if cs.Size() != 3 {
			t.Errorf("Expected size 3, got %d", cs.Size())
		}

		// Test Push
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)
		if cs.Len() != 3 {
			t.Errorf("Expected length 3, got %d", cs.Len())
		}

		// Test circular behavior
		cs.Push(4) // Should overwrite first element
		if cs.Len() != 3 {
			t.Errorf("Expected length 3, got %d", cs.Len())
		}
		if cs.At(0) != 2 {
			t.Errorf("Expected 2, got %d", cs.At(0))
		}
		if cs.At(1) != 3 {
			t.Errorf("Expected 3, got %d", cs.At(1))
		}
		if cs.At(2) != 4 {
			t.Errorf("Expected 4, got %d", cs.At(2))
		}

		// Test Pop
		val := cs.Pop()
		if val != 2 {
			t.Errorf("Expected 2, got %d", val)
		}
		if cs.Len() != 2 {
			t.Errorf("Expected length 2, got %d", cs.Len())
		}
		if cs.At(0) != 3 {
			t.Errorf("Expected 3, got %d", cs.At(0))
		}
		if cs.At(1) != 4 {
			t.Errorf("Expected 4, got %d", cs.At(1))
		}
	})

	t.Run("get and set", func(t *testing.T) {
		cs := NewCircleSlice[string](2)
		cs.Push("a")
		cs.Push("b")

		if cs.At(0) != "a" {
			t.Errorf("Expected 'a', got %s", cs.At(0))
		}
		if cs.At(1) != "b" {
			t.Errorf("Expected 'b', got %s", cs.At(1))
		}

		cs.Set(1, "c")
		if cs.At(1) != "c" {
			t.Errorf("Expected 'c', got %s", cs.At(1))
		}
	})

	t.Run("forEach operations", func(t *testing.T) {
		cs := NewCircleSlice[int](4)
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)

		sum := 0
		cs.ForEach(func(index int, value int) bool {
			sum += value
			return true
		})
		if sum != 6 {
			t.Errorf("Expected sum 6, got %d", sum)
		}

		sum = 0
		cs.ForEachFromEnd(func(index int, value int) bool {
			sum += value
			return true
		})
		if sum != 6 {
			t.Errorf("Expected sum 6, got %d", sum)
		}
	})

	t.Run("filter operations", func(t *testing.T) {
		cs := NewCircleSlice[int](4)
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)
		cs.Push(4)

		evens := cs.Filter(func(v int) bool {
			return v%2 == 0
		})
		if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
			t.Errorf("Expected [2,4], got %v", evens)
		}

		first, found := cs.FiltOne(func(v int) bool {
			return v > 2
		})
		if !found || first != 3 {
			t.Errorf("Expected 3, got %d, found: %v", first, found)
		}
	})

	t.Run("clear and tidy", func(t *testing.T) {
		cs := NewCircleSlice[int](3)
		cs.Push(1)
		cs.Push(2)

		cs.Clear()
		if cs.Len() != 0 {
			t.Errorf("Expected length 0 after clear, got %d", cs.Len())
		}

		cs.Push(1)
		cs.Push(2)
		cs.Tidy()
		if cs.Len() != 2 {
			t.Errorf("Expected length 2 after tidy, got %d", cs.Len())
		}
	})
}
