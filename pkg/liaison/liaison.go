package liaison

import (
	"fmt"
	"net/http"
	"runtime"

	frontierutils "github.com/singchia/frontier/pkg/utils"
	"github.com/singchia/liaison/pkg/entry"
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
	"github.com/singchia/liaison/pkg/liaison/manager/frontierbound"
	"github.com/singchia/liaison/pkg/liaison/manager/iam"
	"github.com/singchia/liaison/pkg/liaison/manager/web"
	"github.com/singchia/liaison/pkg/liaison/repo"
	"github.com/singchia/liaison/pkg/utils"
	"k8s.io/klog/v2"
)

type Liaison struct {
	web           web.Web
	frontierBound frontierbound.FrontierBound
	entry         *entry.Entry
	repo          repo.Repo
	iamService    *iam.IAMService
}

func NewLiaison() (*Liaison, error) {
	err := config.Init()
	if err != nil {
		return nil, err
	}
	// pprof & rlimit
	if config.Conf.Daemon.PProf.Enable {
		runtime.SetCPUProfileRate(config.Conf.Daemon.PProf.CPUProfileRate)
		go func() {
			http.ListenAndServe(config.Conf.Daemon.PProf.Addr, nil)
		}()
	}
	// rlimit
	if config.Conf.Daemon.RLimit.Enable {
		err = frontierutils.SetRLimit(uint64(config.Conf.Daemon.RLimit.NumFile))
		if err != nil {
			klog.Errorf("set rlimit err: %s", err)
			return nil, err
		}
	}
	// repo
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
	controlPlane, err := controlplane.NewControlPlane(config.Conf, repo, frontierBound)
	if err != nil {
		return nil, err
	}
	// IAM service
	iamService := iam.NewIAMService(repo)
	// 设置JWT密钥（必须从配置文件读取）
	if config.Conf.Manager.JWTSecret == "" {
		return nil, fmt.Errorf("JWT secret key is required in configuration file. Please set 'manager.jwt_secret' in your configuration")
	}
	if err := utils.SetJWTSecret(config.Conf.Manager.JWTSecret); err != nil {
		return nil, fmt.Errorf("failed to set JWT secret: %w", err)
	}
	// web layer
	web, err := web.NewWebServer(config.Conf, controlPlane, iamService)
	if err != nil {
		return nil, err
	}
	// entry layer
	entry, err := entry.NewEntry(config.Conf, controlPlane)
	if err != nil {
		return nil, err
	}
	return &Liaison{
		web:           web,
		frontierBound: frontierBound,
		entry:         entry,
		repo:          repo,
		iamService:    iamService,
	}, nil
}

func (l *Liaison) Serve() error {
	return l.web.Serve()
}

func (l *Liaison) Close() error {
	err := l.web.Close()
	if err != nil {
		return err
	}
	err = l.frontierBound.Close()
	if err != nil {
		return err
	}
	err = l.entry.Close()
	if err != nil {
		return err
	}
	err = l.repo.Close()
	if err != nil {
		return err
	}
	return nil
}
