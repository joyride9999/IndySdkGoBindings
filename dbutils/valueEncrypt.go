package dbutils

import (
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
)

const encrypted_key_len = chacha20poly1305.Overhead + chacha20poly1305.NonceSize + chacha20poly1305.KeySize

//SplitValue split the value into key and data
func SplitValue(data []byte) ([]byte, []byte, error) {

	if len(data) < encrypted_key_len {
		return nil, nil, errors.New("invalid structure")
	}

	key := data[:encrypted_key_len]
	value := data[encrypted_key_len:]
	return value, key, nil
}
