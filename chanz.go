// Package chanz provides utility functions for channels.
package chanz

import "sync"

// Stream returns a channel with streamed values.
//
// Closes the returned channel when all values are received.
func Stream[T any](xs ...T) <-chan T {
	out := make(chan T, 1)

	go func() {
		defer close(out)

		for i := range xs {
			out <- xs[i]
		}
	}()

	return out
}

// Join returns a channel that joins all passed channels.
//
// Closes the returned channel when all incoming channels are closed.
func Join[T any](ins ...<-chan T) <-chan T {
	out := make(chan T, 1)

	var wg sync.WaitGroup

	for _, ch := range ins {
		if ch == nil {
			continue
		}

		wg.Add(1)

		go func(ch <-chan T) {
			defer wg.Done()

			for val := range ch {
				out <- val
			}
		}(ch)
	}

	go func() {
		defer close(out)

		wg.Wait()
	}()

	return out
}

// Split returns two channels with split values.
//
// First channel contains all the values that satisfy the condition.
//
// Closes both returned channels when the incoming channel is closed.
func Split[T any](in <-chan T, cond func(val T) bool) (<-chan T, <-chan T) {
	var (
		out1 = make(chan T, 1)
		out2 = make(chan T, 1)
	)

	go func() {
		defer close(out2)
		defer close(out1)

		if in == nil {
			return
		}

		for val := range in {
			if cond(val) {
				out1 <- val
			} else {
				out2 <- val
			}
		}
	}()

	return out1, out2
}

// Filter returns a channel with filtered values.
//
// All values that do not satisfy the condition are discarded.
//
// Closes the returned channel when the incoming channel is closed.
func Filter[T any](in <-chan T, cond func(val T) bool) <-chan T {
	out := make(chan T, 1)

	go func() {
		defer close(out)

		if in == nil {
			return
		}

		for val := range in {
			if cond(val) {
				out <- val
			}
		}
	}()

	return out
}

// Map returns a channel with mapped values.
//
// Closes the returned channel when the incoming channel is closed.
func Map[T any](in <-chan T, fn func(val T) T) <-chan T {
	out := make(chan T, 1)

	go func() {
		defer close(out)

		if in == nil {
			return
		}

		for val := range in {
			out <- fn(val)
		}
	}()

	return out
}

// Remap returns a channel with remapped values.
//
// Closes the returned channel when the incoming channel is closed.
func Remap[T, E any](in <-chan T, fn func(val T) E) <-chan E {
	out := make(chan E, 1)

	go func() {
		defer close(out)

		if in == nil {
			return
		}

		for val := range in {
			out <- fn(val)
		}
	}()

	return out
}

// Reduce reduces the values from the channel into the accumulator.
//
// Runs until the incoming channel is closed.
func Reduce[T any](in <-chan T, acc T, fn func(acc T, val T) T) T {
	if in == nil {
		return acc
	}

	for val := range in {
		acc = fn(acc, val)
	}

	return acc
}

// ReduceDefault reduces the values from the channel into the accumulator.
//
// Accumulator value is the default type value.
//
// Runs until the incoming channel is closed.
func ReduceDefault[T any](in <-chan T, fn func(acc T, val T) T) T {
	var acc T

	return Reduce(in, acc, fn)
}

// Collect collects the values from the channel into a slice.
//
// Runs until the incoming channel is closed.
func Collect[T any](in <-chan T) []T {
	if in == nil {
		return nil
	}

	result := make([]T, 0)

	for val := range in {
		result = append(result, val)
	}

	return result
}
