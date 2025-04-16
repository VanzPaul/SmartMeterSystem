# Simple Makefile for a Go project (Cross-Platform)
ifeq ($(OS),Windows_NT)
    DETECTED_OS := Windows
    TAILWIND_BIN := tailwindcss.exe
    TAILWIND_URL := https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.10/tailwindcss-windows-x64.exe
    BINARY_NAME := main.exe
else
    DETECTED_OS := $(shell uname -s)
    TAILWIND_BIN := tailwindcss
    TAILWIND_URL := https://github.com/tailwindlabs/tailwindcss/releases/download/v3.4.10/tailwindcss-linux-x64
    BINARY_NAME := main
endif

# Build the application
all: build test

templ-install:
ifeq ($(DETECTED_OS),Windows)
	@where templ >nul 2>&1 || ( \
		echo Go's 'templ' is not installed on your machine. & \
		set /p choice="Do you want to install it? [Y/n] " & \
		if /i "!choice!" neq "n" ( \
			go install github.com/a-h/templ/cmd/templ@latest & \
			where templ >nul 2>&1 || ( \
				echo templ installation failed. Exiting... & \
				exit 1 \
			) \
		) else ( \
			echo You chose not to install templ. Exiting... & \
			exit 1 \
		) \
	)
else
	@if ! command -v templ > /dev/null; then \
		read -p "Go's 'templ' is not installed. Install? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/a-h/templ/cmd/templ@latest; \
			command -v templ >/dev/null || { echo "templ installation failed"; exit 1; }; \
		else \
			echo "Installation aborted"; exit 1; \
		fi; \
	fi
endif

tailwind-install:
	@if not exist "$(TAILWIND_BIN)" ( \
		curl -sL "$(TAILWIND_URL)" -o "$(TAILWIND_BIN)" \
	)
ifeq ($(DETECTED_OS),Windows)
	@attrib +x "$(TAILWIND_BIN)" >nul 2>&1
else
	@chmod +x "$(TAILWIND_BIN)"
endif

build: tailwind-install templ-install
	@echo "Building..."
	@templ generate
ifeq ($(DETECTED_OS),Windows)
	@$(TAILWIND_BIN) -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css
else
	@./$(TAILWIND_BIN) -i cmd/web/styles/input.css -o cmd/web/assets/css/output.css
endif
	@go build -o $(BINARY_NAME) cmd/api/main.go

run:
	@go run cmd/api/main.go

docker-run:
ifeq ($(DETECTED_OS),Windows)
	@docker compose up --build || ( \
		echo Falling back to Docker Compose V1 & \
		docker-compose up --build \
	)
else
	@docker compose up --build || ( \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up --build \
	)
endif

docker-down:
ifeq ($(DETECTED_OS),Windows)
	@docker compose down || ( \
		echo Falling back to Docker Compose V1 & \
		docker-compose down \
	)
else
	@docker compose down || ( \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down \
	)
endif

test:
	@echo "Testing..."
	@go test ./... -v

itest:
	@echo "Running integration tests..."
	@go test ./internal/database -v

clean:
	@echo "Cleaning..."
ifeq ($(DETECTED_OS),Windows)
	@del /f $(BINARY_NAME) 2>nul
else
	@rm -f $(BINARY_NAME)
endif

watch:
ifeq ($(DETECTED_OS),Windows)
	@powershell -ExecutionPolicy Bypass -Command "if (Get-Command air -ErrorAction SilentlyContinue) { \
		air; \
		Write-Output 'Watching...'; \
	} else { \
		Write-Output 'Installing air...'; \
		go install github.com/air-verse/air@latest; \
		air; \
		Write-Output 'Watching...'; \
	}"
else
	@if command -v air > /dev/null; then \
		air; \
		echo "Watching..."; \
	else \
		read -p "Go's 'air' is not installed. Install? [Y/n] " choice; \
		if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
			go install github.com/air-verse/air@latest; \
			air; \
			echo "Watching..."; \
		else \
			echo "Installation aborted"; exit 1; \
		fi; \
	fi
endif

.PHONY: all build run test clean watch tailwind-install docker-run docker-down itest templ-install