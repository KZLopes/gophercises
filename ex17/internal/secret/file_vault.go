package secret

import (
	"encoding/json"
	"errors"
	"io"
	"io/fs"
	"main/internal/encrypt"
	"os"
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
		keyValues:   make(map[string]string),
	}
}

func (v *FileVault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()

	err := v.load()
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

	err := v.load()
	if err != nil && err != errFileNotFound {
		return err
	}

	v.keyValues[key] = value
	return v.save()
}

func (v *FileVault) load() error {
	if _, err := os.Stat(v.filepath); errors.Is(err, fs.ErrNotExist) {
		return errFileNotFound
	}
	f, err := os.Open(v.filepath)
	if err != nil {
		return err
	}
	defer f.Close()
	r, err := encrypt.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)
}

func (v *FileVault) save() error {
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	w, err := encrypt.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.writeKeyValues(w)
}

func (v *FileVault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

func (v *FileVault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}
