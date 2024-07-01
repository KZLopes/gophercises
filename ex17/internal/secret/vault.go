package secret

import (
	"errors"
	"main/internal/encrypt"
)

type Vault struct {
	encodingKey string
	keyValues   map[string]string
}

func NewMemoryVault(encodingKey string) Vault {
	return Vault{
		encodingKey: encodingKey,
		keyValues:   make(map[string]string),
	}
}

func (v *Vault) Get(key string) (string, error) {
	hex, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for given key")
	}

	ret, err := encrypt.Decrypt(v.encodingKey, hex)
	if err != nil {
		return "", err
	}

	return ret, nil

}

func (v *Vault) Set(key, value string) error {
	encryptedVal, err := encrypt.Encrypt(v.encodingKey, value)
	if err != nil {
		return err
	}
	v.keyValues[key] = encryptedVal
	return nil
}
