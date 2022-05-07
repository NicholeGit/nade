package id

import (
	// Package xid是一个全球唯一的id生成器库，可以安全地直接在服务器代码中使用。
	// Xid使用Mongo对象ID算法生成具有不同序列化（base64）的全局唯一ID，以使其在作为字符串传输时更短
	// https://docs.mongodb.org/manual/reference/object-id/
	"github.com/rs/xid"
)

type NadeIDService struct {
}

func NewNadeIDService(_ ...interface{}) (interface{}, error) {
	return &NadeIDService{}, nil
}

func (s *NadeIDService) NewID() string {
	return xid.New().String()
}
