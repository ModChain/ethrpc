package ethrpc

import (
	"context"
	"encoding/json"
)

type Handler interface {
	DoCtx(ctx context.Context, method string, args ...any) (json.RawMessage, error)
}

type Api struct {
	Handler
}

func (a *Api) Do(method string, args ...any) (json.RawMessage, error) {
	return a.Handler.DoCtx(context.Background(), method, args...)
}

func (a *Api) To(target any, method string, args ...any) error {
	return a.ToCtx(context.Background(), target, method, args...)
}

func (a *Api) ToCtx(ctx context.Context, target any, method string, args ...any) error {
	v, err := a.DoCtx(ctx, method, args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, target)
}

func (a *Api) BlockNumber(ctx context.Context) (uint64, error) {
	return ReadUint64(a.Handler.DoCtx(ctx, "eth_blockNumber"))
}

func (a *Api) ChainId(ctx context.Context) (uint64, error) {
	return ReadUint64(a.Handler.DoCtx(ctx, "eth_chainId"))
}
