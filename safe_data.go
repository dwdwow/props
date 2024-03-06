package props

import (
	"sync"
)

type SafeRWData[D any] struct {
	Data D
	sync.RWMutex
}

func (s *SafeRWData[D]) Set(d D) (old D) {
	s.Lock()
	defer s.Unlock()
	old = s.Data
	s.Data = d
	return
}

func (s *SafeRWData[D]) Get() D {
	s.RLock()
	defer s.RUnlock()
	return s.Data
}

type SafeRWSlice[D any] struct {
	SafeRWData[[]D]
}

func NewSafeRWSlice[D any]() *SafeRWSlice[D] {
	safeSlice := new(SafeRWSlice[D])
	safeSlice.Data = []D{}
	return safeSlice
}

func (slice *SafeRWSlice[D]) Len() int {
	slice.Lock()
	defer slice.Unlock()
	return len(slice.Data)
}

func (slice *SafeRWSlice[D]) Append(u ...D) {
	slice.Lock()
	slice.Data = append(slice.Data, u...)
	slice.Unlock()
}

func (slice *SafeRWSlice[D]) DeleteOne(filter func(d D) bool) (D, bool) {
	slice.Lock()
	defer slice.Unlock()
	for i, v := range slice.Data {
		if filter(v) {
			if i < len(slice.Data)-1 {
				slice.Data = append(slice.Data[:i], slice.Data[i+1:]...)
				return v, true
			} else {
				slice.Data = slice.Data[:i]
				return v, true
			}
		}
	}
	var null D
	return null, false
}

func (slice *SafeRWSlice[D]) GetAt(i int) D {
	slice.RLock()
	defer slice.RUnlock()
	return slice.Data[i]
}

func (slice *SafeRWSlice[D]) SetAt(i int, d D) D {
	slice.Lock()
	defer slice.Unlock()
	old := slice.Data[i]
	slice.Data[i] = d
	return old
}

func (slice *SafeRWSlice[D]) FindOne(filter func(i int, d D) bool) (D, bool) {
	slice.RLock()
	defer slice.RUnlock()
	for i, v := range slice.Data {
		if filter(i, v) {
			return v, true
		}
	}
	var null D
	return null, false
}

func (slice *SafeRWSlice[D]) Find(filter func(i int, d D) bool) []D {
	slice.RLock()
	defer slice.RUnlock()
	var result []D
	for i, v := range slice.Data {
		if filter(i, v) {
			result = append(result, v)
		}
	}
	return result
}

func (slice *SafeRWSlice[D]) Update(filter func(d D) bool) []D {
	slice.Lock()
	defer slice.Unlock()
	var ns []D
	for _, v := range slice.Data {
		if filter(v) {
			ns = append(ns, v)
		}
	}
	slice.Data = ns
	return ns
}

func (slice *SafeRWSlice[D]) LastOne() (D, bool) {
	slice.Lock()
	defer slice.Unlock()
	dl := len(slice.Data)
	if dl <= 0 {
		var d D
		return d, false
	}
	return slice.Data[dl-1], true
}

type SafeRWMap[K comparable, V any] struct {
	SafeRWData[map[K]V]
}

func NewSafeRWMap[K comparable, V any]() *SafeRWMap[K, V] {
	safeMap := new(SafeRWMap[K, V])
	safeMap.Data = map[K]V{}
	return safeMap
}

func (m *SafeRWMap[K, V]) SetKV(k K, v V) (oldValue V) {
	m.Lock()
	defer m.Unlock()
	if m.Data == nil {
		m.Data = map[K]V{}
	}
	oldValue = m.Data[k]
	m.Data[k] = v
	return
}

func (m *SafeRWMap[K, V]) Delete(k K) V {
	m.Lock()
	defer m.Unlock()
	if m.Data == nil {
		m.Data = map[K]V{}
	}
	delete(m.Data, k)
	return m.Data[k]
}

func (m *SafeRWMap[K, V]) GetVWithOk(k K) (V, bool) {
	m.RLock()
	defer m.RUnlock()
	if m.Data == nil {
		m.Data = map[K]V{}
	}
	v, ok := m.Data[k]
	return v, ok
}

func (m *SafeRWMap[K, V]) GetV(k K) V {
	m.RLock()
	defer m.RUnlock()
	if m.Data == nil {
		m.Data = map[K]V{}
	}
	return m.Data[k]
}

type SafeRWCounter struct {
	SafeRWData[int64]
}

func (counter *SafeRWCounter) Add() int64 {
	counter.Lock()
	defer counter.Unlock()
	counter.Data++
	return counter.Data
}
