package pipes

func Filter[T any](size int, filter func(T) bool, in <-chan T) ChanPull[T] {
	out := make(chan T, size)
	go filterWorker(filter, in, out)
	return out
}

func filterWorker[T any](filter func(T) bool, in <-chan T, out ChanPush[T]) {
	defer close(out)
	for t := range in {
		if filter(t) {
			out <- t
		}
	}
}

func FilterWithError[T any](size int, filter func(T) (bool, error), in <-chan T) (ChanPull[T], ChanPull[error]) {
	out, err := make(chan T, size), make(chan error, size)
	go filterWithErrorWorker(filter, in, out, err)
	return out, err
}

func filterWithErrorWorker[T any](filter func(T) (bool, error), in <-chan T, out chan<- T, err chan<- error) {
	defer func() { close(out); close(err) }()
	for t := range in {
		if keep, er := filter(t); er != nil {
			err <- er
		} else if keep {
			out <- t
		}
	}
}

func FilterWithErrorSink[T any](size int, filter func(T) (bool, error), sink func(error), in <-chan T) ChanPull[T] {
	out := make(chan T, size)
	go filterWithErrorSinkWorker(filter, sink, in, out)
	return out
}

func filterWithErrorSinkWorker[T any](filter func(T) (bool, error), sink func(error), in <-chan T, out chan<- T) {
	defer close(out)
	for t := range in {
		if keep, er := filter(t); er != nil {
			sink(er)
		} else if keep {
			out <- t
		}
	}
}
