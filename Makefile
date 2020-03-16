GO ?= go

GREEN_COLOR = \x1b[32;01m
END_COLOR = \x1b[0m

all: test

.PHONY: test
test:
	@$(GO) test -v

.PHONY: coverage
coverage:
	@$(GO) test -race -coverprofile=coverage.txt -covermode=atomic

.PHONY: run
run:
	@$(GO) run ./_example/basic/main.go

.PHONY: race_detect
race_detect:
	@$(GO) run -race ./_example/basic/main.go

.PHONY: install
install:
	@echo "$(GREEN_COLOR)Installing dependencies...$(END_COLOR)"
	@$(GO) mod download
	@$(GO) mod verify

.PHONY: clean
clean:
	@$(GO) clean -x -i ./...
