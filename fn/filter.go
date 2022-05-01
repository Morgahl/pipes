package fn

type Filter[T any] func(T) bool

var _ Filter[bool] = Not

func Not(b bool) bool {
	return !b
}

type FilterWithError[T any] func(T) (bool, error)
