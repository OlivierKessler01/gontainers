export GOPATH=/home/olivierkessler/go
export PATH=$(GOPATH)/bin:$(shell echo $$PATH)

.PHONY:
setup_hooks:
	git config core.hooksPath .githooks

.PHONY:
.ONESHELL:
build:
	@rm gontainers || true
	@go build  -o gontainers main.go

.PHONY:
debug_run: build
	sudo env GOPATH=$(GOPATH) PATH=$(PATH) dlv debug ./main.go -- run tail -f /dev/null

.PHONY:
debug_list: build
	~/go/bin/dlv debug ./main.go -- list 


