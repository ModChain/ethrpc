package ethrpc

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Note: see https://eth.wiki/json-rpc/API for APIs
// Note: see https://ethereum.org/en/developers/docs/apis/json-rpc/ for APIs

var rpcId uint64

// TODO support ws protocol
type RPC struct {
	host       string
	lag        time.Duration // how long it takes for this endpoint to respond to eth_blockNumber
	block      uint64        // latest block number
	HTTPClient *http.Client
}

// New returns a new instance of RPC to perform requests to the given RPC endpoint
func New(h string) *RPC {
	return &RPC{host: h, HTTPClient: http.DefaultClient}
}

// Do performs a RPC request
func (r *RPC) Do(method string, args ...any) (json.RawMessage, error) {
	return r.DoCtx(context.Background(), method, args...)
}

func (r *RPC) Send(req *Request) (json.RawMessage, error) {
	return r.SendCtx(context.Background(), req)
}

// DoCtx performs a RPC request, taking an optional context that can be cancelled to stop the request
func (r *RPC) DoCtx(ctx context.Context, method string, args ...any) (json.RawMessage, error) {
	return r.SendCtx(ctx, NewRequest(method, args...))
}

func (r *RPC) SendCtx(ctx context.Context, req *Request) (json.RawMessage, error) {
	// JSON RPC over http is simple
	//log.Printf("[RPC] → %s %v", method, args)

	hreq, err := req.HTTPRequest(ctx, r.host)
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTTP request for %s: %w", req.Method, err)
	}

	// post it
	resp, err := r.HTTPClient.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("error while performing %s: %w", req.Method, err)
	}
	defer resp.Body.Close()

	// decode response
	reader := json.NewDecoder(resp.Body)
	var res *Response
	err = reader.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response to %s: %w", req.Method, err)
	}
	if res.Error != nil {
		//log.Printf("[RPC] ← Error: %s", res.Error.Error())
		return nil, fmt.Errorf("RPC error during %s: %w", req.Method, res.Error)
	}

	//log.Printf("[RPC] ← %s", res.Result)

	return res.Result, nil
}

func (r *RPC) To(target any, method string, args ...any) error {
	v, err := r.Do(method, args...)
	if err != nil {
		return err
	}
	return json.Unmarshal(v, target)
}
