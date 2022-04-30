package pipes

type CompareFunc[T any, N comparable] func(T) N

func Router[T any, N comparable](size int, matches []N, compare CompareFunc[T, N], in <-chan T) ([]<-chan T, <-chan T) {
	orElse := make(chan T, size)
	outs := make([]<-chan T, len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i := range outs {
		out := make(chan T, size)
		outs[i] = out
		routes[matches[i]] = out
	}
	go routerWorker(compare, in, routes, orElse)
	return outs, orElse
}

func routerWorker[T any, N comparable](compare CompareFunc[T, N], in <-chan T, routes map[N]chan<- T, orElse chan<- T) {
	defer func() {
		for _, out := range routes {
			close(out)
		}
		close(orElse)
	}()
	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
		} else {
			orElse <- t
		}
	}
}

func RouterWithSink[T any, N comparable](size int, matches []N, compare CompareFunc[T, N], sink SinkFunc[T], in <-chan T) []<-chan T {
	outs := make([]<-chan T, len(matches))
	routes := make(map[N]chan<- T, len(matches))
	for i := range outs {
		out := make(chan T, size)
		outs[i] = out
		routes[matches[i]] = out
	}
	go routerWithSinkWorker(compare, in, routes, sink)
	return outs
}

func routerWithSinkWorker[T any, N comparable](compare CompareFunc[T, N], in <-chan T, routes map[N]chan<- T, sink SinkFunc[T]) {
	defer func() {
		for _, out := range routes {
			close(out)
		}
	}()
	for t := range in {
		if route, exists := routes[compare(t)]; exists {
			route <- t
		} else {
			sink(t)
		}
	}
}
