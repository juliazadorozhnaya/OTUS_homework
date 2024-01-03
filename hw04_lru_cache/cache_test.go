package hw04lrucache

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("eviction due to size", func(t *testing.T) {
		cacheSize := 3
		c := NewCache(cacheSize)

		// Add elements to fill the cache
		for i := 0; i < cacheSize; i++ {
			c.Set(Key(fmt.Sprintf("key%d", i)), i)
		}

		// Add one more element, which should cause the first one to be evicted
		c.Set("newKey", 100)

		_, ok := c.Get("key0")
		require.False(t, ok) // "key0" should have been evicted

		val, ok := c.Get("newKey")
		require.True(t, ok)
		require.Equal(t, 100, val)
	})

	t.Run("eviction of least recently used", func(t *testing.T) {
		cacheSize := 3
		c := NewCache(cacheSize)

		// Add elements to fill the cache
		for i := 0; i < cacheSize; i++ {
			c.Set(Key(fmt.Sprintf("key%d", i)), i)
		}

		// Access some keys to change their 'recently used' status
		c.Get("key1")      // Access 'key1'
		c.Set("key2", 200) // Update 'key2'

		// Add one more element, which should cause 'key0' to be evicted as it's the least recently used
		c.Set("newKey", 300)

		_, ok := c.Get("key0")
		require.False(t, ok) // "key0" should have been evicted

		val, ok := c.Get("key2")
		require.True(t, ok)
		require.Equal(t, 200, val)

		val, ok = c.Get("newKey")
		require.True(t, ok)
		require.Equal(t, 300, val)
	})

	t.Run("clear cache", func(t *testing.T) {
		c := NewCache(2)
		c.Set("key1", 1)
		c.Set("key2", 2)

		c.Clear()

		_, ok := c.Get("key1")
		require.False(t, ok)

		_, ok = c.Get("key2")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
