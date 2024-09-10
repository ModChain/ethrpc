package ethrpc

import (
	"encoding/json"
	"fmt"
)

type Request struct {
	JsonRpc string `json:"jsonrpc"` // 2.0
	Method  string `json:"method"`
	Params  []any  `json:"params"`
	Id      any    `json:"id"`
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
