package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

func initLRUCache(cache *lruCache) {
	cache.items = make(map[Key]*ListItem, cache.capacity)
	cache.queue = NewList()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   string
	value interface{}
}

func NewCache(capacity int) Cache {
	cache := &lruCache{capacity: capacity}
	initLRUCache(cache)

	return cache
}

func (cache *lruCache) getCacheItem(elem *ListItem) *cacheItem {
	return elem.Value.(*cacheItem)
}

func (cache *lruCache) Set(key Key, value interface{}) bool {
	cItem := &cacheItem{
		key:   string(key),
		value: value,
	}
	if item, ok := cache.items[key]; ok {
		item.Value = cItem
		cache.queue.MoveToFront(item)
		return true
	}

	item := cache.queue.PushFront(cItem)
	cache.items[key] = item

	if cache.queue.Len() <= cache.capacity {
		return false
	}

	if lastItem := cache.queue.Back(); lastItem != nil {
		cache.queue.Remove(lastItem)
		delete(cache.items, Key(cache.getCacheItem(lastItem).key))
	}

	return false
}

func (cache *lruCache) Get(key Key) (interface{}, bool) {
	if _, ok := cache.items[key]; !ok {
		return nil, false
	}

	item := cache.items[key]
	cache.queue.PushFront(item)

	return cache.getCacheItem(item).value, true
}

func (cache *lruCache) Clear() {
	initLRUCache(cache)
}
