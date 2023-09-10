.SILENT:

BIN_NAME := $(shell grep -o 'module .*' go.mod | awk '{print $$2}')

build:
	go build .

install: build
	mv ./$(BIN_NAME) $$GOPATH/bin/$(BIN_NAME)

get_abs_path_bin:
	realpath $$GOPATH/bin/$(BIN_NAME)