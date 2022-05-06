package pipes

import "time"

type ChanPull[T any] <-chan T

// Pull is a blocking operation that pulls a T from the channel if available. This blocks while no T
// is available. If the channel is closed and empty, or nil, this will return a zero version of the
// T type.
func (c ChanPull[T]) Pull() T {
	return <-c
}

// PullSafe is a blocking operation that pulls a T from the channel if available. This returns true
// if the T returned is valid, false if the channel is closed and empty, or nil.
func (c ChanPull[T]) PullSafe() (t T, ok bool) {
	t, ok = <-c
	return
}

// TryPull is a non-blocking operation that attempts to pull a T from the channel. This returns true
// if the T returned is valid, false if the channel is closed and empty, or nil.
func (c ChanPull[T]) TryPull() (t T, ok bool) {
	select {
	case t, ok = <-c:
	default: // chan is empty or closed
	}
	return
}

// Drain is a blocking operation that iterates over the channel discarding values until the channel
// is closed and no further elements remain. This returns immeadiately if the channel is closed or
// nil.
func (c ChanPull[T]) Drain() {
	for range c {
	}
}

// Wait is a blocking operation that waits for a value to be returned from the channel. If the
// channel is closed or nil this will immeadiately return.
func (c ChanPull[T]) Wait() {
	<-c
}

func (c ChanPull[T]) FanOut(count, size int) []ChanPull[T] {
	return FanOut(count, size, c)
}

func (c ChanPull[T]) Filter(size int, filter func(T) bool) ChanPull[T] {
	return Filter(size, filter, c)
}

func (c ChanPull[T]) FilterWithError(size int, filter func(T) (bool, error)) (ChanPull[T], ChanPull[error]) {
	return FilterWithError(size, filter, c)
}

func (c ChanPull[T]) FilterWithErrorSink(size int, filter func(T) (bool, error), sink func(error)) ChanPull[T] {
	return FilterWithErrorSink(size, filter, sink, c)
}

// Map returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Map` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Map(size int, mp func(T) any) ChanPull[any] {
	return Map(size, mp, c)
}

// MapWithError returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `MapWithError` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) MapWithError(size int, mp func(T) (any, error)) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, c)
}

// MapWithErrorSink returns any as the type we transform to here due to generics not supporting
// method parameterization. If you need type safety here use the `MapWithErrorSink` function
// directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) MapWithErrorSink(size int, mp func(T) (any, error), sink func(error)) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, c)
}

// Reduce returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Reduce` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Reduce(reduce func(T, any) any, acc any) any {
	return Reduce(reduce, acc, c)
}

// ReduceAndEmit returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `ReduceAndEmit` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) ReduceAndEmit(reduce func(T, any) any, acc any, in <-chan T) ChanPull[any] {
	return ReduceAndEmit(reduce, acc, c)
}

// Window returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Window` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Window(size int, window time.Duration, reduce func(T, any) any, acc func() any) ChanPull[any] {
	return Window(size, window, reduce, acc, c)
}

func (c ChanPull[T]) RoundRobin(size, count int) []ChanPull[T] {
	return RoundRobin(size, count, c)
}

func (c ChanPull[T]) Distribute(size, count int, choose func(T) int) []ChanPull[T] {
	return Distribute(size, count, choose, c)
}

func (c ChanPull[T]) Sink(sink func(T)) {
	Sink(sink, c)
}

func (c ChanPull[T]) SinkWithError(size int, sink func(T) error) ChanPull[error] {
	return SinkWithError(size, sink, c)
}

func (c ChanPull[T]) SinkWithErrorSink(sink func(T) error, errSink func(error)) {
	SinkWithErrorSink(sink, errSink, c)
}

func (c ChanPull[T]) Tap(size int, tap func(T)) ChanPull[T] {
	return Tap(size, tap, c)
}

func (c ChanPull[T]) TapWithError(size int, tap func(T) error) (ChanPull[T], ChanPull[error]) {
	return TapWithError(size, tap, c)
}

func (c ChanPull[T]) TapWithErrorSink(size int, tap func(T) error, sink func(error)) ChanPull[T] {
	return TapWithErrorSink(size, tap, sink, c)
}
