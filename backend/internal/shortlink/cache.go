package shortlink

import (
	"container/list"
	"sync"
)

// DefaultCacheSize — cache 10k mã trong RAM để phần lớn redirect không chạm DB
// (mục 7 architecture.md: phản hồi < 50ms).
const DefaultCacheSize = 10_000

// Target là dữ liệu tối thiểu cần để trả 302.
type Target struct {
	ArticleID int64
	URL       string
	Alive     bool
}

type cacheEntry struct {
	code   string
	target Target
}

// LRUCache là cache LRU an toàn khi dùng đồng thời.
type LRUCache struct {
	mu       sync.Mutex
	capacity int
	ll       *list.List
	items    map[string]*list.Element
}

// NewLRUCache tạo cache. capacity <= 0 => dùng DefaultCacheSize.
func NewLRUCache(capacity int) *LRUCache {
	if capacity <= 0 {
		capacity = DefaultCacheSize
	}
	return &LRUCache{
		capacity: capacity,
		ll:       list.New(),
		items:    make(map[string]*list.Element, capacity),
	}
}

// Get trả về target đã cache.
func (c *LRUCache) Get(code string) (Target, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	el, ok := c.items[code]
	if !ok {
		return Target{}, false
	}
	c.ll.MoveToFront(el)
	return el.Value.(*cacheEntry).target, true
}

// Put ghi target vào cache, đẩy phần tử cũ nhất ra khi đầy.
func (c *LRUCache) Put(code string, t Target) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if el, ok := c.items[code]; ok {
		el.Value.(*cacheEntry).target = t
		c.ll.MoveToFront(el)
		return
	}
	el := c.ll.PushFront(&cacheEntry{code: code, target: t})
	c.items[code] = el
	for c.ll.Len() > c.capacity {
		oldest := c.ll.Back()
		if oldest == nil {
			break
		}
		c.ll.Remove(oldest)
		delete(c.items, oldest.Value.(*cacheEntry).code)
	}
}

// Len trả về số phần tử đang cache.
func (c *LRUCache) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.ll.Len()
}
