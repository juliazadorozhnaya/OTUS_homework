package hw04lrucache //nolint:golint

import (
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
		c := NewCache(3)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache) // Вставка первого элемента
		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache) // Вставка второго элемента
		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache) // Вставка по тому же ключу - замена значения

		wasInCache = c.Set("ddd", 400) // Вставка третьего элемента
		require.False(t, wasInCache)
		val, ok = c.Get("ddd")
		require.True(t, ok)
		require.Equal(t, 400, val)

		wasInCache = c.Set("eee", 500) // Вставка четвертого элемента выше капасити
		require.False(t, wasInCache)
		val, ok = c.Get("eee")
		require.True(t, ok)
		require.Equal(t, 500, val)

		_, ok = c.Get("bbb")
		require.False(t, ok) // Должен отсутствовать и вытолкнуться

		val, ok = c.Get("ccc")
		require.False(t, ok) // Ищем то, чего нет
		require.Nil(t, val)
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

func TestCacheMultithreading(_ *testing.T) {
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
