package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

func (lc *lruCache) Set(key Key, value interface{}) bool {
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
		} else {
			panic("Listitem does not contain cacheItem struct")
		}
	}

	elm := &ListItem{cacheItem{key, value}, nil, nil}
	lc.queue.PushFront(elm)
	lc.items[key] = elm

	return false
}

func (lc *lruCache) Get(key Key) (interface{}, bool) {
	var i interface{}
	flag := false

	if elm, ok := lc.items[key]; ok {
		ci, ok := elm.Value.(cacheItem)
		if ok {
			i = ci.value
		} else {
			panic("Listitem does not contain cacheItem struct")
		}
		lc.queue.MoveToFront(elm)
		flag = true
	}

	return i, flag
}

func (lc *lruCache) Clear() {
	for i := lc.queue.Front(); i != nil; {
		ci, ok := i.Value.(cacheItem)
		if ok {
			delete(lc.items, ci.key)
			lc.queue.Remove(i)
		} else {
			panic("Listitem does not contain cacheItem struct")
		}
	}
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
