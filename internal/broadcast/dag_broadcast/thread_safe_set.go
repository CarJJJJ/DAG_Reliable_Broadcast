package dagbroadcast

import "sync"

type ThreadSafeSet struct {
	Set sync.Map
}

// NewThreadSafeSet 创建一个新的线程安全集合
func NewThreadSafeSet() *ThreadSafeSet {
	return &ThreadSafeSet{}
}

// Add 添加一个元素到集合中
func (s *ThreadSafeSet) Add(key string) {
	s.Set.Store(key, struct{}{})
}

// Remove 从集合中移除一个元素
func (s *ThreadSafeSet) Remove(value interface{}) {
	s.Set.Delete(value)
}

// Contains 检查集合中是否包含某个元素
func (s *ThreadSafeSet) Contains(value interface{}) bool {
	_, exists := s.Set.Load(value)
	return exists
}
