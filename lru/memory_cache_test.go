package lru

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLRUCache(t *testing.T) {
	assert.Panics(t, func() {
		NewLRUCache(-1)
	})

	assert.Panics(t, func() {
		NewLRUCache(0)
	})

	assert.NotPanics(t, func() {
		cache := NewLRUCache(1)
		assert.Equal(t, 1, cache.size)
	})
}

func TestCache_Set(t *testing.T) {
	cache := NewLRUCache(1)
	cache.Set("key1", "val1")
	if cache.evictList.Len() != 1 {
		t.Error("evictList should only include 1 item")
	}

	expect1 := &entry{key: "key1", value: "val1"}
	got1 := cache.evictList.Front().Value.(*entry)
	if !assert.Equal(t, expect1, got1) {
		t.Errorf("evictList should include %v, not %v", expect1, got1)
	}

	cache.Set("key2", "val2")
	if cache.evictList.Len() != 1 {
		t.Error("evictList should only include 1 item")
	}
	expect2 := &entry{key: "key2", value: "val2"}
	got2 := cache.evictList.Front().Value.(*entry)
	if !assert.Equal(t, expect2, got2) {
		t.Errorf("evictList should include %v, not %v", expect2, got2)
	}
}

func TestCache_Get(t *testing.T) {
	cache := NewLRUCache(1)

	val, exist := cache.Get("key")
	assert.Nil(t, val)
	assert.False(t, exist)

	cache.Set("key", "val")
	val, exist = cache.Get("key")
	assert.Equal(t, "val", val)
	assert.True(t, exist)
}
