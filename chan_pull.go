package pipes

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

func (c ChanPull[T]) Map(size int, mp func(T) any) ChanPull[any] {
	return Map(size, mp, c)
}

func (c ChanPull[T]) MapWithError(size int, mp func(T) (any, error)) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, c)
}

func (c ChanPull[T]) MapWithErrorSink(size int, mp func(T) (any, error), sink func(error)) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, c)
}

func (c ChanPull[T]) Tap(size int, tap func(T) T) ChanPull[T] {
	return Tap(size, tap, c)
}

func (c ChanPull[T]) TapWithError(size int, tap func(T) (T, error)) (ChanPull[T], ChanPull[error]) {
	return TapWithError(size, tap, c)
}

func (c ChanPull[T]) TapWithErrorSink(size int, tap func(T) (T, error), sink func(error)) ChanPull[T] {
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
