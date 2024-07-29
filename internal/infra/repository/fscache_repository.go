package repository

import (
	"errors"
	"strconv"
	"time"

	fscache "github.com/iqquee/fs-cache"
)

type fscacheServer struct {
	client fscache.Operations
}

func NewFsCahe() *fscacheServer {
	return &fscacheServer{
		client: fscache.New(),
	}
}

func (f *fscacheServer) Set(key string, value any, expiration int32) (string, error) {
	err := f.client.Memdis().Set(key, value, time.Duration(expiration)*time.Second)
	if err != nil {
		return "", err
	}
	val, err := f.client.Memdis().Get(key)
	if err != nil {
		return "", err
	}
	if value != val {
		return "", errors.New("inconsistency to set the key")
	}
	return val.(string), nil
}

func (f *fscacheServer) Increment(key string, expiration time.Duration) (int64, error) {

	val, err := f.client.Memdis().Get(key)
	if err != nil {
		if err.Error() == "key not found" {
			err := f.client.Memdis().Set(key, 1, time.Duration(expiration)*time.Second)
			if err != nil {
				return 0, err
			}
		}
		return 0, err
	}
	valI, err := strconv.Atoi(val.(string))
	if err != nil {
		return 0, errors.New("value is not an integer or out of range")
	}
	/*
		err = f.client.Memdis().OverWrite(key, valI+1, expiration)
		if err != nil {
			return 0, err
		}
	*/
	return int64(valI), nil
}
