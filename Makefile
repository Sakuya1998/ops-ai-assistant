# 运维AI助手 - Makefile

.PHONY: help build test lint fmt docker-build deploy-local

# 变量定义
APP_NAME=ops-ai-assistant
VERSION=$(shell git describe --tags --always --dirty)
BUILD_TIME=$(shell date -u '+%Y-%m-%d_%H:%M:%S')
GO_VERSION=1.21

# 帮助信息
help:
	@echo "可用命令:"
	@echo "  make build          - 编译所有服务"
	@echo "  make test           - 运行单元测试"
	@echo "  make test-integration - 运行集成测试"
	@echo "  make lint           - 代码检查"
	@echo "  make fmt            - 代码格式化"
	@echo "  make docker-build   - 构建Docker镜像"
	@echo "  make dev-up         - 启动本地开发环境"
	@echo "  make dev-down       - 停止本地开发环境"
	@echo "  make migrate-up     - 执行数据库迁移"
	@echo "  make migrate-down   - 回滚数据库迁移"

# 编译
build:
	@echo "编译 API 服务..."
	go build -o bin/api ./cmd/api
	@echo "编译 Worker 服务..."
	go build -o bin/worker ./cmd/worker
	@echo "编译 CLI 工具..."
	go build -o bin/cli ./cmd/cli

# 运行服务
run-api:
	go run ./cmd/api

run-worker:
	go run ./cmd/worker

# 测试
test:
	go test -v -race -coverprofile=coverage.out ./...

test-integration:
	go test -v -tags=integration ./tests/integration/...

coverage:
	go tool cover -html=coverage.out

# 代码质量
lint:
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	goimports -w .

# Docker
docker-build:
	docker build -t $(APP_NAME)-api:$(VERSION) -f deployments/docker/Dockerfile.api .
	docker build -t $(APP_NAME)-worker:$(VERSION) -f deployments/docker/Dockerfile.worker .

# 本地开发环境
dev-up:
	docker-compose -f deployments/docker/docker-compose.yml up -d

dev-down:
	docker-compose -f deployments/docker/docker-compose.yml down

# 数据库迁移
migrate-up:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/ops_ai?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://user:pass@localhost:5432/ops_ai?sslmode=disable" down 1

migrate-create:
	@read -p "迁移名称: " name; \
	migrate create -ext sql -dir migrations -seq $$name

# 清理
clean:
	rm -rf bin/
	rm -f coverage.out

# 依赖管理
deps:
	go mod download
	go mod tidy

# 生成代码
generate:
	go generate ./...

# Swagger文档
swagger:
	swag init -g cmd/api/main.go -o docs/swagger
