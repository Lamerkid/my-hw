package hw04lrucache

import "sync"

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
	mu       sync.Mutex
}

// Структура, которая позволит связать элемент очереди и кэша.
type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, v interface{}) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Инициализируем новый элемент.
	newItem := &cacheItem{key: key, value: v}
	// Если ключ есть в кэше, то двигаем значение в начало очереди.
	if _, ok := c.items[key]; ok {
		c.items[key] = c.queue.PushFront(newItem)
		return true
	}
	c.items[key] = c.queue.PushFront(newItem)
	// Если превышена емкость, то удаляем последнее значение из очереди
	// и элемент кэша по ключу.
	if c.capacity < c.queue.Len() {
		lastItem := c.queue.Back()
		c.queue.Remove(lastItem)
		delete(c.items, lastItem.Value.(*cacheItem).key)
	}
	return false
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// Если ключа нет в кэше, возвращаем nil, false.
	if c.items[key] == nil {
		return nil, false
	}
	val := c.items[key].Value.(*cacheItem).value
	c.queue.MoveToFront(c.items[key])
	return val, true
}

func (c *lruCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
