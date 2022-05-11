## This is a self-documented Makefile. For usage information, run `make help`:
##
## For more information, refer to https://www.thapaliya.com/en/writings/well-documented-makefiles/

.DEFAULT_GOAL	:= help

.PHONY: start
start: ## run app start
	go run main.go app start

.PHONY: cron
cron: ## run cron start
	go run main.go cron start

.PHONY: dev
dev: ## run dev
	go run main.go dev backend

##@ Helpers
.PHONY: help
help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; \
	printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { \
	printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { \
	printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' \
	$(MAKEFILE_LIST)
