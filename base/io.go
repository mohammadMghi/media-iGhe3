package base

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

type ReadWriterSeeker struct {
	io.LimitedReader
	io.Reader
	io.Writer
	io.Seeker
	io.ByteWriter

	err error
	buf []byte
	n   int
}

func (rws *ReadWriterSeeker) InitializeWriter() {
	if rws.buf == nil {
		rws.buf = make([]byte, 0)
	}
	rws.Writer = bufio.NewWriter(rws)
	var b bytes.Buffer
	rws.ByteWriter = bufio.NewWriter(&b)
}

func (rws *ReadWriterSeeker) WriteTo(w io.Writer) (n int64, err error) {
	n2, err := w.Write(rws.buf)
	return int64(n2), err
}

func (rws *ReadWriterSeeker) Read(p []byte) (n int, err error) {
	return rws.LimitedReader.Read(p)
}

// Flush writes any buffered data to the underlying io.Writer.
func (rws *ReadWriterSeeker) Flush() error {
	if rws.err != nil {
		return rws.err
	}
	if rws.n == 0 {
		return nil
	}
	n, err := rws.Writer.Write(rws.buf[0:rws.n])
	if n < rws.n && err == nil {
		err = io.ErrShortWrite
	}
	if err != nil {
		if n > 0 && n < rws.n {
			copy(rws.buf[0:rws.n-n], rws.buf[n:rws.n])
		}
		rws.n -= n
		rws.err = err
		return err
	}
	rws.n = 0
	return nil
}

func (rws *ReadWriterSeeker) Write(p []byte) (n int, err error) {
	minCap := rws.n + len(p)
	if minCap > cap(rws.buf) {
		buf2 := make([]byte, len(rws.buf), minCap+len(p))
		copy(buf2, rws.buf)
		rws.buf = buf2
	}
	if minCap > len(rws.buf) {
		rws.buf = rws.buf[:minCap]
	}
	copy(rws.buf[rws.n:], p)
	rws.n += len(p)
	return len(p), nil
}

func (rws *ReadWriterSeeker) WriteByte(c byte) error {
	err := rws.ByteWriter.WriteByte(c)
	return err
}

func (rws *ReadWriterSeeker) Seek(offset int64, whence int) (int64, error) {
	newPos, offs := 0, int(offset)
	switch whence {
	case io.SeekStart:
		newPos = offs
	case io.SeekCurrent:
		newPos = rws.n + offs
	case io.SeekEnd:
		newPos = len(rws.buf) + offs
	}
	if newPos < 0 {
		return 0, errors.New("negative result n")
	}
	rws.n = newPos
	return int64(newPos), nil
}

func (rws *ReadWriterSeeker) GetReadSeeker() io.ReadSeeker {
	rws.Reader = bytes.NewReader(rws.buf)
	rws.LimitedReader.R = rws.Reader
	return rws
}

func (rws *ReadWriterSeeker) Close() error {
	return nil
}
