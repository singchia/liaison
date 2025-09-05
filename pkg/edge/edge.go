package edge

import (
	"github.com/singchia/liaison/pkg/edge/config"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/edge/proxy"
	"github.com/singchia/liaison/pkg/edge/scanner"
)

type Edge struct {
	frontierBound frontierbound.FrontierBound
}

func NewEdge() (*Edge, error) {

	err := config.Init()
	if err != nil {
		return nil, err
	}

	frontierBound, err := frontierbound.NewFrontierBound(config.Conf)
	if err != nil {
		return nil, err
	}

	_, err = proxy.NewProxy(frontierBound)
	if err != nil {
		return nil, err
	}

	_, err = scanner.NewScanner(frontierBound)
	if err != nil {
		return nil, err
	}

	return &Edge{
		frontierBound: frontierBound,
	}, nil
}

func (e *Edge) Close() error {
	return e.frontierBound.Close()
}
