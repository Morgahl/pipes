package pipes

import "github.com/curlymon/pipes/fn"

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

func (c ChanPull[T]) Filter(size int, filter fn.Filter[T]) ChanPull[T] {
	return Filter(size, filter, c)
}

func (c ChanPull[T]) FilterWithError(size int, filter fn.FilterWithError[T]) (ChanPull[T], ChanPull[error]) {
	return FilterWithError(size, filter, c)
}

func (c ChanPull[T]) FilterWithErrorSink(size int, filter fn.FilterWithError[T], sink fn.Sink[error]) ChanPull[T] {
	return FilterWithErrorSink(size, filter, sink, c)
}

func (c ChanPull[T]) Map(size int, mp fn.Map[T, any]) ChanPull[any] {
	return Map(size, mp, c)
}

func (c ChanPull[T]) MapWithError(size int, mp fn.MapWithError[T, any]) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, c)
}

func (c ChanPull[T]) MapWithErrorSink(size int, mp fn.MapWithError[T, any], sink fn.Sink[error]) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, c)
}

func (c ChanPull[T]) Tap(size int, tap fn.Map[T, T]) ChanPull[T] {
	return Tap(size, tap, c)
}

func (c ChanPull[T]) TapWithError(size int, tap fn.MapWithError[T, T]) (ChanPull[T], ChanPull[error]) {
	return TapWithError(size, tap, c)
}

func (c ChanPull[T]) TapWithErrorSink(size int, tap fn.MapWithError[T, T], sink fn.Sink[error]) ChanPull[T] {
	return TapWithErrorSink(size, tap, sink, c)
}

func (c ChanPull[T]) Sink(sink fn.Sink[T]) {
	Sink(sink, c)
}

func (c ChanPull[T]) SinkWithError(size int, sink fn.SinkWithError[T]) ChanPull[error] {
	return SinkWithError(size, sink, c)
}

func (c ChanPull[T]) SinkWithErrorSink(sink fn.SinkWithError[T], errSink fn.Sink[error]) {
	SinkWithErrorSink(sink, errSink, c)
}

func (c ChanPull[T]) FanOut(count, size int) []ChanPull[T] {
	return FanOut(count, size, c)
}
