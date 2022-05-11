package kernel

import (
	"github.com/NicholeGit/nade/framework"
	"github.com/NicholeGit/nade/framework/contract"
	"github.com/NicholeGit/nade/framework/util"
	"github.com/pkg/errors"
)

// 引擎服务
type NadeKernelService struct {
	container framework.IContainer // 服务容器
	servers   map[string]contract.IServer
}

// NewNadeKernelService 初始化引擎服务实例
func NewNadeKernelService(params ...interface{}) (interface{}, error) {
	container := params[0].(framework.IContainer)
	return &NadeKernelService{container: container, servers: make(map[string]contract.IServer, 0)}, nil
}

func (s *NadeKernelService) Servers() []contract.IServer {
	l := make([]contract.IServer, 0, len(s.servers))
	for k := range s.servers {
		l = append(l, s.servers[k])
	}
	return l
}

func (s *NadeKernelService) AddServers(srv ...contract.IServer) error {
	for _, v := range srv {
		if _, ok := s.servers[v.Name()]; ok {
			return errors.New("Duplicate server registration")
		}
		s.servers[v.Name()] = v
	}
	return nil
}

func (s *NadeKernelService) GetServer(name string) contract.IServer {
	return s.servers[name]
}

func (s *NadeKernelService) Info() string {
	ps := [][]string{}
	for k, v := range s.servers {
		line := []string{k, v.Addr()}
		ps = append(ps, line)
	}
	return util.Pretty(ps)
}
