package pipes

import "github.com/curlymon/pipes/fn"

func Sink[T any](sink fn.Sink[T], in <-chan T) {
	for t := range in {
		sink(t)
	}
}

func SinkWithError[T any](size int, sink fn.SinkWithError[T], in <-chan T) <-chan error {
	err := make(chan error, size)
	go sinkWithErrorWorker(sink, in, err)
	return err
}

func sinkWithErrorWorker[T any](sink fn.SinkWithError[T], in <-chan T, err chan<- error) {
	defer close(err)
	for t := range in {
		if er := sink(t); er != nil {
			err <- er
		}
	}
}

func SinkWithErrorSink[T any](sink fn.SinkWithError[T], errSink fn.Sink[error], in <-chan T) {
	for t := range in {
		if err := sink(t); err != nil {
			errSink(err)
		}
	}
}
