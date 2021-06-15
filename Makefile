PROJECTNAME=$(shell basename "$(PWD)")
SOL_DIR=./solidity

CENT_EMITTER_ADDR?=0x1
CENT_CHAIN_ID?=0x1
CENT_TO?=0x1234567890
CENT_TOKEN_ID?=0x5
CENT_METADATA?=0x0

.PHONY: help run build install license
all: help

help: Makefile
	@echo
	@echo "Choose a make command to run in "$(PROJECTNAME)":"
	@echo
	@$(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'
	@echo

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

get-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest

.PHONY: lint
lint:
	if [ ! -f ./bin/golangci-lint ]; then \
		$(MAKE) get-lint; \
	fi;
	./bin/golangci-lint run ./... --timeout 5m0s

lint-fix:
	if [ ! -f ./bin/golangci-lint ]; then \
		$(MAKE) get-lint; \
	fi;
	./bin/golangci-lint run ./... --timeout 5m0s --fix

build:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	cd cmd/relay && env GOARCH=amd64 go build -o ../../build/relay
	cd cmd/gencos && env GOARCH=amd64 go build -o ../../build/gencos
	cd cmd/solvault && env GOARCH=amd64 go build -o ../../build/solvault
	cd cmd/soltool && env GOARCH=amd64 go build -o ../../build/soltool

install:
	@echo "  >  \033[32mInstalling rtoken-relay...\033[0m "
	cd cmd/relay && go install
	cd cmd/gencos && go install
	cd cmd/solvault && go install
	cd cmd/soltool && go install

## license: Adds license header to missing files.
license:
	@echo "  >  \033[32mAdding license headers...\033[0m "
	GO111MODULE=off go get -u github.com/google/addlicense
	addlicense -c "Stafi Protocol" -f ./scripts/header.txt -y 2020 .

## Install dependency subkey
install-subkey:
	curl https://getsubstrate.io -sSf | bash -s -- --fast
	cargo install --force --git https://github.com/paritytech/substrate subkey

## Runs go test for all packages except the solidity bindings
test:
	@echo "  >  \033[32mRunning tests...\033[0m "
	go test `go list ./... | grep -v bindings | grep -v e2e`


clean:
	rm -rf build/
