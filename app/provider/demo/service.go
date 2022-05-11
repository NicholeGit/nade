package demo

import "github.com/NicholeGit/nade/framework"

type Service struct {
	container framework.IContainer
}

func NewService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	return &Service{container: container}, nil
}
