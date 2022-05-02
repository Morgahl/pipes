package pipes

func Sink[T any](sink func(T), in <-chan T) {
	for t := range (<-chan T)(in) {
		sink(t)
	}
}

func SinkWithError[T any](size int, sink func(T) error, in <-chan T) ChanPull[error] {
	err := make(chan error, size)
	go sinkWithErrorWorker(sink, in, err)
	return err
}

func sinkWithErrorWorker[T any](sink func(T) error, in <-chan T, err chan<- error) {
	defer close(err)
	for t := range in {
		if er := sink(t); er != nil {
			err <- er
		}
	}
}

func SinkWithErrorSink[T any](sink func(T) error, errSink func(error), in <-chan T) {
	for t := range in {
		if err := sink(t); err != nil {
			errSink(err)
		}
	}
}
