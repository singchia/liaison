package controlplane

import "github.com/singchia/liaison/pkg/liaison/internal/repo"

// @title Liaison Swagger API
// @version 1.0
// @description Liaison Swagger API
// @contact.name Austin Zhai
// @contact.email singchia@163.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

type controlPlane struct {
	repo repo.Repo
}
