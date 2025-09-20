package proxy

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"strings"
	"sync"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/sirupsen/logrus"
)

type Proxy interface{}

type proxy struct {
	frontierBound frontierbound.FrontierBound
}

func NewProxy(frontierBound frontierbound.FrontierBound) (Proxy, error) {
	proxy := &proxy{
		frontierBound: frontierBound,
	}

	proxy.frontierBound.RegisterStreamHandler(proxy.proxy)

	return proxy, nil
}

func (p *proxy) proxy(ctx context.Context, stream geminio.Stream) {
	// 读取前4个字节获取meta长度
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(stream, lengthBuf)
	if err != nil {
		logrus.Errorf("proxy stream read meta length err: %s", err)
		return
	}
	length := binary.BigEndian.Uint32(lengthBuf)
	dataBuf := make([]byte, length)
	_, err = io.ReadFull(stream, dataBuf)
	if err != nil {
		logrus.Errorf("proxy stream read meta data err: %s", err)
		return
	}

	meta := stream.Meta()
	var dst proto.Dst
	if err := json.Unmarshal(meta, &dst); err != nil {
		logrus.Errorf("proxy stream meta unmarshal err: %s", err)
		return
	}

	conn, err := net.Dial("tcp", dst.Addr)
	if err != nil {
		logrus.Errorf("proxy stream dial err: %s", err)
		return
	}

	wg := new(sync.WaitGroup)
	wg.Add(2)

	go func() {
		defer wg.Done()

		_, err := io.Copy(conn, stream)
		if err != nil && !IsErrClosed(err) {
			logrus.Errorf("read stream, src: %s, dst: %s; to conn, src: %s, dst: %s; err: %s",
				stream.RemoteAddr(), stream.LocalAddr(), conn.LocalAddr(), conn.RemoteAddr(), err)
		}
		_ = stream.Close()
		_ = conn.Close()
	}()

	go func() {
		defer wg.Done()

		_, err := io.Copy(stream, conn)
		if err != nil && !IsErrClosed(err) {
			logrus.Errorf("read conn, src: %s, dst: %s; to stream, src: %s, dst: %s; err: %s",
				conn.LocalAddr(), conn.RemoteAddr(), stream.RemoteAddr(), stream.LocalAddr(), err)
		}
		_ = stream.Close()
		_ = conn.Close()
	}()

	wg.Wait()

	// TODO: some statistics
}

func IsErrClosed(err error) bool {
	if strings.Contains(err.Error(), net.ErrClosed.Error()) {
		return true
	}
	return false
}
