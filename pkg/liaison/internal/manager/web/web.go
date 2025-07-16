package web

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	"github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/internal/config"
	"github.com/singchia/liaison/pkg/utils"
)

type Web interface {
	Serve() error
	Close() error
}

type web struct {
	app *kratos.App
}

func NewWebServer(conf *config.Configuration) (Web, error) {
	web := &web{}

	listen := &conf.Manager.Listen
	ln, err := utils.Listen(listen)
	if err != nil {
		return nil, err
	}
	opts := []http.ServerOption{
		http.Middleware(recovery.Recovery()),
		http.Listener(ln),
	}
	srv := http.NewServer(opts...)
	v1.RegisterLiaisonServiceHTTPServer(srv, web)

	web.app = kratos.New(
		kratos.Name("liaison"),
		kratos.Server(srv),
	)

	return web, nil
}

func (web *web) Serve() error {
	err := web.app.Run()
	if err != nil {
		return err
	}
	return nil
}

func (web *web) Close() error {
	return web.app.Stop()
}
