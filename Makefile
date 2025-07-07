.PHONY:
setup_hooks:
	git config core.hooksPath .githooks

.PHONY:
.ONESHELL:
build:
	@rm main || true
	@go build  main.go 
	./main run ls -a

.PHONY:
debug: 
	~/go/bin/dlv debug ./main.go -- run ls -a

