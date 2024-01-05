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

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (l *lruCache) Set(key Key, value interface{}) bool {
	if item, ok := l.items[key]; ok {
		// Обновляем значение и перемещаем элемент в начало списка
		item.Value = value
		l.queue.MoveToFront(item)
		return true
	}

	newItem := l.queue.PushFront(value)
	l.items[key] = newItem

	if l.queue.Len() > l.capacity {
		// Удаляем последний элемент из списка и словаря
		lastItem := l.queue.Back()
		if lastItem != nil {
			l.queue.Remove(lastItem)
			for k, v := range l.items {
				if v == lastItem {
					delete(l.items, k)
					break
				}
			}
		}
	}

	return false
}

func (l *lruCache) Get(key Key) (interface{}, bool) {
	if item, ok := l.items[key]; ok {
		l.queue.MoveToFront(item)
		return item.Value, true
	}
	return nil, false
}

func (l *lruCache) Clear() {
	l.queue = NewList()
	l.items = make(map[Key]*ListItem)
}
