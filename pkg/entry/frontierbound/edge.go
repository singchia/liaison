package frontierbound

import (
	"context"

	"github.com/singchia/geminio"
)

func (fb *frontierBound) OpenStream(ctx context.Context, edgeID uint64) (geminio.Stream, error) {
	stream, err := fb.svc.OpenStream(ctx, edgeID)
	if err != nil {
		return nil, err
	}
	return stream, nil
}
