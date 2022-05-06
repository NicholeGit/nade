package framework

import (
	"fmt"
	"sync"

	"github.com/pkg/errors"
)

type IDepend interface {
	DependOn() []string
}

var providers = make(map[string]IServiceProvider)

func Register(p IServiceProvider) error {
	if _, ok := providers[p.Name()]; ok {
		return errors.New(p.Name() + " is register repeat !")
	}
	providers[p.Name()] = p
	return nil
}

// IContainer 是一个服务容器，提供绑定服务和获取服务的功能
type IContainer interface {
	// BindAll 绑定所有 Provider
	BindAll() error

	// IsBind 关键字凭证是否已经绑定服务提供者
	IsBind(key string) bool

	// Make 根据关键字凭证获取一个服务，
	Make(key string) (interface{}, error)
	// MustMake 根据关键字凭证获取一个服务，如果这个关键字凭证未绑定服务提供者，那么会panic。
	// 所以在使用这个接口的时候请保证服务容器已经为这个关键字凭证绑定了服务提供者。
	MustMake(key string) interface{}
	// MakeNew 根据关键字凭证获取一个服务，只是这个服务并不是单例模式的
	// 它是根据服务提供者注册的启动函数和传递的params参数实例化出来的
	// 这个函数在需要为不同参数启动不同实例的时候非常有用
	MakeNew(key string, params []interface{}) (interface{}, error)
}

// 编译检查是否实现所有接口
var _ IContainer = &NadeContainer{}

// NadeContainer 是服务容器的具体实现
type NadeContainer struct {
	// providers 存储注册的服务提供者，key为字符串凭证
	providers map[string]IServiceProvider
	// instance 存储具体的实例，key为字符串凭证
	instances map[string]interface{}
	// lock 用于锁住对容器的变更操作
	lock sync.RWMutex

	// 处理依赖缓存
	providerChan chan IServiceProvider
}

// NewNadeContainer 创建一个服务容器
func NewNadeContainer() *NadeContainer {
	return &NadeContainer{
		providers:    map[string]IServiceProvider{},
		instances:    map[string]interface{}{},
		lock:         sync.RWMutex{},
		providerChan: make(chan IServiceProvider, len(providers)),
	}
}

// PrintProviders 输出服务容器中注册的关键字
func (n *NadeContainer) PrintProviders() []string {
	var ret []string
	for _, provider := range n.providers {
		name := provider.Name()
		line := fmt.Sprint(name)
		ret = append(ret, line)
	}
	return ret
}

func (n *NadeContainer) BindAll() error {
	// 所有的插件
	for _, p := range providers {
		n.providerChan <- p
	}

	num := len(n.providerChan)
	for num > 0 {
		p := <-n.providerChan
		canBind := true
		if deps, ok := p.(IDepend); ok {
			depends := deps.DependOn()
			for _, dependName := range depends {
				if p := n.findServiceProvider(dependName); p == nil {
					// 有未加载的前置插件
					canBind = false
					break
				}
			}
		}
		if canBind {
			err := n.bind(p)
			if err != nil {
				return errors.Wrap(err, p.Name()+" bind is error!")
			}
			num--
		} else {
			// 如果插件不能被 bind, 这个插件就塞入到最后一个队列
			n.providerChan <- p
		}
	}
	return nil
}

// Bind 将服务容器和关键字做了绑定
func (n *NadeContainer) bind(provider IServiceProvider) error {
	fmt.Println("Bind", provider.Name())
	n.lock.Lock()
	key := provider.Name()
	n.providers[key] = provider
	n.lock.Unlock()

	// if provider is not defer
	if provider.IsDefer() == false {
		if err := provider.Boot(n); err != nil {
			return err
		}
		// 实例化方法
		params := provider.Params(n)
		method := provider.Register(n)
		instance, err := method(params...)
		if err != nil {
			fmt.Println("bind service provider ", key, " error: ", err)
			return errors.New(err.Error())
		}
		n.instances[key] = instance
	}
	return nil
}

func (n *NadeContainer) IsBind(key string) bool {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.findServiceProvider(key) != nil
}

func (n *NadeContainer) findServiceProvider(key string) IServiceProvider {
	if sp, ok := n.providers[key]; ok {
		return sp
	}
	return nil
}

func (n *NadeContainer) Make(key string) (interface{}, error) {
	return n.make(key, nil, false)
}

func (n *NadeContainer) MustMake(key string) interface{} {
	serv, err := n.make(key, nil, false)
	if err != nil {
		panic("container not contain key " + key)
	}
	return serv
}

func (n *NadeContainer) MakeNew(key string, params []interface{}) (interface{}, error) {
	return n.make(key, params, true)
}

func (n *NadeContainer) newInstance(sp IServiceProvider, params []interface{}) (interface{}, error) {
	// force new a
	if err := sp.Boot(n); err != nil {
		return nil, err
	}
	if params == nil {
		params = sp.Params(n)
	}
	method := sp.Register(n)
	ins, err := method(params...)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return ins, err
}

// 真正的实例化一个服务
func (n *NadeContainer) make(key string, params []interface{}, forceNew bool) (interface{}, error) {
	n.lock.RLock()
	defer n.lock.RUnlock()
	// 查询是否已经注册了这个服务提供者，如果没有注册，则返回错误
	sp := n.findServiceProvider(key)
	if sp == nil {
		return nil, errors.New("contract " + key + " have not register")
	}

	if forceNew {
		return n.newInstance(sp, params)
	}

	// 不需要强制重新实例化，如果容器中已经实例化了，那么就直接使用容器中的实例
	if ins, ok := n.instances[key]; ok {
		return ins, nil
	}

	// 容器中还未实例化，则进行一次实例化
	inst, err := n.newInstance(sp, nil)
	if err != nil {
		return nil, err
	}
	n.instances[key] = inst
	return inst, nil
}

// NameList 列出容器中所有服务提供者的字符串凭证
func (n *NadeContainer) NameList() []string {
	var ret []string
	for _, provider := range n.providers {
		name := provider.Name()
		ret = append(ret, name)
	}
	return ret
}
