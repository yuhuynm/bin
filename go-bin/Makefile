.PHONY: run build tidy check-air install-air

APP_NAME := go-bin
MAIN := ./cmd/server
BIN := ./tmp/$(APP_NAME)
GOBIN := $(shell go env GOBIN)
GOPATH := $(shell go env GOPATH)
AIR := $(shell command -v air 2>/dev/null)

ifeq ($(strip $(AIR)),)
ifeq ($(strip $(GOBIN)),)
AIR := $(GOPATH)/bin/air
else
AIR := $(GOBIN)/air
endif
endif

run: check-air
	$(AIR) -c .air.toml

build:
	go build -o $(BIN) $(MAIN)

tidy:
	go mod tidy

check-air:
	@[ -x "$(AIR)" ] || (echo "air is not installed. Run: make install-air"; exit 1)

install-air:
	go install github.com/air-verse/air@latest
