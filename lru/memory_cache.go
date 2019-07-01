package lru

import (
	"container/list"
	"errors"
)

var _ Cacher = (*Cache)(nil)

type Cacher interface {
	Set(key, value interface{})
	Get(key interface{}) (val interface{}, exist bool)
}

func NewLRUCache(size int) *Cache {
	if size <= 0 {
		panic(errors.New("cache size should be positive"))
	}

	return &Cache{
		size:      size,
		items:     make(map[interface{}]*list.Element, size),
		evictList: list.New(),
	}
}

type Cache struct {
	size      int
	items     map[interface{}]*list.Element
	evictList *list.List
}

type entry struct {
	key   interface{}
	value interface{}
}

func (c *Cache) Set(key, value interface{}) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value = value
		return
	}

	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	if c.evictList.Len() <= c.size {
		return
	}

	c.removeOldest()
}

func (c *Cache) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

func (c *Cache) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
}

func (c *Cache) Get(key interface{}) (val interface{}, exist bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return nil, false
}
