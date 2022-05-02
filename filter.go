package pipes

import "github.com/curlymon/pipes/fn"

func Filter[T any](size int, filter fn.Filter[T], in ChanPull[T]) ChanPull[T] {
	out := make(chan T, size)
	go filterWorker(filter, in, out)
	return out
}

func filterWorker[T any](filter fn.Filter[T], in ChanPull[T], out ChanPush[T]) {
	defer close(out)
	for t := range in {
		if filter(t) {
			out <- t
		}
	}
}

func FilterWithError[T any](size int, filter fn.FilterWithError[T], in <-chan T) (<-chan T, <-chan error) {
	out, err := make(chan T, size), make(chan error, size)
	go filterWithErrorWorker(filter, in, out, err)
	return out, err
}

func filterWithErrorWorker[T any](filter fn.FilterWithError[T], in <-chan T, out chan<- T, err chan<- error) {
	defer func() { close(out); close(err) }()
	for t := range in {
		if keep, er := filter(t); er != nil {
			err <- er
		} else if keep {
			out <- t
		}
	}
}

func FilterWithErrorSink[T any](size int, filter fn.FilterWithError[T], sink fn.Sink[error], in <-chan T) <-chan T {
	out := make(chan T, size)
	go filterWithErrorSinkWorker(filter, sink, in, out)
	return out
}

func filterWithErrorSinkWorker[T any](filter fn.FilterWithError[T], sink fn.Sink[error], in <-chan T, out chan<- T) {
	defer close(out)
	for t := range in {
		if keep, er := filter(t); er != nil {
			sink(er)
		} else if keep {
			out <- t
		}
	}
}
