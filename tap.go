package pipes

func Tap[T any](size int, tap func(T) T, in <-chan T) ChanPull[T] {
	return Map(size, tap, in)
}

func TapWithError[T any](size int, tap func(T) (T, error), in <-chan T) (ChanPull[T], ChanPull[error]) {
	return MapWithError(size, tap, in)
}

func TapWithErrorSink[T any](size int, tap func(T) (T, error), sink func(error), in <-chan T) ChanPull[T] {
	return MapWithErrorSink(size, tap, sink, in)
}
