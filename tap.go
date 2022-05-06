package pipes

func Tap[T any](size int, tap func(T), in <-chan T) ChanPull[T] {
	out := make(chan T, size)

	go tapWorker(tap, in, out)

	return out
}

func tapWorker[T any](tap func(T), in <-chan T, out chan<- T) {
	defer close(out)

	for t := range in {
		tap(t)
		out <- t
	}
}

func TapWithError[T any](size int, tap func(T) error, in <-chan T) (ChanPull[T], ChanPull[error]) {
	out, err := make(chan T, size), make(chan error, size)

	go tapWithErrorWorker(tap, in, out, err)

	return out, err
}

func tapWithErrorWorker[T any](tap func(T) error, in <-chan T, out chan<- T, err chan<- error) {
	defer func() { close(out); close(err) }()

	for t := range in {
		if er := tap(t); er != nil {
			err <- er
		}
		out <- t
	}
}

func TapWithErrorSink[T any](size int, tap func(T) error, sink func(error), in <-chan T) ChanPull[T] {
	out := make(chan T, size)

	go tapWithErrorSinkWorker(tap, sink, in, out)

	return out
}

func tapWithErrorSinkWorker[T any](mp func(T) error, sink func(error), in <-chan T, out chan<- T) {
	defer close(out)

	for t := range in {
		if er := mp(t); er != nil {
			sink(er)
		}
		out <- t
	}
}
