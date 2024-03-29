package env

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"path"
	"strings"

	"github.com/NicholeGit/nade/framework/contract"
	"github.com/pkg/errors"
)

const AppEnv = "APP_ENV"

var _ contract.IEnv = &NadeEnv{}

// NadeEnv 是 Env 的具体实现
type NadeEnv struct {
	folder string            // 代表.env所在的目录
	maps   map[string]string // 保存所有的环境变量
}

// NewNadeEnv 有一个参数，.env文件所在的目录
//  example: NewNadeEnv("/envfolder/") 会读取文件: /envfolder/.env
//  .env的文件格式 FOO_ENV=BAR
func NewNadeEnv(params ...interface{}) (interface{}, error) {
	if len(params) != 1 {
		return nil, errors.New("NewHadeEnv param error")
	}

	// 读取folder文件
	folder := params[0].(string)

	// 实例化
	nadeEnv := &NadeEnv{
		folder: folder,
		// 实例化环境变量，APP_ENV默认设置为开发环境
		maps: map[string]string{AppEnv: contract.EnvDevelopment},
	}

	// 解析folder/.env文件
	file := path.Join(folder, ".env")
	// 读取.env文件, 不管任意失败，都不影响后续

	// 打开文件.env
	fi, err := os.Open(file)
	if err == nil {
		defer fi.Close()

		// 读取文件
		br := bufio.NewReader(fi)
		for {
			// 按照行进行读取
			line, _, c := br.ReadLine()
			if c == io.EOF {
				break
			}
			// 按照等号解析
			s := bytes.SplitN(line, []byte{'='}, 2)
			// 如果不符合规范，则过滤
			const param = 2
			if len(s) < param {
				continue
			}
			// 保存map
			key := string(s[0])
			val := string(s[1])
			nadeEnv.maps[key] = val
		}
	}

	// 获取当前程序的环境变量，并且覆盖.env文件下的变量
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		const param = 2
		if len(pair) < param {
			continue
		}
		nadeEnv.maps[pair[0]] = pair[1]
	}

	// 返回实例
	return nadeEnv, nil
}

// AppEnv 获取表示当前APP环境的变量APP_ENV
func (en *NadeEnv) AppEnv() string {
	return en.Get(AppEnv)
}

// IsExist 判断一个环境变量是否有被设置
func (en *NadeEnv) IsExist(key string) bool {
	_, ok := en.maps[key]
	return ok
}

// Get 获取某个环境变量，如果没有设置，返回""
func (en *NadeEnv) Get(key string) string {
	if val, ok := en.maps[key]; ok {
		return val
	}
	return ""
}

// All 获取所有的环境变量，.env和运行环境变量融合后结果
func (en *NadeEnv) All() map[string]string {
	return en.maps
}
