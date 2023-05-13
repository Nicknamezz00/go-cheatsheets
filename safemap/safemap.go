/*
 * MIT License
 *
 * Copyright (c) 2023 Runze Wu
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package safemap

import (
	"fmt"
	"sync"
)

type SafeMap[K comparable, V any] struct {
	mu   sync.RWMutex
	data map[K]V
}

func NewSafeMap[K comparable, V any]() *SafeMap[K, V] {
	return &SafeMap[K, V]{
		data: make(map[K]V),
		// `mu` initiated with default value
	}
}

func (m *SafeMap[K, V]) Insert(key K, value V) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.data[key] = value
}

func (m *SafeMap[K, V]) Get(key K) (V, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	value, ok := m.data[key]
	if !ok {
		return value, fmt.Errorf("key %v not found", key)
	}
	return value, nil
}

func (m *SafeMap[K, V]) Update(key K, value V) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key %v not found", key)
	}

	m.data[key] = value
	return nil
}

func (m *SafeMap[K, V]) Delete(key K) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ok := m.data[key]
	if !ok {
		return fmt.Errorf("key %v not found", key)
	}

	delete(m.data, key)
	return nil
}

func (m *SafeMap[K, V]) HasKey(key K) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.data[key]
	return ok
}
