# recordio for Go

An easy way to read and write records. A record is data wrapped with some
metadata that can delimit a single peice of data from a large set.

The current protocol is simple and fairly naive:

    [magic-number] [compressed-size] [data]
    |<---------- header ---------->|

Overhead per record: 16B (8B + 8B)

### Compression

Uses the [Go port](https://code.google.com/p/snappy-go) of Google's
[snappy](https://code.google.com/p/snappy) compression library.

All records are unconditionally compressed. Snappy might add a few additional
bytes of overhead.

### Usage

See test/recordio_test.go for usage.

### TODOs

1. Integrity check: something like SHA-256 on the compressed data.
2. Better examples.
