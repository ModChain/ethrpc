package ethrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

// Note: see https://eth.wiki/json-rpc/API for APIs
// Note: see https://ethereum.org/en/developers/docs/apis/json-rpc/ for APIs

var rpcId uint64

// TODO support ws protocol
type RPC struct {
	host  string
	lag   time.Duration // how long it takes for this endpoint to respond to eth_blockNumber
	block uint64        // latest block number
}

// New returns a new instance of RPC to perform requests to the given RPC endpoint
func New(h string) *RPC {
	return &RPC{host: h}
}

// Do performs a RPC request
func (r *RPC) Do(method string, args ...any) (json.RawMessage, error) {
	return r.DoCtx(context.Background(), method, args...)
}

// DoCtx performs a RPC request, taking an optional context that can be cancelled to stop the request
func (r *RPC) DoCtx(ctx context.Context, method string, args ...any) (json.RawMessage, error) {
	// JSON RPC is simple
	req := &rpcRequest{
		JsonRpc: "2.0",
		Method:  method,
		Params:  args,
		Id:      atomic.AddUint64(&rpcId, 1),
	}

	//log.Printf("[RPC] → %s %v", method, args)

	reqEnc, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode %s request: %w", method, err)
	}

	hreq, err := http.NewRequestWithContext(ctx, "POST", r.host, bytes.NewReader(reqEnc))
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTTP request for %s: %w", method, err)
	}
	hreq.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(reqEnc)), nil }
	hreq.Header.Set("Content-Type", "application/json")

	// post it
	resp, err := http.DefaultClient.Do(hreq)
	if err != nil {
		return nil, fmt.Errorf("error while performing %s: %w", method, err)
	}
	defer resp.Body.Close()

	// decode response
	reader := json.NewDecoder(resp.Body)
	var res rpcResponse
	err = reader.Decode(&res)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response to %s: %w", method, err)
	}
	if res.Error != nil {
		//log.Printf("[RPC] ← Error: %s", res.Error.Error())
		return nil, fmt.Errorf("RPC error during %s: %w", method, res.Error)
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
