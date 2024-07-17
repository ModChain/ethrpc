package ethrpc

import (
	"encoding/json"
	"fmt"
)

type rpcRequest struct {
	JsonRpc string `json:"jsonrpc"` // 2.0
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      any    `json:"id"`
}

type rpcResponse struct {
	JsonRpc string          `json:"jsonrpc"` // 2.0
	Result  json.RawMessage `json:"result"`
	Error   *rpcErrorObject `json:"error,omitempty"`
	Id      any             `json:"id"`
}

// rpcResponseIntf is same as rpcResponse except Result is a any
type rpcResponseIntf struct {
	JsonRpc string          `json:"jsonrpc"` // 2.0
	Result  any             `json:"result"`
	Error   *rpcErrorObject `json:"error,omitempty"`
	Id      any             `json:"id"`
}

type rpcErrorObject struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (e *rpcErrorObject) Error() string {
	return fmt.Sprintf("jsonrpc error %d: %s", e.Code, e.Message)
}
