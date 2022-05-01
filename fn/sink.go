package fn

type Sink[T any] func(T)

type SinkWithError[T any] func(T) error
