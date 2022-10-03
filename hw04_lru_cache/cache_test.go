package hw04lrucache

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

	t.Run("check capacity", func(t *testing.T) {
		c := NewCache(3)
		b := c.Set("i1", 1)
		require.False(t, b)

		b = c.Set("i2", 2)
		require.False(t, b)

		b = c.Set("i3", 3)
		require.False(t, b)

		_, _ = c.Get("i1")
		_, _ = c.Get("i2")

		b = c.Set("i4", 4)
		require.False(t, b)

		_, ok := c.Get("i3")
		require.False(t, ok)
	})
	t.Run("check func clear", func(t *testing.T) {
		c := NewCache(3)
		b := c.Set("i1", 1)
		require.False(t, b)

		b = c.Set("i2", 2)
		require.False(t, b)

		b = c.Set("i3", 3)
		require.False(t, b)

		c.Clear()

		_, b = c.Get("i1")
		require.False(t, b)
		_, b = c.Get("i2")
		require.False(t, b)
		_, b = c.Get("i3")
		require.False(t, b)

		b = c.Set("i1", 1)
		require.False(t, b)

		b = c.Set("i2", 2)
		require.False(t, b)

		b = c.Set("i3", 3)
		require.False(t, b)
	})
}

func TestCacheMultithreading(t *testing.T) {

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
