.PHONY: all test build clean fmt lint

all: fmt lint test build

# 运行测试
test:
	go test -v ./...

# 运行测试并生成覆盖率报告
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# 构建
build:
	go build ./...

# 清理
clean:
	go clean
	rm -f coverage.out coverage.html

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
lint:
	golangci-lint run --fix

# 安装依赖
deps:
	go mod download
	go mod tidy

# 运行示例
example-basic:
	go run examples/basic/main.go

example-message:
	go run examples/message/main.go

example-stream:
	go run examples/stream/main.go
