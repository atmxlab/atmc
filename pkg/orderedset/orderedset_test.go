package orderedset_test

import (
	"testing"

	"github.com/atmxlab/atmc/pkg/orderedset"
	"github.com/stretchr/testify/require"
)

func TestOrderedSet_SetAndGet(t *testing.T) {
	t.Parallel()

	t.Run("set_and_get", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		for i := 0; i < 10; i++ {
			v, ok := os.Get(i)
			require.True(t, ok)
			require.Equal(t, i+1, v)
		}
	})

	t.Run("get_not_exist_key", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		v, ok := os.Get(5000)
		require.False(t, ok)
		require.Zero(t, v)
	})
}

func TestOrderedSet_Delete(t *testing.T) {
	t.Parallel()

	t.Run("delete", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		os.Delete(1)

		_, ok := os.Get(1)
		require.False(t, ok)

		v, ok := os.Get(0)
		require.True(t, ok)
		require.Equal(t, 0+1, v)
	})

	t.Run("delete_not_exist_key", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		os.Delete(5000)
	})

	t.Run("delete_last_element", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		os.Set(0, 1)

		os.Delete(0)

		_, ok := os.Get(0)
		require.False(t, ok)
	})
}

func TestOrderedSet_Iterate(t *testing.T) {
	t.Parallel()

	t.Run("iterate", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		i := 0
		for k, v := range os.Iterator() {
			require.Equal(t, i, k)
			require.Equal(t, i+1, v)
			i++
		}
	})

	t.Run("iterate_after_override", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		os.Set(1, 2)

		i := 0
		for k, v := range os.Iterator() {
			require.Equal(t, i, k)
			require.Equal(t, i+1, v)
			i++
		}
	})
}

func TestOrderedSet_KeysAndValues(t *testing.T) {
	t.Parallel()

	t.Run("values", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		expectedValues := make([]int, 0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
			expectedValues = append(expectedValues, i+1)
		}

		require.Equal(t, expectedValues, os.Values())
	})

	t.Run("values_after_delete", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		deleteIndex := 4
		expectedValues := make([]int, 0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
			if i != deleteIndex {
				expectedValues = append(expectedValues, i+1)
			}
		}

		os.Delete(deleteIndex)

		require.Equal(t, expectedValues, os.Values())
	})

	t.Run("keys", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		expectedKeys := make([]int, 0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
			expectedKeys = append(expectedKeys, i)
		}

		require.Equal(t, expectedKeys, os.Keys())
	})

	t.Run("values_after_delete", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		deleteIndex := 4
		expectedKeys := make([]int, 0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
			if i != deleteIndex {
				expectedKeys = append(expectedKeys, i)
			}
		}

		os.Delete(deleteIndex)

		require.Equal(t, expectedKeys, os.Keys())
	})
}

func TestOrderedSet_Params(t *testing.T) {
	t.Parallel()

	t.Run("len", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		require.Equal(t, 10, os.Len())
	})

	t.Run("len_after_delete", func(t *testing.T) {
		t.Parallel()

		os := orderedset.New[int, int](0)

		for i := 0; i < 10; i++ {
			os.Set(i, i+1)
		}

		os.Delete(2)
		os.Delete(5)

		require.Equal(t, 8, os.Len())
	})
}
