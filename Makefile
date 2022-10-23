## Makefile

.PHONY: setup
setup: ## Install all the build and lint dependencies
	go mod download

.PHONY: fmt
fmt:
	find . -name '*.go' -not -wholename './vendor/*' | while read -r file; do goimports -w "$$file"; done

.PHONY: check
check: ## Run all the linters
	staticcheck ./...	

.PHONY: build
build: cellserv_global cellserv_account cellserv_root cellserv_data cellserv_log cellserv_view cellserv_center cellserv_game cellserv_logic cellserv_gate

## 各进程编译目标
.PHONY: cellserv_global
cellserv_global:
	go build -o ./bin/cellserv_global ./cmd/cellserv_global/main.go

.PHONY: cellserv_account
cellserv_account:
	go build -o ./bin/cellserv_account ./cmd/account/main.go

.PHONY: cellserv_root
cellserv_root:
	go build -o ./bin/cellserv_root ./cmd/cellserv_root/main.go

.PHONY: cellserv_data
cellserv_data:
	go build -o ./bin/cellserv_data ./cmd/cellserv_data/main.go


.PHONY: cellserv_log
cellserv_log:
	go build -o ./bin/cellserv_log ./cmd/cellserv_log/main.go

.PHONY: cellserv_view
cellserv_view:
	go build -o ./bin/cellserv_view ./cmd/cellserv_view/main.go

.PHONY: cellserv_center
cellserv_center:
	go build -o ./bin/cellserv_center ./cmd/cellserv_center/main.go

.PHONY: cellserv_game
cellserv_game:
	go build -o ./bin/cellserv_game ./cmd/cellserv_game/main.go

.PHONY: cellserv_logic
cellserv_logic:
	go build -o ./bin/cellserv_logic ./cmd/cellserv_logic/main.go

.PHONY: cellserv_gate
cellserv_gate:
	go build -o ./bin/cellserv_gate ./cmd/cellserv_gate/main.go

.PHONY:clean

clean: ## Remove temporary files
	go clean
	rm -rf bin

.PHONY:code
code:
	go test cmd/generate/*

# Absolutely awesome: http://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := build
