PROJECTNAME=$(shell basename "$(PWD)")

all: build

get:
	@echo "  >  \033[32mDownloading & Installing all the modules...\033[0m "
	go mod tidy && go mod download

fmt:
	go fmt ./...


get-lint:
	curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s latest

lint:
	golangci-lint run ./... --skip-files ".+_test.go"


build:
	@echo "  >  \033[32mBuilding binary...\033[0m "
	cd cmd/relay && go build -o ../../build/relay
	cd cmd/solvault && go build -o ../../build/solvault
	cd cmd/soltool && go build -o ../../build/soltool

install:
	@echo "  >  \033[32mInstalling rtoken-relay...\033[0m "
	cd cmd/relay && go install
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

abi:
	@echo " > \033[32mGenabi...\033[0m "
	abigen --abi ./bindings/StakeERC20Portal/stakeportal_abi.json --pkg stake_erc20_portal --type StakeERC20Portal --out ./bindings/StakeERC20Portal/StakeERC20Portal.go
	abigen --abi ./bindings/StakeNativePortal/stakenativeportal_abi.json --pkg stake_native_portal --type StakeNativePortal --out ./bindings/StakeNativePortal/StakeNativePortal.go
	abigen --abi ./bindings/MultisigOnchain/multisigonchain_abi.json --pkg multisig_onchain --type MultisigOnchain --out ./bindings/MultisigOnchain/MultisigOnchain.go
	abigen --abi ./bindings/Staking/staking_abi.json --pkg staking --type Staking --out ./bindings/Staking/Staking.go

## Runs go test for all packages except the solidity bindings
test:
	@echo "  >  \033[32mRunning tests...\033[0m "
	go test `go list ./... | grep -v bindings | grep -v e2e`


clean:
	rm -rf build/


.PHONY: help run build install license lint lint-fix get-lint abi