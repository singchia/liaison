package transport

import (
	"io"
)

// TrafficCounter 流量统计器
type TrafficCounter struct {
	BytesIn  int64
	BytesOut int64
}

// CountingReader 统计读取流量的 Reader
type CountingReader struct {
	reader io.Reader
	counter *TrafficCounter
	isInbound bool // true 表示入站（从客户端到服务器），false 表示出站（从服务器到客户端）
}

func NewCountingReader(reader io.Reader, counter *TrafficCounter, isInbound bool) *CountingReader {
	return &CountingReader{
		reader:    reader,
		counter:   counter,
		isInbound: isInbound,
	}
}

func (cr *CountingReader) Read(p []byte) (n int, err error) {
	n, err = cr.reader.Read(p)
	if n > 0 {
		if cr.isInbound {
			cr.counter.BytesIn += int64(n)
		} else {
			cr.counter.BytesOut += int64(n)
		}
	}
	return n, err
}

// CountingWriter 统计写入流量的 Writer
type CountingWriter struct {
	writer io.Writer
	counter *TrafficCounter
	isInbound bool
}

func NewCountingWriter(writer io.Writer, counter *TrafficCounter, isInbound bool) *CountingWriter {
	return &CountingWriter{
		writer:    writer,
		counter:   counter,
		isInbound: isInbound,
	}
}

func (cw *CountingWriter) Write(p []byte) (n int, err error) {
	n, err = cw.writer.Write(p)
	if n > 0 {
		if cw.isInbound {
			cw.counter.BytesIn += int64(n)
		} else {
			cw.counter.BytesOut += int64(n)
		}
	}
	return n, err
}
