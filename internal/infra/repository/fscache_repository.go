package repository

import (
	"time"

	fscache "github.com/iqquee/fs-cache"
)

type fscacheServer struct {
	memdis *fscache.Memdis
}

func NewFsCahe(client *fscache.Memdis) *fscacheServer {
	return &fscacheServer{
		memdis: client,
	}
}

func (f *fscacheServer) Set(key string, value any, expiration int32) (string, error) {

	return "", nil
}

func (f *fscacheServer) Increment(key string, expiration time.Duration) (int64, error) {

	return 0, nil
}

func (f *fscacheServer) Get(key string) (string, error) {

	return "", nil
}

func (f *fscacheServer) Del(key string) (int64, error) {

	return 0, nil
}

func (r *fscacheServer) Debug(val string) string {
	return "[FSCASE]: " + val
}
