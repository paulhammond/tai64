# Tai64

Tai64 is a Go implementation of the [TAI64 and TAI64N timestamp formats](http://cr.yp.to/daemontools/tai64n.html).

## Usage

Install with `go get`:

    go get github.com/paulhammond/tai64

Then parse some TAI64N timestamps:

	time, err := tai64.ParseTai64n("@4000000037c219bf2ef02e94")

Full documentation is at http://godoc.org/github.com/paulhammond/tai64.

## License

MIT license, see LICENSE.txt for details.
