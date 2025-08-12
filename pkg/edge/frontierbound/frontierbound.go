package frontierbound

import (
	"errors"

	"github.com/singchia/liaison/pkg/edge/config"
)

type FrontierBound interface {
}

type frontierBound struct{}

func NewFrontierBound(conf *config.Configuration) (FrontierBound, error) {
	dial := conf.Manager.Dial
	if len(dial.Addrs) == 0 {
		return nil, errors.New("dial addr is empty")
	}

	return nil, nil
}
