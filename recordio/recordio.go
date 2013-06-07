// Package recordio provides a simple way to read and write records. A record
// is wrapped bytes. The "wrap" is a delimiting header starting with a magic
// number followed by the compressed size of the payload (bytes). The payload
// is compressed using the Snappy compression algorithm.
package recordio

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"

	"code.google.com/p/snappy-go/snappy"
)

const (
	magicNumber int64 = 0x25f9c3e0
)

// Writes records to a writer.
type RecordWriter struct {
	w io.WriteCloser
}

// Creates a new RecordWriter.
func NewRecordWriter(w io.WriteCloser) *RecordWriter {
	return &RecordWriter{w: w}
}

// Writes bytes as a single record. Use multiple Write()s to write multiple
// records before closing the writer.
func (rw *RecordWriter) Write(b []byte) (int, error) {
	buf := new(bytes.Buffer)

	if err := binary.Write(buf, binary.LittleEndian, magicNumber); err != nil {
		return -1, err
	}

	compressed, err := snappy.Encode(nil, b)
	if err != nil {
		return -1, err
	}
	var compressedSize int64 = int64(len(compressed))
	if err := binary.Write(buf, binary.LittleEndian, compressedSize); err != nil {
		return -1, err
	}
	if n, err := buf.Write(compressed); err != nil || int64(n) != compressedSize {
		return -1, err
	}

	return rw.w.Write(buf.Bytes())
}

// Closes the writer.
func (rw *RecordWriter) Close() error {
	return rw.w.Close()
}

// Reads records from a reader.
type RecordReader struct {
	r io.ReadCloser
}

// Creates a new RecordReader.
func NewRecordReader(r io.ReadCloser) *RecordReader {
	return &RecordReader{r: r}
}

// Reads a single record. Use multiple Read()s to read multiple records before
// closing the reader.
func (rr *RecordReader) ReadNext() ([]byte, error) {

	var magic int64
	if err := binary.Read(rr.r, binary.LittleEndian, &magic); err != nil {
		return nil, err
	}
	if magic != magicNumber {
		return nil, errors.New("Invalid record")
	}

	var compressedSize int64
	if err := binary.Read(rr.r, binary.LittleEndian, &compressedSize); err != nil {
		return nil, err
	}

	d := make([]byte, compressedSize)
	if n, err := rr.r.Read(d); err != nil || int64(n) != compressedSize {
		return nil, err
	}

	return snappy.Decode(nil, d)
}

// Closes the reader.
func (rr *RecordReader) Close() error {
	return rr.r.Close()
}
