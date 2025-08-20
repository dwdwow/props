package props

import "testing"

func TestCircleSlice(t *testing.T) {
	cs := NewCirculaSlice[int](3)

	// Test initial state
	if cs.Size() != 3 {
		t.Errorf("Expected size 3, got %d", cs.Size())
	}
	if cs.LastOneIndex() != -1 {
		t.Errorf("Expected initial index -1, got %d", cs.LastOneIndex())
	}

	// Test pushing values
	cs.Push(1)
	if cs.LastOneIndex() != 0 {
		t.Errorf("Expected index 0 after first push, got %d", cs.LastOneIndex())
	}

	cs.Push(2)
	cs.Push(3)
	cs.Push(4) // Should overwrite first value

	// Test getting LastOneIndex
	if cs.LastOneIndex() != 0 {
		t.Errorf("Expected index 2 after pushing 4, got %d", cs.LastOneIndex())
	}

	// Test getting values
	val, ok := cs.Get(0) // Should get last inserted value (4)
	if !ok || val != 4 {
		t.Errorf("Expected value 4 at offset 0, got %d", val)
	}

	val, ok = cs.Get(-1) // Should get second to last value (3)
	if !ok || val != 3 {
		t.Errorf("Expected value 3 at offset -1, got %d", val)
	}

	val, ok = cs.Get(-2) // Should get third to last value (2)
	if !ok || val != 2 {
		t.Errorf("Expected value 2 at offset -2, got %d", val)
	}

	val, ok = cs.Get(1) // Should get first value (1)
	if !ok || val != 2 {
		t.Errorf("Expected value 1 at offset 1, got %d", val)
	}

	val, ok = cs.Get(2) // Should get second value (2)
	if !ok || val != 3 {
		t.Errorf("Expected value 2 at offset 2, got %d", val)
	}

	// Test invalid offsets
	_, ok = cs.Get(-3) // At size limit, should be valid
	if ok {
		t.Error("Expected ok for offset -3")
	}

	_, ok = cs.Get(-4) // Beyond size limit, should be invalid
	if ok {
		t.Error("Expected not ok for offset -4")
	}

	_, ok = cs.Get(3) // At size limit, should be valid
	if ok {
		t.Error("Expected ok for offset 3")
	}

	_, ok = cs.Get(4) // Beyond size limit, should be invalid
	if ok {
		t.Error("Expected not ok for offset 4")
	}
}
