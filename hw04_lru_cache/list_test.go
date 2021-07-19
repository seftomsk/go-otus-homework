package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})
}

// The new cases

func TestRemoveFromEmptyList(t *testing.T) {
	l := NewList()
	item := &ListItem{Value: 10}
	require.Panics(t, func() {
		l.Remove(item)
	})
	require.Equal(t, 0, l.Len())
	require.Nil(t, l.Front())
	require.Nil(t, l.Back())
}

func TestRemoveStrangerItem(t *testing.T) {
	l := NewList()
	l.PushFront(10)
	l.PushBack(20)
	item := &ListItem{Value: 10}
	require.Panics(t, func() {
		l.Remove(item)
	})
	require.Equal(t, 2, l.Len())
	require.Equal(t, 10, l.Front().Value)
	require.Equal(t, 20, l.Back().Value)
}

func TestCallMethodsWithNil(t *testing.T) {
	t.Run("List.Remove", func(t *testing.T) {
		l := NewList()
		require.Panics(t, func() {
			l.Remove(nil)
		})
		require.Equal(t, 0, l.Len())
	})

	t.Run("List.MoveToFront", func(t *testing.T) {
		l := NewList()
		require.Panics(t, func() {
			l.MoveToFront(nil)
		})
		require.Equal(t, 0, l.Len())
	})
}

func createListWithNumbers(count int) List {
	l := NewList()
	for i := 1; i <= count; i++ {
		l.PushBack(i * 10)
	}
	return l
}

func TestRemoveTheFirstOrMiddleOrLastItem(t *testing.T) {
	const count = 4
	t.Run("List.Remove('the first item')", func(t *testing.T) {
		l := createListWithNumbers(count)
		item := l.Front()
		l.Remove(item)
		require.Equal(t, 3, l.Len())
	})

	t.Run("List.Remove('the middle item')", func(t *testing.T) {
		l := createListWithNumbers(count)
		item := l.Front().Next
		l.Remove(item)
		require.Equal(t, 3, l.Len())
	})

	t.Run("List.Remove('the last item')", func(t *testing.T) {
		l := createListWithNumbers(count)
		item := l.Back()
		l.Remove(item)
		require.Equal(t, 3, l.Len())
	})
}

func TestOnlyOneItemIsFirstAndLast(t *testing.T) {
	l := NewList()
	l.PushFront(20)
	require.NotNil(t, l.Front())
	require.NotNil(t, l.Back())
	l = NewList()
	item := &ListItem{Value: 10}
	l.PushFront(item)
	require.NotNil(t, l.Front())
	require.NotNil(t, l.Back())
}
