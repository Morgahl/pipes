package pipes

import "github.com/curlymon/pipes/fn"

type SourceFunc[T any] func() T

func Source[T any](repeat, size int, source SourceFunc[T]) <-chan T {
	out := make(chan T, size)
	go func() {
		defer close(out)
		for i := 0; repeat == RepeatForever || repeat > 0; i++ {
			out <- source()
		}
	}()
	return out
}

type SourceWithErrorFunc[T any] func() (T, error)

func SourceWithError[T any](repeat, size int, source SourceWithErrorFunc[T]) (<-chan T, <-chan error) {
	out, err := make(chan T, size), make(chan error, size)
	go func() {
		defer func() { close(err); close(out) }()
		for i := 0; repeat == RepeatForever || repeat > 0; i++ {
			if v, er := source(); er != nil {
				err <- er
			} else {
				out <- v
			}
		}
	}()
	return out, err
}

func SourceWithErrorSink[T any](repeat, size int, source SourceWithErrorFunc[T], sink fn.Sink[error]) <-chan T {
	out := make(chan T, size)
	go func() {
		defer close(out)
		for i := 0; repeat == RepeatForever || repeat > 0; i++ {
			if v, err := source(); err != nil {
				sink(err)
			} else {
				out <- v
			}
		}
	}()
	return out
}
