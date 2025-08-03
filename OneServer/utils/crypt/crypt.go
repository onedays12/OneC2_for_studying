package crypt

import (
	"crypto/rc4"
	"errors"
)

func RC4Crypt(data []byte, key []byte) ([]byte, error) {
	rc4crypt, errcrypt := rc4.NewCipher(key)
	if errcrypt != nil {
		return nil, errors.New("rc4 crypt error")
	}
	decryptData := make([]byte, len(data))
	rc4crypt.XORKeyStream(decryptData, data)
	return decryptData, nil
}
