package chanz_test

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/JFAexe/chanz"
)

func Test_Chanz(t *testing.T) {
	t.Parallel()

	var (
		xs = []int{1, 2, 3, 4, 5, 6}
		os = []int{1, 3, 5}
		es = []int{2, 4, 6}
		ss = []string{"1", "2", "3", "4", "5", "6"}
	)

	var (
		valEven   = func(val int) bool { return val%2 == 0 }
		valPlus   = func(val int) int { return val + 1 }
		valString = func(val int) string { return strconv.Itoa(val) }
		accSum    = func(acc int, val int) int { return acc + val }
	)

	t.Run("Stream", func(t *testing.T) {
		t.Parallel()

		var count int

		for range chanz.Stream(xs...) {
			count++
		}

		require.Equal(t, len(xs), count)
	})

	t.Run("Join", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, 21, chanz.ReduceDefault(chanz.Join(chanz.Stream(os...), chanz.Stream(es...)), accSum))
	})

	t.Run("Split", func(t *testing.T) {
		t.Parallel()

		even, odd := chanz.Split(chanz.Stream(xs...), valEven)

		go func(<-chan int) {
			require.Equal(t, 12, chanz.ReduceDefault(even, accSum))
		}(even)

		require.Equal(t, 9, chanz.ReduceDefault(odd, accSum))
	})

	t.Run("Filter", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, 12, chanz.ReduceDefault(chanz.Filter(chanz.Stream(xs...), valEven), accSum))
	})

	t.Run("Map", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, es, chanz.Collect(chanz.Map(chanz.Stream(os...), valPlus)))
	})

	t.Run("Remap", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, ss, chanz.Collect(chanz.Remap(chanz.Stream(xs...), valString)))
	})

	t.Run("Reduce", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, 42, chanz.Reduce(nil, 42, accSum))
		require.Equal(t, 42, chanz.Reduce(chanz.Stream(xs...), 21, accSum))
	})

	t.Run("ReduceDefault", func(t *testing.T) {
		t.Parallel()

		require.Equal(t, 0, chanz.ReduceDefault(nil, accSum))
		require.Equal(t, 21, chanz.ReduceDefault(chanz.Stream(xs...), accSum))
	})

	t.Run("Collect", func(t *testing.T) {
		t.Parallel()

		require.Nil(t, chanz.Collect[int](nil))
		require.Equal(t, xs, chanz.Collect(chanz.Stream(xs...)))
	})
}
