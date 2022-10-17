package hw04lrucache

import "sync"

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	sync.Mutex
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
	lc.Lock()
	defer lc.Unlock()
	if elm, ok := lc.items[key]; ok {
		elm.Value = cacheItem{key, value}
		lc.queue.MoveToFront(elm)
		return true
	}

	if lc.queue.Len() >= lc.capacity {
		elm := lc.queue.Back()
		ci, ok := elm.Value.(cacheItem)
		if ok {
			delete(lc.items, ci.key)
			lc.queue.Remove(elm)
		}
	}

	elm := lc.queue.PushFront(cacheItem{key, value})
	lc.items[key] = elm

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	lc.Lock()
	defer lc.Unlock()
	if elm, ok := lc.items[key]; ok {
		ci, ok := elm.Value.(cacheItem)
		if ok {
			lc.queue.MoveToFront(elm)
			return ci.value, true
		}
	}
	return nil, false
}

func (lc *lruCache) Clear() {
	lc.Lock()
	defer lc.Unlock()
	lc.queue = NewList()
	lc.items = make(map[Key]*ListItem, lc.capacity)
}

type cacheItem struct {
	key   Key
	value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
