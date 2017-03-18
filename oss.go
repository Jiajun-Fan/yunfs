package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"io"
	"os"
	"strconv"
)

type Oss interface {
	Stat(fname string) error
	Size(fname string) uint64
	Create(name string) (io.WriteCloser, error)
	Open(name string) (io.ReadCloser, error)
}

type LocalOss struct {
	Base string
}

func NewLocalOss(base string) Oss {
	oss := &LocalOss{}
	oss.Base = base
	return oss
}

func (l *LocalOss) Stat(fname string) error {
	fn := fmt.Sprintf("%s/%s", l.Base, fname)
	_, err := os.Stat(fn)
	return err
}

func (l *LocalOss) Size(fname string) uint64 {
	return 0
}

func (l *LocalOss) Create(name string) (io.WriteCloser, error) {
	fn := fmt.Sprintf("%s/%s", l.Base, name)
	return os.Create(fn)
}

func (l *LocalOss) Open(name string) (io.ReadCloser, error) {
	fn := fmt.Sprintf("%s/%s", l.Base, name)
	return os.Open(fn)
}

func MakeOss(config OssConfig) Oss {
	if config.Type == "local" {
		return NewLocalOss(config.EndPoint)
	} else if config.Type == "aliyun" {
		return NewAliyunOss(config)
	} else {
		Fatal(fmt.Sprintf("oss type '%s' is not implemented", config.Type))
		return nil
	}
}

type AliyunOss struct {
	Key        string
	Secret     string
	BucketName string
	EndPoint   string
	client     *oss.Client
	bucket     *oss.Bucket
}

type AliyunFile struct {
	key    string
	bucket *oss.Bucket
	buffer bytes.Buffer
}

func NewAliyunFile(name string, bucket *oss.Bucket) io.WriteCloser {
	f := &AliyunFile{}
	f.bucket = bucket
	f.key = name
	return f
}

func (af *AliyunFile) Write(input []byte) (int, error) {
	af.buffer.Write(input)
	return len(input), nil
}

func (af *AliyunFile) Close() error {
	return af.bucket.PutObject(af.key, &af.buffer)
}

func NewAliyunOss(config OssConfig) Oss {
	aliyun := &AliyunOss{}
	aliyun.Key = config.Key
	aliyun.Secret = config.Secret
	aliyun.BucketName = config.Bucket
	aliyun.EndPoint = config.EndPoint
	if client, err := oss.New(config.EndPoint, config.Key, config.Secret); err != nil {
		panic(err.Error())
	} else {
		aliyun.client = client
	}
	if bucket, err := aliyun.client.Bucket(config.Bucket); err != nil {
		panic(err.Error())
	} else {
		aliyun.bucket = bucket
	}
	return aliyun
}

func (a *AliyunOss) Stat(fname string) error {
	if ok, _ := a.bucket.IsObjectExist(fname); ok {
		return nil
	} else {
		return errors.New("no such file")
	}
}

func (a *AliyunOss) Size(fname string) uint64 {
	if head, err := a.bucket.GetObjectMeta(fname); err == nil {
		v, _ := strconv.Atoi(head["Content-Length"][0])
		return uint64(v)
	} else {
		return 0
	}
}

func (a *AliyunOss) Create(name string) (io.WriteCloser, error) {
	return NewAliyunFile(name, a.bucket), nil
}

func (a *AliyunOss) Open(name string) (io.ReadCloser, error) {
	return a.bucket.GetObject(name)
}
