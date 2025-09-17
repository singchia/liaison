package liaison

import (
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
	"github.com/singchia/liaison/pkg/liaison/manager/frontierbound"
	"github.com/singchia/liaison/pkg/liaison/manager/web"
	"github.com/singchia/liaison/pkg/liaison/repo"
)

type Liaison struct {
	web           web.Web
	controlPlane  controlplane.ControlPlane
	frontierBound frontierbound.FrontierBound
}

func NewLiaison() (*Liaison, error) {
	err := config.Init()
	if err != nil {
		return nil, err
	}
	repo, err := repo.NewRepo(config.Conf)
	if err != nil {
		return nil, err
	}
	// frontier bound
	frontierBound, err := frontierbound.NewFrontierBound(config.Conf, repo)
	if err != nil {
		return nil, err
	}
	// service layer
	controlPlane, err := controlplane.NewControlPlane(repo, frontierBound)
	if err != nil {
		return nil, err
	}
	// web layer
	web, err := web.NewWebServer(config.Conf, controlPlane)
	if err != nil {
		return nil, err
	}
	// entry layer
	return &Liaison{
		web:           web,
		controlPlane:  controlPlane,
		frontierBound: frontierBound,
	}, nil
}

func (l *Liaison) Serve() error {
	return l.web.Serve()
}

func (l *Liaison) Close() error {
	return l.web.Close()
}
