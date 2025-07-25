export GOPATH=/home/olivierkessler/go
export PATH=$(GOPATH)/bin:$(shell echo $$PATH)
GOFLAGS = 

help: 
	@echo " "
	@grep '^[^.#]\+:\s\+.*#' Makefile | \
	sed "s/\(.\+\):\s*\(.*\) #\s*\(Usage\s*\`.*\`\)/`printf "\033[93m"`  \1`printf "\033[0m"`	`printf "\033[31m"` \3`printf "\033[0m"` [\2]/" | \
	expand -45
	@echo " "

.PHONY:
setup_hooks:
	git config core.hooksPath .githooks

.PHONY:
.ONESHELL:
build: # Usage `make build` Compile gontainers 
	@rm gontainers || true
	@go build -o gontainers main.go

.PHONY:
debug_run: build # Usage `make debug_run` Step debug the run command
	@sudo env GOPATH=$(GOPATH) PATH=$(PATH) dlv debug ./main.go -- run tail -f /dev/null --verbose

.PHONY:
debug_list: build # Usage `make debug_run` Step debug the list command
	~/go/bin/dlv debug ./main.go -- list 

.PHONY:
install_cri_tools:
	VERSION="v1.30.0"
	curl -LO https://github.com/kubernetes-sigs/cri-tools/releases/download/${VERSION}/crictl-${VERSION}-linux-amd64.tar.gz
	sudo tar -C /usr/local/bin -xzf crictl-${VERSION}-linux-amd64.tar.gz

.PHONY:
test: #Usage `make test` Run the tests
	@cd process && sudo env GOPATH=$(GOPATH) PATH=$(PATH) go test


