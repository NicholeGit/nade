## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://www.thapaliya.com/en/writings/well-documented-makefiles/

GO_RUN		= go run -race main.go

SRC=$(shell find . -name "*.go")

ifeq (, $(shell which richgo))
$(warning "could not find richgo in $(PATH), run: go get github.com/kyoh86/richgo")
endif

ifeq (, $(shell which golangci-lint))
$(warning "could not find golangci-lint in $(PATH), run: curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh")
endif

.DEFAULT_GOAL	:= help

.PHONY: start
start: ## run app start
	$(GO_RUN) app start

.PHONY: cron
cron: ## run cron start
	$(GO_RUN) cron start

.PHONY: dev
dev: ## run dev
	$(GO_RUN) dev backend

.PHONY: fmt
fmt: ## gofmt
	$(info ******************** checking formatting ********************)
	@test -z $(shell gofmt -l $(SRC)) || (gofmt -d $(SRC); exit 1)

.PHONY: test
test: lint  ## test
	$(info ******************** running tests ********************)
	richgo test -v ./...  -covermode=atomic
	# go test  -covermode=atomic  ./...  覆盖率

lint:  ## golangci-lint
	$(info ******************** running lint tools ********************)
	golangci-lint run -v



##@ Helpers
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; \
	printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { \
	printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { \
	printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
	$(MAKEFILE_LIST)
