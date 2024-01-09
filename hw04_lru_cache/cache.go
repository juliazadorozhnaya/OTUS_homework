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

type cacheItem struct {
	Key   Key
	Value interface{}
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	if item, ok := l.items[key]; ok {
		// Обновляем значение и перемещаем элемент в начало списка
		item.Value = cacheItem{key, value}
		l.queue.MoveToFront(item)
		return true
	}

	newItem := l.queue.PushFront(cacheItem{key, value})
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		// Удаляем последний элемент из списка и словаря
		lastItem := l.queue.Back()
		if lastItem != nil {
			last := l.queue.Back()
			item := last.Value.(cacheItem)
			l.queue.Remove(last)
			delete(l.items, item.Key)
		}
	}
	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value.(cacheItem).Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.mu.Lock()
	defer l.mu.Unlock()
	for item := range l.items {
		l.queue.Remove(l.items[item])
		delete(l.items, item)
	}
}
