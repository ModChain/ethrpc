package ethrpc

import (
	"encoding/json"
	"math/big"
)

// ReadUint64 decodes the return value and passes it as a uint64.
//
// This can be used as: res, err := ReadUint64(target.Do("eth_blockNumber"))
func ReadUint64(v json.RawMessage, e error) (uint64, error) {
	if e != nil {
		return 0, e
	}

	var v2 uint64
	err := json.Unmarshal(v, &v2)
	return v2, err
}

// ReadBigInt can decode a json-encoded bigint in various ways, including
// if it is a number literal or a string.
func ReadBigInt(v json.RawMessage, e error) (*big.Int, error) {
	if e != nil {
		return nil, e
	}

	res := new(big.Int)
	err := json.Unmarshal(v, &res)
	return res, err
}

// ReadString decodes the return value as a string and returns it
func ReadString(v json.RawMessage, e error) (string, error) {
	if e != nil {
		return "", e
	}

	var v2 string
	err := json.Unmarshal(v, &v2)
	return v2, err
}

// ReadTo returns a setter function that will return an error if an error happens. This is
// a bit convoluted because of limitation in Go's syntax, but this could be used as:
//
// err = ReadTo(&block)(target.Do("eth_getBlockByNumber", "0x1b4", true))
func ReadTo(target any) func(v json.RawMessage, e error) error {
	return func(v json.RawMessage, e error) error {
		if e != nil {
			return e
		}
		return json.Unmarshal(v, target)
	}
}
