package pipes

import "github.com/curlymon/pipes/fn"

func Tap[T any](size int, tap fn.Map[T, T], in <-chan T) <-chan T {
	return Map(size, tap, in)
}

func TapWithError[T any](size int, tap fn.MapWithError[T, T], in <-chan T) (<-chan T, <-chan error) {
	return MapWithError(size, tap, in)
}

func TapWithErrorSink[T any](size int, tap fn.MapWithError[T, T], sink fn.Sink[error], in <-chan T) <-chan T {
	return MapWithErrorSink(size, tap, sink, in)
}
