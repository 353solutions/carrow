# carrow - Go bindings to Apache Arrow via C++-API

Access to [Arrow C++](https://arrow.apache.org/docs/cpp/) from Go.

## FAQ

#### Why Not [Apache Arrow for Go](https://github.com/apache/arrow/tree/master/go)?

We'd like to share memory between Go & Python and the current arrow bindings
don't have that option. Since `pyarrow` uses the `C++` Arrow under the hood, we
can just pass a s a pointer.

Also, the C++ Arrow library is more maintained than the Go one and have more
features.
