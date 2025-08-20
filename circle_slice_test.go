package props

import "testing"

func TestCircleSlice(t *testing.T) {
	t.Run("basic operations", func(t *testing.T) {
		cs := NewCirculaSlice[int](5)

		// Test initial state
		if cs.Len() != 0 {
			t.Errorf("Expected length 0, got %d", cs.Len())
		}

		// Test Push and Len
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)
		if cs.Len() != 3 {
			t.Errorf("Expected length 3, got %d", cs.Len())
		}

		// Test At
		if val, ok := cs.At(0); !ok || val != 1 {
			t.Errorf("Expected 1 at index 0, got %d", val)
		}
		if val, ok := cs.At(1); !ok || val != 2 {
			t.Errorf("Expected 2 at index 1, got %d", val)
		}

		// Test Pop
		if val, ok := cs.Pop(); !ok || val != 1 {
			t.Errorf("Expected pop value 1, got %d", val)
		}
		if cs.Len() != 2 {
			t.Errorf("Expected length 2 after pop, got %d", cs.Len())
		}

		// Test circular behavior
		cs.Push(4)
		cs.Push(5)
		cs.Push(6) // Should overwrite the oldest value

		expected := []int{2, 3, 4, 5, 6}
		idx := 0
		cs.Iterate(func(_ int, val int) {
			if val != expected[idx] {
				t.Errorf("Expected %d at position %d, got %d", expected[idx], idx, val)
			}
			idx++
		})
	})

	t.Run("filter operations", func(t *testing.T) {
		cs := NewCirculaSlice[int](5)
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)
		cs.Push(4)

		// Test Filter
		evens := cs.Filter(func(n int) bool {
			return n%2 == 0
		})
		if len(evens) != 2 || evens[0] != 2 || evens[1] != 4 {
			t.Error("Filter didn't return expected even numbers")
		}

		// Test FilterOne
		if val, ok := cs.FiltOne(func(n int) bool {
			return n > 3
		}); !ok || val != 4 {
			t.Error("FilterOne didn't find expected value")
		}
	})

	t.Run("set and remove", func(t *testing.T) {
		cs := NewCirculaSlice[int](5)
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)

		// Test Set
		if ok := cs.Set(1, 10); !ok {
			t.Error("Set operation failed")
		}
		if val, ok := cs.At(1); !ok || val != 10 {
			t.Errorf("Expected 10 at index 1, got %d", val)
		}

		// Test Remove
		if val, ok := cs.Remove(1); !ok || val != 10 {
			t.Errorf("Expected to remove 10, got %d", val)
		}
		if cs.Len() != 2 {
			t.Errorf("Expected length 2 after remove, got %d", cs.Len())
		}
		if val, ok := cs.At(1); !ok || val != 3 {
			t.Errorf("Expected 3 at index 1, got %d", val)
		}
	})

	t.Run("clear and tidy", func(t *testing.T) {
		cs := NewCirculaSlice[int](5)
		cs.Push(1)
		cs.Push(2)
		cs.Push(3)

		cs.Clear()
		if cs.Len() != 0 {
			t.Errorf("Expected length 0 after clear, got %d", cs.Len())
		}

		cs.Push(1)
		cs.Push(2)
		cs.Tidy()
		// After tidy, the slice should still maintain its length
		if cs.Len() != 2 {
			t.Errorf("Expected length 2 after tidy, got %d", cs.Len())
		}
	})
}
