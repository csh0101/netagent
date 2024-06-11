# 定义目标文件名
TARGET = netagent.git

# 定义Go源文件
SRC = main.go

# 伪目标
.PHONY: all build clean run

# 默认目标
all: build

# 编译目标
build:
	@echo "Building $(TARGET)..."
	@go build -o $(TARGET) $(SRC)
	@echo "Build completed."

# 清理目标
clean:
	@echo "Cleaning up..."
	@rm -f $(TARGET)
	@echo "Clean completed."

# 运行目标
run: build
	@echo "Running $(TARGET)..."
	@./$(TARGET)

# 帮助信息
help:
	@echo "Usage:"
	@echo "  make         - Build the project (default

