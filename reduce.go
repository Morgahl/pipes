package pipes

type ReduceFunc[T any, Acc any] func(T, Acc) Acc

// TODO: implement reduce functionality
