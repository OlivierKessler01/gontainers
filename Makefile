.PHONY:
setup_hooks:
	git config core.hooksPath .githooks

.PHONY:
.ONESHELL:
build:
	@rm main || true
	@go build  main.go 
	./main list ls -a

.PHONY:
debug: 
	~/go/bin/dlv debug ./main.go -- list ls -a

