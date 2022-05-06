# Pipes - Type Safe Functional Channels

This library makes use of the newly released Go Generics in `v1.18` to enable the melding of functional programming paradigms, and common Go function call patterns, all over top of Go's excellent channels implementation in a `type safe` manner. Additionally, this library manages the lifecycle of channels and goroutines in a clear and easily reasoned manner, alleviating the details of implementation from the library user.

## Dependencies

This library currently only imports the standard library and has no 3rd party dependencies. The current goal is to maintain this as long as possible. However, I am not opposed to importing 3rd party libraries for testing purposes or if a good case is made for the main libraries and it is not more appropriate for the implementation to exist in another repository. This is intended to be a straightforward implementation only at this time.

## Versioning

This library will make use of semantic versioning and is a work in progress and will be released as `v0.0.1` once:

* A stable-ish API exists in at least the `pipes` package.
* 80-100% test coverage of at least the `pipes` package.
* 80-100% test coverage of all known "gotcha" issues with channels.
* General documentation is available if not explicitly completed.
* Documentation of channel "gotcha" issues with channels.

# Usage

Future home of `README.md` library usage documentation... I'm getting to it. For now follow the soon to be completed first pass of the core `pipes` package.

# Examples

This library includes an `examples` folder containing functional, if sometimes contrived, example implementations of library functionality.

# Contributors and the Curious

Contributions are welcome at this time, however, for the moment, inclusion will be selective in scope to the near term core goals of the project. Once a more stable version is released and we can look at next stages if appropriate.

## Decisions to make (RFC)

### Feel free to open an issue to discuss one of the below decisions, if one does not already exist. This should provide an opportunity for a documented discussion of the issue in a more in depth manner at the time it becomes a requested feature.

1. If our functionality is passed a `nil` `chan` we should "gracefully" skip launching a worker for it, thus preventing accidental blocked reads for no actual useful reason. This should allow for channel lifecycle management as expected.

    * The `Chan[T]`, `ChanPull[T]`, and `ChanPush[T]` should continue to implement the same functionality as expected from standard channel usage in Go on their associated API methods.
    * **This behavior however does not make any effort to warn the library user that the channel was `nil` and that this is likely a bug in their implementation. We should decide whether we should implement this behavior with or without the ability to indicate this issue to the library user.**
    * Document this behavior along with the functional implementations.

2. Parameter ordering on methods and functions should be made consistent and friendly to the library user.

    * Nothing fancy going on here. We just want to make things nice to work with generally.

3. Does any of our implementation need to be moved into an `internal` package?

    * Generally speaking I want to make every package a public API and not make use of `internal`.
    * However this may be a better solution to gathering common worker code patterns as these may only be interesting to any library contributors.

## Decisions made

### Feel free to open an issue to provide an alternate implementation or just a discussion of the issue more in depth, if one does not already exist. This should provide an opportunity for a documented discussion of the issue in a more in depth manner at the time it becomes a requested feature.

1. Some functionality on `Chan[T]` and `ChanPull[T]` cannot be implemented in a type safe manner currently due to the lack of support for parameterized methods. The Functional based implementation should be used instead to ensure type safe implementations.

    References:

    * [Go Generics Proposal - No Parameterized Methods](https://go.googlesource.com/proposal/+/refs/heads/master/design/43651-type-parameters.md#No-parameterized-methods)

    Affected methods at time of writing are below but may not represent a full list:

    * `Map`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's mutated parameter. Use the `Map` function instead if type safety is needed.
    * `MapWithError`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's mutated parameter. Use the `MapWithError` function instead if type safety is needed.
    * `MapWithErrorSink`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's mutated parameter. Use the `MapWithErrorSink` function instead if type safety is needed.
    * `Reduce`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's accumulator parameter. Use the `Reduce` function instead if type safety is needed.
    * `ReduceAndEmit`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's accumulator parameter. Use the `ReduceAndEmit` function instead if type safety is needed.
    * `Window`: currently implemented on `Chan[T]` and `ChanPull[T]` works with `any` as it's accumulator parameter. Use the `Window` function instead if type safety is needed.
    * `Router`: currently *NOT* implemented on any `Chan*[T]` as we cannot easily make use of the `comparable` constraint. Use the `Router` function instead.
    * `RouterWithSink`: currently *NOT* implemented on any `Chan*[T]` as we cannot easily make use of the `comparable` constraint. Use the `RouterWithSink` function instead.

2. `FanIn` will not be able to be used on the `Chan[T]` and `ChanPush[T]` types as `FanIn` as implemented currently will always close the `out` channel. This complicates the reasoning of the channel lifecycle when used from the perspective of `Chan[T]` and `ChanPush[T]`.
