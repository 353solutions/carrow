# carrow - Go bindings to Apache Arrow via C++-API

We're going to implement just a subset of [Arrow](https://arrow.apache.org/) to
make it usable in [frames](https://github.com/v3io/frames).

## FAQ

#### Why note [Apache Arrow for Go](https://github.com/apache/arrow/tree/master/go)?

We'd like to share memory between Go & Python and the current arrow bindings
don't have that option. Since `pyarrow` uses the `C++` Arrow under the hood, we
can just pass a s a pointer.

Also, the C++ Arrow library is more maintained than the Go one and have more
features.
