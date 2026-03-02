# MVTable Server Makefile

# Go 构建变量
BINARY_NAME=mvtable
MAIN_FILE=.
GO=go
GOFLAGS=-v
LDFLAGS=-ldflags "-s -w"

# 构建模式
BUILD_MODE=release
GOOS=$(shell $(GO) env GOOS)
GOARCH=$(shell $(GO) env GOARCH)

# 环境变量
ENV ?= dev
CONFIG_FILE=configs/config.$(ENV).yaml

# 颜色定义
RED=\033[0;31m
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

.PHONY: help build run dev test clean deps install init generate swag migrate docker-build docker-run docker-stop

# 默认目标：显示帮助
help:
	@echo ""
	@echo "$(GREEN)MVTable Server Makefile$(NC)"
	@echo ""
	@echo "用法: make <目标> [ENV=<环境>]"
	@echo ""
	@echo "可用目标:"
	@echo "  $(GREEN)help$(NC)       - 显示帮助信息"
	@echo "  $(GREEN)build$(NC)      - 编译应用程序"
	@echo "  $(GREEN)run$(NC)        - 运行应用程序"
	@echo "  $(GREEN)dev$(NC)        - 开发模式运行（热重载）"
	@echo "  $(GREEN)test$(NC)       - 运行测试"
	@echo "  $(GREEN)clean$(NC)       - 清理构建产物"
	@echo "  $(GREEN)deps$(NC)       - 下载依赖"
	@echo "  $(GREEN)install$(NC)    - 安装依赖"
	@echo "  $(GREEN)init$(NC)       - 初始化项目（deps + wire + swag）"
	@echo "  $(GREEN)generate$(NC)   - 生成代码（wire, swag）"
	@echo "  $(GREEN)swag$(NC)       - 生成 Swagger API 文档"
	@echo "  $(GREEN)wire$(NC)       - 生成依赖注入代码"
	@echo "  $(GREEN)migrate$(NC)    - 运行数据库迁移"
	@echo "  $(GREEN)docker-build$(NC) - 构建 Docker 镜像"
	@echo "  $(GREEN)docker-run$(NC)   - 运行 Docker 容器"
	@echo "  $(GREEN)docker-stop$(NC)  - 停止 Docker 容器"
	@echo ""
	@echo "示例:"
	@echo "  make build             # 为当前平台编译"
	@echo "  make run ENV=dev      # 使用 dev 配置运行"
	@echo "  make run ENV=prod     # 使用 prod 配置运行"
	@echo "  make dev ENV=dev      # 开发模式热重载"
	@echo "  make docker-build     # 构建 Docker 镜像"
	@echo ""

# 编译应用程序
build:
	@echo "$(YELLOW)正在编译 MVTable server...$(NC)"
	@mkdir -p bin
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME) $(MAIN_FILE)
	@echo "$(GREEN)编译完成! 二进制文件: bin/$(BINARY_NAME)$(NC)"

# 为不同平台编译
build-linux:
	@echo "$(YELLOW)为 Linux 平台编译...$(NC)"
	GOOS=linux GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-linux-amd64 $(MAIN_FILE)

build-darwin:
	@echo "$(YELLOW)为 macOS 平台编译...$(NC)"
	GOOS=darwin GOARCH=amd64 $(GO) build $(GOFLAGS) $(LDFLAGS) -o bin/$(BINARY_NAME)-darwin-amd64 $(MAIN_FILE)

# 运行应用程序
run:
	@echo "$(YELLOW)正在运行 MVTable server (ENV=$(ENV))...$(NC)"
	$(GO) run $(MAIN_FILE) server -c $(CONFIG_FILE) -e $(ENV)

# 开发模式运行（热重载）
dev:
	@echo "$(YELLOW)启动开发模式（热重载）...$(NC)"
	@which air > /dev/null 2>&1 || $(GO) install github.com/air-verse/air@latest
	@cd $(shell pwd) && air

# 清理构建产物
clean:
	@echo "$(YELLOW)清理构建产物...$(NC)"
	rm -rf bin/
	rm -rf tmp/
	rm -f coverage.out coverage.html
	@echo "$(GREEN)清理完成!$(NC)"

# 下载依赖
deps:
	@echo "$(YELLOW)下载依赖...$(NC)"
	$(GO) mod download
	$(GO) mod tidy
	@echo "$(GREEN)依赖下载完成!$(NC)"

# 安装依赖和工具
install:
	@echo "$(YELLOW)安装项目依赖...$(NC)"
	$(GO) mod download
	@echo "$(GREEN)安装开发工具...$(NC)"
	$(GO) install github.com/air-verse/air@latest
	$(GO) install github.com/swaggo/swag/cmd/swag@latest
	$(GO) install github.com/google/wire/cmd/wire@latest
	@echo "$(GREEN)安装完成!$(NC)"

# 初始化项目（依赖 + wire + swag）
init: deps generate swag
	@echo "$(GREEN)项目初始化完成!$(NC)"

# 生成代码（wire, swag）
generate: wire swag

# 生成 Swagger API 文档
swag:
	@echo "$(YELLOW)生成 Swagger 文档...$(NC)"
	@which swag > /dev/null 2>&1 || $(GO) install github.com/swaggo/swag/cmd/swag@latest
	swag init -g main.go -o docs
	@echo "$(GREEN)Swagger 文档生成完成!$(NC)"

# 生成依赖注入代码
wire:
	@echo "$(YELLOW)生成依赖注入代码...$(NC)"
	@which wire > /dev/null 2>&1 || $(GO) install github.com/google/wire/cmd/wire@latest
	cd internal/di && wire
	@echo "$(GREEN)Wire 代码生成完成!$(NC)"

# 运行数据库迁移
migrate:
	@echo "$(YELLOW)运行数据库迁移...$(NC)"
	@if [ -f "scripts/create.sql" ]; then \
		echo "请手动运行迁移: psql -h <host> -U <user> -d <database> -f scripts/create.sql"; \
	else \
		echo "迁移脚本未找到!"; \
	fi

# 显示版本
version:
	$(GO) run $(MAIN_FILE) version

# 运行代码检查（需要 golangci-lint）
lint:
	@echo "$(YELLOW)运行代码检查...$(NC)"
	@which golangci-lint > /dev/null 2>&1 || $(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	golangci-lint run

# 格式化代码
fmt:
	@echo "$(YELLOW)格式化代码...$(NC)"
	$(GO) fmt ./...
	gofmt -w -s .
	@echo "$(GREEN)代码格式化完成!$(NC)"
