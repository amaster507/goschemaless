.PHONY: help init init-service fmt tidy build build-all run test watch watch-run clean

help:
	@echo "Go Monorepo Makefile"
	@echo ""
	@echo "Usage:"
	@echo "  make init <module>              Initialize root module (e.g., github.com/user/project)"
	@echo "  make init-service <name>        Add a new service to the workspace"
	@echo "  make build <service>            Build a specific service"
	@echo "  make build-all                  Build all services"
	@echo "  make run <service>              Run a specific service"
	@echo "  make test <service>             Run tests for a service"
	@echo "  make watch <service>            Watch and rebuild a service on changes"
	@echo "  make watch-run <service>        Watch and run a service on changes"
	@echo "  make fmt                        Format all Go code"
	@echo "  make tidy                       Tidy dependencies"
	@echo "  make clean                      Remove bin/ directory"

init:
	@MODULE="$(filter-out init,$(MAKECMDGOALS))"; \
	if [ -z "$$MODULE" ]; then \
		echo "Usage: make init github.com/user/project"; \
		exit 1; \
	fi; \
	go mod init $$MODULE; \
	mkdir -p cmd internal pkg; \
	go work init ./; \
	echo "Project initialized at $$MODULE"
	@:

init-service:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make init-service <name>"; \
		exit 1; \
	fi; \
	if [ ! -f go.mod ]; then \
		echo "Error: go.mod not found. Run 'make init <module>' first."; \
		exit 1; \
	fi; \
	ROOT_MODULE=$$(grep '^module ' go.mod | awk '{print $$2}'); \
	MODULE="$$ROOT_MODULE/$$SERVICE"; \
	mkdir -p cmd/$$SERVICE; \
	(cd cmd/$$SERVICE && go mod init $$MODULE); \
	go work use ./cmd/$$SERVICE; \
	echo "package main\n\nfunc main() {\n\t// TODO: $$SERVICE\n}" > cmd/$$SERVICE/main.go; \
	echo "Service '$$SERVICE' created at $$MODULE"
	@:

fmt:
	go fmt ./...

tidy:
	go mod tidy -v

build:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make build <service>"; \
		exit 1; \
	fi; \
	go build -o bin/$$SERVICE ./cmd/$$SERVICE
	@:

build-all:
	@for dir in $$(ls -d cmd/*/ 2>/dev/null | sed 's|cmd/||g' | sed 's|/||g'); do \
		echo "Building $$dir..."; \
		go build -o bin/$$dir ./cmd/$$dir || exit 1; \
	done
	@echo "All services built"

run:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make run <service>"; \
		exit 1; \
	fi; \
	go run ./cmd/$$SERVICE
	@:

test:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make test <service>"; \
		exit 1; \
	fi; \
	if [ -d internal ]; then \
		go test ./internal/...; \
	fi; \
	go test ./cmd/$$SERVICE/...
	@:

clean:
	rm -rf bin/

watch:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make watch <service>"; \
		exit 1; \
	fi; \
	echo "Watching $$SERVICE for changes..."; \
	go run github.com/cosmtrek/air@latest -- -c <(echo 'root = "cmd/'"$$SERVICE"'"\nfull_bin = "bin/'"$$SERVICE"'"') 2>/dev/null || \
	while true; do \
		inotifywait -r -e modify cmd/$$SERVICE 2>/dev/null && make build SERVICE=$$SERVICE; \
	done
	@:

watch-run:
	@SERVICE="$(word 2,$(MAKECMDGOALS))"; \
	if [ -z "$$SERVICE" ]; then \
		echo "Usage: make watch-run <service>"; \
		exit 1; \
	fi; \
	echo "Watching and running $$SERVICE..."; \
	while true; do \
		go run ./cmd/$$SERVICE; \
		inotifywait -r -e modify cmd/$$SERVICE internal pkg go.* 2>/dev/null && echo "Restarting..."; \
	done
	@:

%:
	@:

