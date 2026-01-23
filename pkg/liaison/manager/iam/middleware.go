package iam

import (
	"context"
	"net/http"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	kratoshttp "github.com/go-kratos/kratos/v2/transport/http"
	"github.com/jumboframes/armorigo/log"
)

// AuthMiddleware JWT认证中间件
func AuthMiddleware(iamService *IAMService) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// 从context中获取HTTP请求信息
			if httpReq, ok := kratoshttp.RequestFromServerContext(ctx); ok {
				// 检查是否为IAM相关的接口，这些接口不需要认证
				path := httpReq.URL.Path
				if isIAMEndpoint(path) {
					log.Debugf("Skipping IAM endpoint authentication: %s", path)
					return handler(ctx, req)
				}

				// 获取Authorization头
				authHeader := httpReq.Header.Get("Authorization")
				if authHeader == "" {
					log.Warnf("No authentication token provided")
					return nil, errors.New(http.StatusUnauthorized, "UNAUTHORIZED", "No authentication token provided")
				}

				// 检查Bearer前缀
				tokenParts := strings.SplitN(authHeader, " ", 2)
				if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
					log.Warnf("Invalid token format")
					return nil, errors.New(http.StatusUnauthorized, "UNAUTHORIZED", "Invalid token format")
				}

				tokenString := tokenParts[1]

				// 验证token
				user, err := iamService.GetUserByToken(tokenString)
				if err != nil {
					log.Warnf("Token validation failed: %v", err)
					return nil, errors.New(http.StatusUnauthorized, "UNAUTHORIZED", "Token validation failed")
				}

				// 将用户信息添加到context中
				ctx = context.WithValue(ctx, "user_id", user.ID)
				ctx = context.WithValue(ctx, "user_email", user.Email)
				ctx = context.WithValue(ctx, "user", user)

				log.Debugf("User authentication successful: %s", user.Email)
			}

			return handler(ctx, req)
		}
	}
}

// isIAMEndpoint 判断是否为IAM相关的接口，这些接口不需要认证
func isIAMEndpoint(path string) bool {
	// 不需要认证的接口路径
	noAuthPaths := []string{
		"/api/v1/iam/login",  // 用户登录
		"/api/v1/iam/logout", // 用户登出
		"/install.sh",        // 安装脚本
	}

	// 安装包路径不需要认证
	if strings.HasPrefix(path, "/packages/") {
		return true
	}

	for _, noAuthPath := range noAuthPaths {
		if path == noAuthPath {
			return true
		}
	}

	return false
}
