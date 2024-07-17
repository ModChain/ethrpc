package ethrpc

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// ReadUint64 decodes the return value and passes it as a uint64.
//
// This can be used as: res, err := ReadUint64(target.Do("eth_blockNumber"))
func ReadUint64(v any, e error) (uint64, error) {
	if e != nil {
		return 0, e
	}

	switch k := v.(type) {
	case uint64:
		return k, nil
	case string:
		return strconv.ParseUint(k, 0, 64)
	case json.RawMessage:
		var v2 any
		err := json.Unmarshal(k, &v2)
		return ReadUint64(v2, err)
	default:
		return 0, fmt.Errorf("unsupported type %T", v)
	}
}

// ReadString decodes the return value as a string and returns it
func ReadString(v any, e error) (string, error) {
	if e != nil {
		return "", e
	}

	switch k := v.(type) {
	case string:
		return k, nil
	case json.RawMessage:
		var v2 string
		err := json.Unmarshal(k, &v2)
		return v2, err
	default:
		return "", fmt.Errorf("unsupported type %T", v)
	}
}

// ReadTo returns a setter function that will return an error if an error happens. This is
// a bit convoluted because of limitation in Go's syntax, but this could be used as:
//
// err = ReadTo(&block)(target.RPC("eth_getBlockByNumber", "0x1b4", true))
func ReadTo(target any) func(v any, e error) error {
	return func(v any, e error) error {
		if e != nil {
			return e
		}
		return json.Unmarshal(v.(json.RawMessage), target)
	}
}
