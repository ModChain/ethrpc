package ethrpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync/atomic"
)

type Request struct {
	JsonRpc string `json:"jsonrpc"` // 2.0
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      any    `json:"id"`
}

// NewRequest makes a new [Request] fit to use with methods like Send.
func NewRequest(method string, params ...any) *Request {
	if params == nil {
		// make sure it is not nil so json will encode it as "[]" and not "null"
		params = []any{}
	}
	req := &Request{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
		Id:      atomic.AddUint64(&rpcId, 1),
	}
	return req
}

// HTTPRequest returns a [http.Request] for the given json-rpc request.
func (req *Request) HTTPRequest(ctx context.Context, host string) (*http.Request, error) {
	reqEnc, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to encode %s request: %w", req.Method, err)
	}

	hreq, err := http.NewRequestWithContext(ctx, "POST", host, bytes.NewReader(reqEnc))
	if err != nil {
		return nil, fmt.Errorf("failed to generate HTTP request for %s: %w", req.Method, err)
	}
	hreq.GetBody = func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(reqEnc)), nil }
	hreq.Header.Set("Content-Type", "application/json")

	return hreq, nil
}

type Response struct {
	JsonRpc string          `json:"jsonrpc"` // 2.0
	Result  json.RawMessage `json:"result"`
	Error   *ErrorObject    `json:"error,omitempty"`
	Id      any             `json:"id"`
}

// RPCResponseIntf is same as rpcResponse except Result is a any
type ResponseIntf struct {
	JsonRpc string       `json:"jsonrpc"` // 2.0
	Result  any          `json:"result"`
	Error   *ErrorObject `json:"error,omitempty"`
	Id      any          `json:"id"`
}

type ErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *ErrorObject) Error() string {
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}
