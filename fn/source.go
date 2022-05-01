package fn

type Source[T any] func() T

type SourceWithError[T any] func() (T, error)
