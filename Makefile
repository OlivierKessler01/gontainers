export GOPATH=/home/olivierkessler/go
export PATH=$(GOPATH)/bin:$(shell echo $$PATH)
GOFLAGS = 
RED    := \033[0;31m
GREEN  := \033[0;32m
YELLOW := \033[0;33m
NC     := \033[0m


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
	@sudo env GOPATH=$(GOPATH) PATH=$(PATH) dlv debug ./main.go -- run --verbose --command="tail -f /dev/null"

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
	@sudo env GOPATH=$(GOPATH) PATH=$(PATH) go test



################################## K8S node/virtualization ###############################3
.PHONY: 
k8s-build-image: #Usage `make k8s-build-vm-image` Build a k8s-compatible VM image from a ubuntu base image.
	# 2. Create cloud-init files
	rm k8s-virt/cloud-init.iso || true #Delete VM image
	rm k8s-virt/state.img || true #Delete VM disk storage to start anew
	cp k8s-virt/ubuntu.img k8s-virt/state.img
	echo 'instance-id: gontainers' > meta-data
	@echo '#cloud-config' > user-data
	@echo 'users:' >> user-data
	@echo '  - name: user' >> user-data
	@echo '    ssh-authorized-keys:' >> user-data
	@echo '      - $(shell cat ~/.ssh/id_rsa.pub)' >> user-data
	@echo '    sudo: ALL=(ALL) NOPASSWD:ALL' >> user-data
	@echo '    shell: /bin/bash' >> user-data
	@echo '' >> user-data
	@echo 'packages:' >> user-data
	@echo '  - apt-transport-https' >> user-data
	@echo '  - ca-certificates' >> user-data
	@echo '  - curl' >> user-data
	@echo '' >> user-data
	@echo 'write_files:' >> user-data
	@echo '  - path: /etc/apt/sources.list.d/kubernetes.list' >> user-data
	@echo '    content: |' >> user-data
	@echo '      deb https://apt.kubernetes.io/ kubernetes-xenial main' >> user-data
	@echo 'runcmd:' >> user-data
	@echo '  - apt-get update' >> user-data
	@echo '  - sudo apt-get install -y apt-transport-https ca-certificates curl gpg' >> user-data
	@echo '  - curl -fsSL https://pkgs.k8s.io/core:/stable:/v1.33/deb/Release.key | sudo gpg --dearmor -o /etc/apt/keyrings/kubernetes-apt-keyring.gpg' >> user-data
	@echo '  - echo 'deb [signed-by=/etc/apt/keyrings/kubernetes-apt-keyring.gpg] https://pkgs.k8s.io/core:/stable:/v1.33/deb/ /' | sudo tee /etc/apt/sources.list.d/kubernetes.list' >> user-data
	@echo '  - apt-get update' >> user-data
	@echo '  - apt-get install -y kubelet cri-tools' >> user-data
	@echo '  - mkdir -p /etc/kubelet/manifests' >> user-data
	# 3. Build cloud-init ISO
	#sudo dnf install cloud-utils
	cloud-localds k8s-virt/cloud-init.iso user-data meta-data

.PHONY:
k8s-run-vm: #Usage `make k8s-test-setup` Launch a KVM VM and launch kubelet to make it a k8s-node.
	rm k8s-virt/state.img || true #Delete VM disk storage to start anew
	cp k8s-virt/ubuntu.img k8s-virt/state.img
	qemu-img resize k8s-virt/state.img 20G
	qemu-system-x86_64 \
	  -m 4096 -smp 2 \
	  -nographic \
	  -enable-kvm \
	  -drive file=k8s-virt/state.img,format=qcow2,if=virtio \
	  -drive file=k8s-virt/cloud-init.iso,format=raw,if=virtio \
	  -netdev user,id=n1,hostfwd=tcp::2222-:22 \
	  -device virtio-net-pci,netdev=n1

k8s-test: #Usage `make k8s-test` Run the ./gontainers grpc server into the VM. Then run kubelet.
	@echo -e "$(GREEN)Uploading config and running the grpc server$(NC)"
	scp -P 2222 ./k8s-virt/kubelet_config.yml user@localhost:/home/user/
	ssh -p 2222 user@localhost sudo killall gontainers
	scp -P 2222 ./gontainers user@localhost:/home/user/
	ssh -p 2222 user@localhost "chmod +x ./gontainers"
	ssh -p 2222 user@localhost "sudo ./gontainers server -v" & true
	sleep 1
	@echo -e "$(GREEN)Testing the grpc server by calling the Version endpoint$(NC)"
	ssh -p 2222 user@localhost "sudo crictl --runtime-endpoint unix:///var/run/gontainers.sock version"
	@echo -e "$(GREEN)Running kubelet$(NC)"
	ssh -p 2222 user@localhost "sudo kubelet --config=kubelet_config.yml"
