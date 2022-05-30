output = "./build"
package = "git-hook-bot"

.PHONY: help

## help: Show all make commands.
help: Makefile
	@echo "Choose a command run in "$(package)":"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'

## install: Install dependencies.
install:
	@echo " >  installing dependencies..."
	@go mod tidy

## clean: Clean build cache.
clean:
	@echo " >  Cleaning build cache..."
	@if [ -e "${output}" ]; then rm -rf "${output}" ; fi

## build: Build binary.
build: install clean
	@echo " >  Building binary..."
	@go build -o "${output}/${package}" cmd/main.go

## run: Execute build file
run: build
	@echo " >  Execute binary file..."
	@"${output}/${package}"
