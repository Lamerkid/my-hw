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

		require.True(t, l.Front().Next == l.Back(), "2 items: First item should be connected to last")
		require.True(t, l.Back().Prev == l.Front(), "2 items: Last item should be connected to first")

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

	t.Run("check items", func(t *testing.T) {
		l := NewList()

		l.PushBack(`world`) // [world]
		l.PushFront("🙃")    // [🙃, world]
		l.PushFront(0x123)  // [291, 🙃, world]
		require.Equal(t, 3, l.Len())

		l.MoveToFront(l.Back())  // [world, 291, 🙃]
		l.MoveToFront(l.Front()) // [world, 291, 🙃]

		require.True(t, l.Front().Prev == nil, "First item should have Prev = nil")
		require.True(t, l.Front().Next == l.Back().Prev, "Middle item should be accesed from first and last")
		require.True(t, l.Back().Next == nil, "Last item should have Next = nil")

		elems := make([]interface{}, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value)
		}
		require.Equal(t, []interface{}{`world`, 0x123, "🙃"}, elems)

		l.Remove(l.Back()) // [world, 291]

		require.True(t, l.Front().Next == l.Back(), "2 items: First item should be connected to last")
		require.True(t, l.Back().Prev == l.Front(), "2 items: Last item should be connected to first")

		l.Remove(l.Front()) // [291]
		l.Remove(l.Front()) // []

		require.True(t, l.Front() == nil, "No items should be in the list")
		require.True(t, l.Back() == nil, "No items should be in the list")
		require.Equal(t, 0, l.Len())
	})
}
