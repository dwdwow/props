package props

import (
	"sync"
	"testing"
)

func TestSafeRWData(t *testing.T) {
	data := SafeRWData[int]{Data: 42}

	data.Lock()
	if data.Data != 42 {
		t.Errorf("Expected 42, got %d", data.Data)
	}
	data.Data = 24
	data.Unlock()

	data.RLock()
	if data.Data != 24 {
		t.Errorf("Expected 24, got %d", data.Data)
	}
	data.RUnlock()
}

func TestSafeRWSlice(t *testing.T) {
	slice := NewSafeRWSlice[int]()

	// Test Add
	slice.Append(1, 2, 3)
	if len(slice.Data) != 3 {
		t.Errorf("Expected length 3, got %d", len(slice.Data))
	}

	// Test FindOne
	val, found := slice.FindOne(func(i int, d int) bool {
		return d == 2
	})
	if !found || val != 2 {
		t.Error("FindOne failed to find value 2")
	}

	// Test Find
	results := slice.Find(func(i int, d int) bool {
		return d > 1
	})
	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	// Test Update
	updated := slice.Update(func(d int) bool {
		return d < 3
	})
	if len(updated) != 2 {
		t.Errorf("Expected 2 items after update, got %d", len(updated))
	}

	// Test LastOne
	last, ok := slice.LastOne()
	if !ok || last != 2 {
		t.Error("LastOne failed")
	}
}

func TestSafeRWMap(t *testing.T) {
	m := NewSafeRWMap[string, int]()

	// Test SetKV and GetV
	m.SetKV("one", 1)
	if v := m.GetV("one"); v != 1 {
		t.Errorf("Expected 1, got %d", v)
	}

	// Test SetIfNotExists
	currentValue, exists := m.SetIfNotExists("one", 2)
	if !exists || currentValue != 1 {
		t.Errorf("Expected 1, got %d", currentValue)
	}

	currentValue, exists = m.SetIfNotExists("two", 2)
	if exists || currentValue != 2 {
		t.Errorf("Expected 0, got %d", currentValue)
	}

	m.Delete("two")

	// Test GetVWithOk
	if v, ok := m.GetVWithOk("one"); !ok || v != 1 {
		t.Error("GetVWithOk failed")
	}

	// Test GetVWithNew
	v := m.GetVWithNew("two", func() int { return 2 })
	if v != 2 {
		t.Errorf("Expected 2, got %d", v)
	}

	// Test HasKey
	if !m.HasKey("one") {
		t.Error("HasKey failed")
	}

	// Test Delete
	m.Delete("one")
	if m.HasKey("one") {
		t.Error("Delete failed")
	}
	m.Delete("two")
	if m.HasKey("two") {
		t.Error("Delete failed")
	}

	// Test Keys
	m.SetKV("a", 1)
	m.SetKV("b", 2)
	keys := m.Keys()
	if len(keys) != 2 {
		t.Errorf("Expected 2 keys, got %d", len(keys))
	}

	// Test Iterate
	count := 0
	m.Iterate(func(k string, v int) {
		count++
	})
	if count != 2 {
		t.Errorf("Expected 2 iterations, got %d", count)
	}

	// Test Find
	results := m.Find(func(k string, v int) bool {
		return v > 1
	})
	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	// Test DeleteMany
	deleted := m.DeleteMany([]string{"a", "b"})
	if len(deleted) != 2 {
		t.Errorf("Expected 2 deleted items, got %d", len(deleted))
	}
	if m.HasKey("a") || m.HasKey("b") {
		t.Error("DeleteMany failed")
	}

	m.SetKV("a", 1)
	m.SetKV("b", 2)

	// Test DeleteManyWithFilter
	deleted = m.DeleteManyWithFilter(func(k string, v int) bool {
		return v > 1
	})
	if len(deleted) != 1 {
		t.Errorf("Expected 1 deleted item, got %d", len(deleted))
	}
	if m.HasKey("b") {
		t.Error("DeleteManyWithFilter failed")
	}
}

func TestSafeRWCounter(t *testing.T) {
	counter := &SafeRWCounter{}

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			counter.Add()
		}()
	}
	wg.Wait()

	if counter.Data != 100 {
		t.Errorf("Expected counter value 100, got %d", counter.Data)
	}
}
