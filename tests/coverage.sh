#!/bin/bash

go test -v -covermode=set -coverprofile=./report/coverage.out  ./...  \
    -coverpkg=github.com/NicholeGit/nade/framework/command

echo "========="

go tool cover -func=./report/coverage.out  # 分析 out 文件，得到文字输出

go tool cover -html=./report/coverage.out -o=./report/coverage.html # 分析 out 文件，得到html



# 主要就是使用converpkg参数，把代码覆盖率限制在controller层。
#go test -coverpkg xxx/controllers/... -coverprofile=report/coverage.out ./...
#go tool cover -html=report/coverage.out -o report/coverage.html
#open report/coverage.html
