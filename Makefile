CONFIG_PATH := $(PWD)/config/local.test.yml
export CONFIG_PATH
.PHONY: test
test:
	@echo "Using CONFIG_PATH=$(CONFIG_PATH)"
	@go test -v ./... 
