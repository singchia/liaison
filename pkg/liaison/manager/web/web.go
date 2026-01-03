package web

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/middleware/recovery"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	v1 "github.com/singchia/liaison/api/v1"
	"github.com/singchia/liaison/pkg/liaison/config"
	"github.com/singchia/liaison/pkg/liaison/manager/controlplane"
	"github.com/singchia/liaison/pkg/liaison/manager/iam"
	"github.com/singchia/liaison/pkg/utils"
)

// @title Liaison Swagger API
// @version 1.0
// @description Liaison Swagger API
// @contact.name Austin Zhai
// @contact.email singchia@163.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
type Web interface {
	Serve() error
	Close() error
}

type web struct {
	app *kratos.App

	// deps
	controlPlane controlplane.ControlPlane
	iamService   *iam.IAMService
}

func NewWebServer(conf *config.Configuration, controlPlane controlplane.ControlPlane, iamService *iam.IAMService) (Web, error) {
	web := &web{
		controlPlane: controlPlane,
		iamService:   iamService,
	}

	listen := &conf.Manager.Listen
	ln, err := utils.Listen(listen)
	if err != nil {
		return nil, err
	}
	// 创建认证中间件
	authMiddleware := iam.AuthMiddleware(web.iamService)

	opts := []kratoshttp.ServerOption{
		kratoshttp.Middleware(recovery.Recovery()),
		kratoshttp.Middleware(authMiddleware),
		kratoshttp.Listener(ln),
	}
	srv := kratoshttp.NewServer(opts...)
	v1.RegisterLiaisonServiceHTTPServer(srv, web)

	// 文件服务
	err = web.serveFiles(conf, srv)
	if err != nil {
		return nil, err
	}

	web.app = kratos.New(
		kratos.Name("liaison"),
		kratos.Server(srv),
	)

	return web, nil
}

func (web *web) serveFiles(conf *config.Configuration, srv *kratoshttp.Server) error {
	// 安装脚本服务
	installScriptPath := filepath.Join("dist", "edge", "install.sh")
	if _, err := os.Stat(installScriptPath); err == nil {
		srv.HandleFunc("/install.sh", func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, installScriptPath)
		})
	}

	// 安装包服务
	packagesDir := conf.Manager.PackagesDir
	if packagesDir == "" {
		packagesDir = "/opt/liaison/packages"
	}
	// 确保 packagesDir 是绝对路径
	packagesDirAbs, err := filepath.Abs(packagesDir)
	if err != nil {
		return err
	}
	packagesPath := filepath.Join(packagesDirAbs, "edge")

	if _, err := os.Stat(packagesPath); err == nil {
		// http.FileServer 使用绝对路径时，内置了路径穿越保护（会拒绝包含 .. 的路径）
		fileServer := http.FileServer(http.Dir(packagesPath))
		srv.HandlePrefix("/packages/edge/", http.StripPrefix("/packages/edge/", fileServer))
	}

	// 前端文件服务
	webDir := conf.Manager.WebDir
	// 确保 webDir 是绝对路径
	webDirAbs, err := filepath.Abs(webDir)
	if err != nil {
		return err
	}
	if _, err := os.Stat(webDirAbs); err != nil {
		return err
	}
	fileServer := http.FileServer(http.Dir(webDirAbs))
	// 前端文件服务：作为 fallback，处理所有非 API、非 install.sh、非 packages 的路径
	srv.HandlePrefix("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		// 如果不是 API、install.sh 和 packages 路径，就走 web 服务
		if !strings.HasPrefix(path, "/api/") &&
			path != "/install.sh" &&
			!strings.HasPrefix(path, "/packages/") {
			// 检查文件是否存在

			filePath := filepath.Join(webDirAbs, path)
			if info, err := os.Stat(filePath); err == nil && !info.IsDir() {
				// 文件存在，直接使用 FileServer
				fileServer.ServeHTTP(w, r)
				return
			}
			// 文件不存在，返回 index.html（SPA 路由）
			http.ServeFile(w, r, filepath.Join(webDirAbs, "index.html"))
			return
		}
		// 其他路径由其他处理器处理，这里不做处理
	}))
	return nil
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
