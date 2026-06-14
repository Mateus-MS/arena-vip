.PHONY: dev build generate css css-watch install

TAILWIND_VERSION := v3.4.17

ifeq ($(OS),Windows_NT)
TAILWIND_BIN := ./bin/tailwindcss.exe
TAILWIND_URL := https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-windows-x64.exe
else
UNAME_S := $(shell uname -s)
UNAME_M := $(shell uname -m)

ifeq ($(UNAME_S),Darwin)
TAILWIND_OS := macos
else
TAILWIND_OS := linux
endif

ifeq ($(UNAME_M),arm64)
TAILWIND_ARCH := arm64
else ifeq ($(UNAME_M),aarch64)
TAILWIND_ARCH := arm64
else
TAILWIND_ARCH := x64
endif

TAILWIND_BIN := ./bin/tailwindcss
TAILWIND_URL := https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-$(TAILWIND_OS)-$(TAILWIND_ARCH)
endif

install:
	go mod tidy
	mkdir -p ./bin
	curl -sL $(TAILWIND_URL) -o $(TAILWIND_BIN)
	chmod +x $(TAILWIND_BIN)

generate:
	templ generate

css:
	$(TAILWIND_BIN) -i ./static/css/input.css -o ./static/css/output.css --minify

css-watch:
	$(TAILWIND_BIN) -i ./static/css/input.css -o ./static/css/output.css --watch

build: generate css
	go build -o ./bin/arena-vip .

dev: generate css
	go run .
