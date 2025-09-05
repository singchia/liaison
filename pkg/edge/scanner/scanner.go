package scanner

import (
	"context"
	"encoding/json"

	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/edge/frontierbound"
	"github.com/singchia/liaison/pkg/proto"
	"github.com/sirupsen/logrus"
)

type Scanner interface{}

type scanner struct {
	frontierBound frontierbound.FrontierBound
}

func NewScanner(frontierBound frontierbound.FrontierBound) (Scanner, error) {

	s := &scanner{
		frontierBound: frontierBound,
	}

	// 注册函数
	err := frontierBound.RegisterRPCHandler("scan_application", s.scan)
	if err != nil {
		logrus.Errorf("scanner register func err: %s", err)
		return nil, err
	}

	return s, nil
}

func (s *scanner) scan(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var task proto.ScanApplicationTaskResult
	err := json.Unmarshal(req.Data(), &task)
	if err != nil {
		rsp.SetError(err)
		return
	}

	// 开始扫描
}

func (s *scanner) Close() error {
	return nil
}
