package pipes

type SinkFunc[T any] func(T)

func Sink[T any](sink SinkFunc[T], in <-chan T) {
	for t := range in {
		sink(t)
	}
}

type SinkWithErrorFunc[T any] func(T) error

func SinkWithError[T any](size int, sink SinkWithErrorFunc[T], in <-chan T) <-chan error {
	err := make(chan error, size)
	go sinkWithErrorWorker(sink, in, err)
	return err
}

func sinkWithErrorWorker[T any](sink SinkWithErrorFunc[T], in <-chan T, err chan<- error) {
	defer close(err)
	for t := range in {
		if er := sink(t); er != nil {
			err <- er
		}
	}
}

func SinkWithErrorSink[T any](sink SinkWithErrorFunc[T], errSink SinkFunc[error], in <-chan T) {
	for t := range in {
		if err := sink(t); err != nil {
			errSink(err)
		}
	}
}
