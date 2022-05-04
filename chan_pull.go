package pipes

import "time"

type ChanPull[T any] <-chan T

func (c ChanPull[T]) Pull() T {
	return <-c
}

func (c ChanPull[T]) PullSafe() (t T, ok bool) {
	t, ok = <-c
	return
}

func (c ChanPull[T]) TryPull() (t T, ok bool) {
	select {
	case t, ok = <-c:
	default: // chan is empty or closed
	}
	return
}

func (c ChanPull[T]) Drain() {
	for range c {
	}
}

func (c ChanPull[T]) Wait() {
	<-c
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

// Map returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `Map` funciton directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Map(size int, mp func(T) any) ChanPull[any] {
	return Map(size, mp, c)
}

// MapWithError returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `MapWithError` funciton directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) MapWithError(size int, mp func(T) (any, error)) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, c)
}

// MapWithErrorSink returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `MapWithErrorSink` funciton directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) MapWithErrorSink(size int, mp func(T) (any, error), sink func(error)) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, c)
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

func (c ChanPull[T]) Sink(sink func(T)) {
	Sink(sink, c)
}

func (c ChanPull[T]) SinkWithError(size int, sink func(T) error) ChanPull[error] {
	return SinkWithError(size, sink, c)
}

func (c ChanPull[T]) SinkWithErrorSink(sink func(T) error, errSink func(error)) {
	SinkWithErrorSink(sink, errSink, c)
}

func (c ChanPull[T]) FanOut(count, size int) []ChanPull[T] {
	return FanOut(count, size, c)
}

// Reduce returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `Reduce` function directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Reduce(reduce func(T, any) any, acc any) any {
	return Reduce(reduce, acc, c)
}

// ReduceAndEmit returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `ReduceAndEmit` function directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) ReduceAndEmit(reduce func(T, any) any, acc any, in <-chan T) ChanPull[any] {
	return ReduceAndEmit(reduce, acc, c)
}

// Window returns any as the type we transform to here due to generics not supporting method parameterization.
// If you need type safety here use the `Window` function directly.
// https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c ChanPull[T]) Window(size int, window time.Duration, reduce func(T, any) any, acc func() any) ChanPull[any] {
	return Window(size, window, reduce, acc, c)
}
