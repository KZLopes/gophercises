package secret

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"main/internal/encrypt"
	"os"
	"strings"
	"sync"
)

var errFileNotFound = errors.New("vault: file not found")

type FileVault struct {
	encodingKey string
	filepath    string
	keyValues   map[string]string

	mutex sync.Mutex
}

func NewFileVault(encodingKey, filepath string) *FileVault {
	return &FileVault{
		encodingKey: encodingKey,
		filepath:    filepath,
		// keyValues:   make(map[string]string),
	}
}

func (v *FileVault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil {
		return "", err
	}

	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for given key")
	}

	return value, nil
}

func (v *FileVault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.loadKeyValues()
	if err != nil && err != errFileNotFound {
		return err
	}

	v.keyValues[key] = value
	err = v.saveKeyValues()
	if err != nil {
		return err
	}
	return nil
}

func (v *FileVault) loadKeyValues() error {
	if _, err := os.Stat(v.filepath); errors.Is(err, fs.ErrNotExist) {
		v.keyValues = make(map[string]string)
		return errFileNotFound
	}
	f, err := os.Open(v.filepath)
	if err != nil {
		return nil
	}
	defer f.Close()

	var sb strings.Builder
	_, err = io.Copy(&sb, f)
	if err != nil {
		return err
	}
	decryptedJSON, err := encrypt.Decrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}
	r := strings.NewReader(decryptedJSON)
	dec := json.NewDecoder(r)
	err = dec.Decode(&v.keyValues)
	if err != nil {
		return err
	}
	return nil
}

func (v *FileVault) saveKeyValues() error {
	var sb strings.Builder

	enc := json.NewEncoder(&sb)
	err := enc.Encode(v.keyValues)
	if err != nil {
		return err
	}

	encryptedJSON, err := encrypt.Encrypt(v.encodingKey, sb.String())
	if err != nil {
		return err
	}

	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(encryptedJSON)
	if err != nil {
		return err
	}
	//
	// 	err = os.WriteFile(v.filepath, []byte(encryptedJSON), 7001)
	// 	if err != nil {
	// 		return err
	// 	}

	return nil
}
