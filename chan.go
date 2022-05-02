package pipes

type Chan[T any] chan T

func New[T any](len int) Chan[T] {
	return make(chan T, len)
}

func (c Chan[T]) Close() {
	close(c)
}

func (c Chan[T]) Push(t T) {
	c <- t
}

func (c Chan[T]) TryPush(t T) (ok bool) {
	select {
	case c <- t:
		return true
	default:
		return false
	}
}

func (c Chan[T]) Pull() T {
	return <-c
}

func (c Chan[T]) PullSafe() (t T, ok bool) {
	t, ok = <-c
	return
}

func (c Chan[T]) TryPull() (t T, ok bool) {
	select {
	case t, ok = <-c:
	default: // chan is empty or closed
	}
	return
}

func (c Chan[T]) Drain() {
	for range c {
	}
}

func (c Chan[T]) Wait() {
	<-c
}

func (c Chan[T]) Filter(size int, filter func(T) bool) ChanPull[T] {
	return Filter(size, filter, thunkChanPull(c))
}

func (c Chan[T]) FilterWithError(size int, filter func(T) (bool, error)) (ChanPull[T], ChanPull[error]) {
	return FilterWithError(size, filter, thunkChanPull(c))
}

func (c Chan[T]) FilterWithErrorSink(size int, filter func(T) (bool, error), sink func(error)) ChanPull[T] {
	return FilterWithErrorSink(size, filter, sink, thunkChanPull(c))
}

func (c Chan[T]) Map(size int, mp func(T) any) ChanPull[any] {
	return Map(size, mp, thunkChanPull(c))
}

func (c Chan[T]) MapWithError(size int, mp func(T) (any, error)) (ChanPull[any], ChanPull[error]) {
	return MapWithError(size, mp, thunkChanPull(c))
}

func (c Chan[T]) MapWithErrorSink(size int, mp func(T) (any, error), sink func(error)) ChanPull[any] {
	return MapWithErrorSink(size, mp, sink, thunkChanPull(c))
}

func (c Chan[T]) Tap(size int, tap func(T)) ChanPull[T] {
	return Tap(size, tap, thunkChanPull(c))
}

func (c Chan[T]) TapWithError(size int, tap func(T) error) (ChanPull[T], ChanPull[error]) {
	return TapWithError(size, tap, thunkChanPull(c))
}

func (c Chan[T]) TapWithErrorSink(size int, tap func(T) error, sink func(error)) ChanPull[T] {
	return TapWithErrorSink(size, tap, sink, thunkChanPull(c))
}

func (c Chan[T]) Sink(sink func(T)) {
	Sink(sink, thunkChanPull(c))
}

func (c Chan[T]) SinkWithError(size int, sink func(T) error) ChanPull[error] {
	return SinkWithError(size, sink, thunkChanPull(c))
}

func (c Chan[T]) SinkWithErrorSink(sink func(T) error, errSink func(error)) {
	SinkWithErrorSink(sink, errSink, thunkChanPull(c))
}

func (c Chan[T]) FanOut(count, size int) []ChanPull[T] {
	return FanOut(count, size, thunkChanPull(c))
}

// ChanPush should be a zero cost conversion of Chan[T] to it's ChanPush[T] variant
func (c Chan[T]) ChanPush() ChanPush[T] {
	return thunkChanPush(c)
}

// ChanPull should be a zero cost conversion of Chan[T] to it's ChanPull[T] variant
func (c Chan[T]) ChanPull() ChanPull[T] {
	return thunkChanPull(c)
}

// ChanPull should be a zero cost conversion of Chan[T] to it's ChanPush[T] and ChanPull[T] variants
func (c Chan[T]) ChanPushPull() (ChanPush[T], ChanPull[T]) {
	return thunkChanPush(c), thunkChanPull(c)
}

// thunkChanPull is a thunk that converts the type only and after compile should be optimized out
func thunkChanPull[T any](c chan T) ChanPull[T] {
	return c
}

// thunkChanPush is a thunk that converts the type only and after compile should be optimized out
func thunkChanPush[T any](c chan T) ChanPush[T] {
	return c
}
