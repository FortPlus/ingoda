GO=go
PROJECT=autobot
PROJECT_MAIN=$(PROJECT).go

all:	build test

run:
		$(GO) run $(PROJECT_MAIN)

build: get-dependencies build-default

test:
	$(GO) test

get-dependencies:
	$(GO) mod tidy

build-default:
	$(GO) build  -o ../../bin/$(PROJECT)$(GOEXE) $(PROJECT_MAIN)


build-windows: export GOOS=windows
build-windows: export GOARCH=amd64
build-windows:
	$(GO) build  -o ../../bin/$(PROJECT)-win-amd64.exe $(PROJECT_MAIN)

build-linux: export GOOS=linux
build-linux: export GOARCH=arm
build-linux:
	$(GO) build  -o ../../bin/$(PROJECT)-win-amd64 $(PROJECT_MAIN)


