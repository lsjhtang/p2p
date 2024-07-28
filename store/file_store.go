package store

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"io"
	"os"
	"path"
	"strings"
)

func CASPathTransportFromFun(key string) PathKey {
	hash := sha1.Sum([]byte(key))
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize
	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}
	return PathKey{
		PathName: strings.Join(paths, "/"),
		Original: hashStr,
	}
}

type PathKey struct {
	PathName string
	Original string
}

type PathTransportFromFun func(string) PathKey

type Opts struct {
	Root                 string
	PathTransportFromFun PathTransportFromFun
}

type FileStore struct {
	Opts
}

func NewFileStore(opts Opts) *FileStore {
	if opts.PathTransportFromFun == nil {
		opts.PathTransportFromFun = CASPathTransportFromFun
	}
	return &FileStore{Opts: opts}
}

func (s *FileStore) Write(key string, r io.Reader) error {
	pathKey := s.PathTransportFromFun(key)
	return s.writeStream(pathKey, r)
}

func (s *FileStore) writeStream(pathKey PathKey, r io.Reader) error {
	err := os.MkdirAll(s.fullPathDir(pathKey), os.ModePerm)
	fs, err := os.OpenFile(s.fullFileName(pathKey), os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer fs.Close()

	_, err = io.Copy(fs, r)
	return err
}

func (s *FileStore) Has(key string) bool {
	pathKey := s.PathTransportFromFun(key)
	_, err := os.Stat(s.fullFileName(pathKey))
	if err != nil {
		return false
	}
	return true
}

func (s *FileStore) Remove(key string) error {
	pathKey := s.PathTransportFromFun(key)
	return os.Remove(s.fullFileName(pathKey))
}

func (s *FileStore) RemoveAll(key string) error {
	pathKey := s.PathTransportFromFun(key)
	firstPath := pathKey.PathName[:strings.Index(pathKey.PathName, "/")]
	return os.RemoveAll(path.Join(s.Root, firstPath))
}

func (s *FileStore) Read(key string) (io.Reader, error) {
	pathKey := s.PathTransportFromFun(key)
	return s.readStream(pathKey)
}

func (s *FileStore) readStream(pathKey PathKey) (io.Reader, error) {
	fs, err := os.OpenFile(s.fullFileName(pathKey), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return fs, err
	}
	defer fs.Close()
	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, fs)

	return buf, err
}

func (s *FileStore) fullPathDir(pathKey PathKey) string {
	return path.Join(s.Root, pathKey.PathName)
}

func (s *FileStore) fullFileName(pathKey PathKey) string {
	return path.Join(s.fullPathDir(pathKey), pathKey.Original)
}
