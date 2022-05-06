package pipes

import "time"

type Chan[T any] chan T

// New returns a new Chan with the given type T. This is essentially a chan T and can be used the
// same way as one would use a channel in Go under normal syntax usages. However this variant has
// methods for functional operations common to the use and lifecycle of channels.
//
// Passing a len of 0 will create an unbuffered channel.
func New[T any](len int) Chan[T] {
	return make(chan T, len)
}

// Close closes the channel. Any attempts to push to a closed channel will panic. Closing an already
// closed channel will return immeadiately.
func (c Chan[T]) Close() {
	close(c)
}

// Push is a blocking operation that pushes a T onto the channel. This blocks while the channel is
// full. This will panic if the channel is closed.
func (c Chan[T]) Push(t T) {
	c <- t
}

// TryPush is a non-blocking operation that attempts to push a T onto the channel. This returns true
// if the T was successfully pushed, false if the channel was blocked. It is exceedingly unlikely
// that you will ever successfully push onto an unbuffered channel as this requires the scheduler to
// have a recieveing goroutine ready.
func (c Chan[T]) TryPush(t T) (ok bool) {
	select {
	case c <- t:
		return true
	default:
		return false
	}
}

// Pull is a blocking operation that pulls a T from the channel if available. This blocks while no T
// is available. If the channel is closed this will return a zero version of the T type.
func (c Chan[T]) Pull() T {
	return <-c
}

// PullSafe is a blocking operation that pulls a T from the channel if available. This returns true
// if the T returned is valid, false if the channel is nil.
func (c Chan[T]) PullSafe() (t T, ok bool) {
	t, ok = <-c
	return
}

// TryPull is a non-blocking operation that attempts to pull a T from the channel. This returns true
// if the T returned is valid, false if the channel is nil or empty.
func (c Chan[T]) TryPull() (t T, ok bool) {
	select {
	case t, ok = <-c:
	default: // chan is empty or closed
	}
	return
}

// Drain is a blocking operation that iterates over the channel discarding values until the channel
// is closed and no further elements remain.
func (c Chan[T]) Drain() {
	for range c {
	}
}

// Wait is a blocking operation that waits for a value to be returned from the channel.
func (c Chan[T]) Wait() {
	<-c
}

func (c Chan[T]) Filter(size int, filter func(T) bool) ChanPull[T] {
	return Filter(size, filter, c)
}

func (c Chan[T]) FilterWithError(size int, filter func(T) (bool, error)) (ChanPull[T], ChanPull[error]) {
	return FilterWithError(size, filter, c)
}

func (c Chan[T]) FilterWithErrorSink(size int, filter func(T) (bool, error), sink func(error)) ChanPull[T] {
	return FilterWithErrorSink(size, filter, sink, c)
}

// Map returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Map` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) Map(size int, mp func(T) any) ChanPull[any] {
	return Map(size, mp, c)
}

// MapWithError returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `MapWithError` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) MapWithError(size int, mp func(T) (any, error)) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, c)
}

// MapWithErrorSink returns any as the type we transform to here due to generics not supporting
// method parameterization. If you need type safety here use the `MapWithErrorSink` function
// directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) MapWithErrorSink(size int, mp func(T) (any, error), sink func(error)) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, c)
}

func (c Chan[T]) Tap(size int, tap func(T)) ChanPull[T] {
	return Tap(size, tap, c)
}

func (c Chan[T]) TapWithError(size int, tap func(T) error) (ChanPull[T], ChanPull[error]) {
	return TapWithError(size, tap, c)
}

func (c Chan[T]) TapWithErrorSink(size int, tap func(T) error, sink func(error)) ChanPull[T] {
	return TapWithErrorSink(size, tap, sink, c)
}

func (c Chan[T]) Sink(sink func(T)) {
	Sink(sink, c)
}

func (c Chan[T]) SinkWithError(size int, sink func(T) error) ChanPull[error] {
	return SinkWithError(size, sink, c)
}

func (c Chan[T]) SinkWithErrorSink(sink func(T) error, errSink func(error)) {
	SinkWithErrorSink(sink, errSink, c)
}

func (c Chan[T]) FanOut(count, size int) []ChanPull[T] {
	return FanOut(count, size, c)
}

func (c Chan[T]) FanIn(ins ...<-chan T) {
	FanInExisting(c, ins...)
}

func (c Chan[T]) RoundRobin(size, count int) []ChanPull[T] {
	return RoundRobin(size, count, c)
}

func (c Chan[T]) Distribute(size, count int, choose func(T) int) []ChanPull[T] {
	return Distribute(size, count, choose, c)
}

// Reduce returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Reduce` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) Reduce(reduce func(T, any) any, acc any) any {
	return Reduce(reduce, acc, c)
}

// ReduceAndEmit returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `ReduceAndEmit` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) ReduceAndEmit(reduce func(T, any) any, acc any, in <-chan T) ChanPull[any] {
	return ReduceAndEmit(reduce, acc, c)
}

// Window returns any as the type we transform to here due to generics not supporting method
// parameterization. If you need type safety here use the `Window` function directly.
//
// ref: https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods
func (c Chan[T]) Window(size int, window time.Duration, reduce func(T, any) any, acc func() any) ChanPull[any] {
	return Window(size, window, reduce, acc, c)
}

// ChanPush is a zero cost conversion of Chan[T] to it's ChanPush[T] variant.
func (c Chan[T]) ChanPush() ChanPush[T] {
	return zcaChanPush(c)
}

// ChanPull is a zero cost conversion of Chan[T] to it's ChanPull[T] variant.
func (c Chan[T]) ChanPull() ChanPull[T] {
	return zcaChanPull(c)
}

// ChanPull is a zero cost conversion of Chan[T] to it's ChanPush[T] and ChanPull[T] variants.
func (c Chan[T]) ChanPushPull() (ChanPush[T], ChanPull[T]) {
	return zcaChanPush(c), zcaChanPull(c)
}

// zcaChanPull is a zero cost abstraction that converts the type only and after compile should be
// optimized out.
func zcaChanPull[T any](c chan T) ChanPull[T] {
	return c
}

// zcaChanPush is a zero cost abstraction that converts the type only and after compile should be
// optimized out.
func zcaChanPush[T any](c chan T) ChanPush[T] {
	return c
}
