package frontierbound

import (
	"context"
	"encoding/json"

	"github.com/jumboframes/armorigo/log"
	"github.com/singchia/geminio"
	"github.com/singchia/liaison/pkg/proto"
)

// reportTrafficMetric 上报流量统计
func (fb *frontierBound) reportTrafficMetric(ctx context.Context, req geminio.Request, rsp geminio.Response) {
	var trafficReq proto.ReportTrafficMetricRequest
	if err := json.Unmarshal(req.Data(), &trafficReq); err != nil {
		log.Errorf("unmarshal traffic metric request error: %s", err)
		rsp.SetError(err)
		return
	}

	// 记录流量到流量统计器（如果存在）
	// 注意：这里需要访问流量统计器，但 frontierBound 目前没有这个字段
	// 我们需要在 NewFrontierBound 中传入流量统计器，或者通过 repo 来记录
	// 暂时先通过 repo 直接记录（每分钟落盘由流量统计器处理）
	// 这里先记录日志，后续可以通过流量统计器来记录
	log.Debugf("received traffic metric: proxyID=%d, applicationID=%d, bytesIn=%d, bytesOut=%d",
		trafficReq.ProxyID, trafficReq.ApplicationID, trafficReq.BytesIn, trafficReq.BytesOut)

	// 通过流量统计器记录流量
	if fb.trafficCollector != nil {
		fb.trafficCollector.RecordTraffic(trafficReq.ProxyID, trafficReq.ApplicationID, trafficReq.BytesIn, trafficReq.BytesOut)
	}
}
