# Makefile E2E测试使用指南

## 新增的Makefile Targets

### 测试相关Targets

```bash
# 安装测试依赖
make test-deps

# 运行所有单元测试
make test

# 运行简化E2E测试（推荐用于快速验证）
make test-e2e

# 运行完整E2E测试（需要默认用户）
make test-e2e-full

# 运行所有测试（单元测试 + E2E测试）
make test-all
```

### 构建相关Targets

```bash
# 构建所有组件并安装测试依赖
make all

# 构建manager服务
make liaison

# 构建edge服务
make liaison-edge

# 构建Linux版本
make linux
```

### API生成相关Targets

```bash
# 生成protobuf代码
make gen-api

# 生成swagger文档
make gen-swagger
```

## 使用示例

### 开发流程

1. **首次设置**：
   ```bash
   make all  # 构建所有组件并安装依赖
   ```

2. **日常开发**：
   ```bash
   make test-e2e  # 快速验证功能
   ```

3. **完整测试**：
   ```bash
   make test-all  # 运行所有测试
   ```

4. **API更新后**：
   ```bash
   make gen-api  # 重新生成protobuf代码
   make liaison  # 重新构建
   make test-e2e  # 验证更新
   ```

### CI/CD集成

在CI/CD流水线中可以这样使用：

```bash
# 构建和测试
make all
make test-all

# 或者分步执行
make test-deps
make liaison
make test-e2e
```

## 注意事项

1. **E2E测试要求**：
   - 服务器需要能够启动
   - 端口8080需要可用
   - 数据库需要可写权限

2. **测试失败处理**：
   - 如果健康检查失败，检查服务器是否正常启动
   - 如果认证测试失败，检查中间件是否正确配置
   - 如果端口冲突，停止现有服务或修改配置

3. **环境要求**：
   - Go 1.23+
   - Docker（用于protobuf生成）
   - curl（用于健康检查）

## 故障排除

### 常见问题

1. **"Server failed to start"**：
   - 检查端口8080是否被占用
   - 检查配置文件路径是否正确
   - 检查数据库权限

2. **"Health check failed"**：
   - 确认健康检查端点已正确注册
   - 检查服务器是否完全启动

3. **"Authentication test failed"**：
   - 确认认证中间件已正确配置
   - 检查IAM服务是否正常工作

### 调试技巧

```bash
# 手动启动服务器进行调试
./bin/liaison -c etc/liaison.yaml

# 在另一个终端运行测试
make test-e2e

# 查看详细测试输出
cd test/e2e && go test -v .
```
