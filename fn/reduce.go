package fn

type Reduce[T any, Acc any] func(T, Acc) Acc

type ReduceWithError[T any, Acc any] func(T, Acc) (Acc, error)
