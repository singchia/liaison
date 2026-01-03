# Edge 安装脚本

## 目录结构

```
dist/edge/
├── install.sh      # 安装脚本（会根据操作系统自动下载对应的安装包）
└── README.md       # 本文件
```

## 安装包目录

安装包需要放置在 `/opt/liaison/packages/edge/` 目录下，文件命名格式：

- Linux: `liaison-edge-linux-amd64.tar.gz` 或 `liaison-edge-linux-arm64.tar.gz`
- macOS: `liaison-edge-darwin-amd64.tar.gz` 或 `liaison-edge-darwin-arm64.tar.gz`
- Windows: `liaison-edge-windows-amd64.tar.gz`

## 安装脚本使用

安装脚本会：
1. 自动检测操作系统和架构
2. 从服务器下载对应的安装包（tar.gz 格式）
3. 解压并安装到 `/opt/liaison/edge/` 目录
4. 创建配置文件，包含 access_key、secret_key 和 server_addr

### 使用方法

```bash
# 基本用法（使用默认服务器地址 localhost:8080）
curl -sSL http://your-server:8080/install.sh | bash -s -- \
  --access-key=YOUR_ACCESS_KEY \
  --secret-key=YOUR_SECRET_KEY

# 指定服务器地址（不需要 http:// 前缀）
curl -sSL http://your-server:8080/install.sh | bash -s -- \
  --access-key=YOUR_ACCESS_KEY \
  --secret-key=YOUR_SECRET_KEY \
  --server-addr=your-server:8080
```

### 参数说明

- `--access-key`: Access Key（必需）
- `--secret-key`: Secret Key（必需）
- `--server-addr`: 服务器地址，格式为 `host:port`（可选，默认：localhost:8080）
- `--help`: 显示帮助信息

## 配置说明

在 `liaison.yaml` 中配置：

```yaml
manager:
  server_url: "http://your-server:8080"  # 服务器 URL，用于生成安装命令（包含协议）
  packages_dir: "/opt/liaison/packages"   # 安装包目录（可选，默认 /opt/liaison/packages）
```

注意：`server_url` 用于生成安装命令的下载 URL（需要包含协议），而安装脚本中的 `--server-addr` 参数只需要地址部分（不需要 `http://` 前缀）。

## API 端点

- `/install.sh` - 安装脚本下载
- `/packages/edge/{package-name}` - 安装包下载

这些端点不需要认证即可访问。

