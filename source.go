package pipes

import "github.com/curlymon/pipes/fn"

const RepeatForever = -1

func Source[T any](repeat, size int, source fn.Source[T]) <-chan T {
	out := make(chan T, size)
	go sourceWorker(repeat, source, out)
	return out
}

func sourceWorker[T any](repeat int, source fn.Source[T], out chan<- T) {
	defer close(out)
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		out <- source()
	}
}

func SourceWithError[T any](repeat, size int, source fn.SourceWithError[T]) (<-chan T, <-chan error) {
	out, err := make(chan T, size), make(chan error, size)
	go sourceWithErrorWorker(repeat, source, out, err)
	return out, err
}

func sourceWithErrorWorker[T any](repeat int, source fn.SourceWithError[T], out chan<- T, err chan<- error) {
	defer func() { close(err); close(out) }()
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		if v, er := source(); er != nil {
			err <- er
		} else {
			out <- v
		}
	}
}

func SourceWithErrorSink[T any](repeat, size int, source fn.SourceWithError[T], sink fn.Sink[error]) <-chan T {
	out := make(chan T, size)
	go sourceWithErrorSinkWorker(repeat, source, sink, out)
	return out
}

func sourceWithErrorSinkWorker[T any](repeat int, source fn.SourceWithError[T], sink fn.Sink[error], out chan<- T) {
	defer close(out)
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		if v, err := source(); err != nil {
			sink(err)
		} else {
			out <- v
		}
	}
}
