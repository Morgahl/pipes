package pipes

const RepeatForever = -1

func Source[T any](repeat, size int, source func() T) ChanPull[T] {
	out := make(chan T, size)
	go sourceWorker(repeat, source, out)
	return out
}

func sourceWorker[T any](repeat int, source func() T, out chan<- T) {
	defer close(out)
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		out <- source()
	}
}

func SourceWithError[T any](repeat, size int, source func() (T, error)) (ChanPull[T], ChanPull[error]) {
	out, err := make(chan T, size), make(chan error, size)
	go sourceWithErrorWorker(repeat, source, out, err)
	return out, err
}

func sourceWithErrorWorker[T any](repeat int, source func() (T, error), out chan<- T, err chan<- error) {
	defer func() { close(err); close(out) }()
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		if v, er := source(); er != nil {
			err <- er
		} else {
			out <- v
		}
	}
}

func SourceWithErrorSink[T any](repeat, size int, source func() (T, error), sink func(error)) ChanPull[T] {
	out := make(chan T, size)
	go sourceWithErrorSinkWorker(repeat, source, sink, out)
	return out
}

func sourceWithErrorSinkWorker[T any](repeat int, source func() (T, error), sink func(error), out chan<- T) {
	defer close(out)
	for i := 0; repeat == RepeatForever || repeat > 0; i++ {
		if v, err := source(); err != nil {
			sink(err)
		} else {
			out <- v
		}
	}
}
