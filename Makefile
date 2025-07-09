.PHONY:
setup_hooks:
	git config core.hooksPath .githooks

.PHONY:
.ONESHELL:
build:
	@rm gontainers || true
	@go build  -o gontainers main.go

.PHONY:
debug: 
	~/go/bin/dlv debug ./main.go -- run ls -a

