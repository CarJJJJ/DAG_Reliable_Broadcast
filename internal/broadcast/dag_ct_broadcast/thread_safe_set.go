package dagctbroadcast

import (
	"sync"
	"sync/atomic"
)

type ThreadSafeSet struct {
	Set    sync.Map
	length atomic.Int32
}

// NewThreadSafeSet 创建一个新的线程安全集合
func NewThreadSafeSet() *ThreadSafeSet {
	return &ThreadSafeSet{
		length: atomic.Int32{},
	}
}

// Add 添加一个元素到集合中
func (s *ThreadSafeSet) Add(key string) {
	// LoadOrStore 返回: 实际值和是否已存在
	_, loaded := s.Set.LoadOrStore(key, struct{}{})
	// 只有当元素不存在时才增加计数
	if !loaded {
		s.length.Add(1)
	}
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

// Length 返回集合的长度
func (s *ThreadSafeSet) Length() int {
	return int(s.length.Load())
}

// ToSlice 转换为切片
func (s *ThreadSafeSet) ToSlice() []string {
	slice := make([]string, 0)
	s.Set.Range(func(key, value interface{}) bool {
		slice = append(slice, key.(string))
		return true
	})
	return slice
}
