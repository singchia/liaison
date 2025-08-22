package frontierbound

import (
	"encoding/json"
	"errors"
	"math/rand"
	"net"

	"github.com/singchia/geminio"
	"github.com/singchia/geminio/client"
	"github.com/singchia/liaison/pkg/edge/config"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/singchia/liaison/pkg/utils"
	"github.com/sirupsen/logrus"
)

type FrontierBound interface {
}

type frontierBound struct {
	end geminio.End
}

func NewFrontierBound(conf *config.Configuration) (FrontierBound, error) {
	dial := conf.Manager.Dial
	if len(dial.Addrs) == 0 {
		return nil, errors.New("dial addr is empty")
	}

	meta := proto.Meta{
		AccessKey: conf.Manager.Auth.AccessKey,
		SecretKey: conf.Manager.Auth.SecretKey,
	}
	data, err := json.Marshal(meta)
	if err != nil {
		return nil, err
	}
	opt := client.NewEndOptions()
	opt.SetMeta(data)

	dialer := func() (net.Conn, error) {
		conn, err := utils.Dial(&dial, rand.Intn(len(dial.Addrs)))
		if err != nil {
			logrus.Errorf("frontlas new informer, dial err: %s", err)
			return nil, err
		}
		return conn, nil
	}
	end, err := client.NewRetryEndWithDialer(dialer, opt)
	if err != nil {
		logrus.Errorf("frontlas new retry end err: %s", err)
		return nil, err
	}

	return &frontierBound{
		end: end,
	}, nil
}
