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

	t.Run("purge logic", func(t *testing.T) {
		const capacity = 5
		c := NewCache(capacity)
		for i := 1; i <= capacity; i++ {
			c.Set(Key("a"+strconv.Itoa(i)), i*10)
		}
		for i := 1; i <= capacity; i++ {
			val, _ := c.Get(Key("a" + strconv.Itoa(i)))
			require.Equal(t, i*10, val)
		}
		c.Clear()
		for i := 1; i <= capacity; i++ {
			val, ok := c.Get(Key("a" + strconv.Itoa(i)))
			require.Nil(t, val)
			require.False(t, ok)
		}
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

// The new cases

func TestCapacityPermission(t *testing.T) {
	c := NewCache(3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	c.Set("d", 4)
	v, ok := c.Get("a")
	require.Nil(t, v)
	require.False(t, ok)
}

func TestPushOutOfElements(t *testing.T) {
	c := NewCache(3)
	c.Set("a", 1)
	c.Set("b", 2)
	c.Set("c", 3)
	// c = [c=3, b=2, a=1]
	c.Set("b", 10)
	// c = [b=10, c=3, a=1]
	c.Set("c", 20)
	// c = [c=20, b=10, a=1]
	c.Get("a")
	// c = [a=1, c=20, b=10]
	c.Set("d", 40)
	// c = [d=40, a=1, c=20]
	v, ok := c.Get("b")
	require.Nil(t, v)
	require.False(t, ok)
	_, ok = c.Get("a")
	require.True(t, ok)
	_, ok = c.Get("c")
	require.True(t, ok)
	_, ok = c.Get("d")
	require.True(t, ok)
}
