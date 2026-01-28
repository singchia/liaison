# Edge 安装脚本

## 目录结构

```
dist/edge/
├── install.sh      # 安装脚本（Linux/macOS，会根据操作系统自动下载对应的安装包）
├── install.ps1     # 安装脚本（Windows PowerShell，推荐）
├── install.bat     # 安装脚本（Windows 批处理脚本）
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

#### Linux/macOS

```bash
# 基本用法（使用默认服务器地址 localhost:8080）
curl -sSL http://your-server:8080/install.sh | bash -s -- \
  --access-key=YOUR_ACCESS_KEY \
  --secret-key=YOUR_SECRET_KEY

# 指定服务器地址（不需要 http:// 前缀）
curl -sSL http://your-server:8080/install.sh | bash -s -- \
  --access-key=YOUR_ACCESS_KEY \
  --secret-key=YOUR_SECRET_KEY \
  --server-http-addr=your-server:443 \
  --server-edge-addr=your-server:30012
```

#### Windows

**方法 1：使用 PowerShell 脚本（推荐）**

```powershell
# 下载并运行 PowerShell 脚本
powershell -ExecutionPolicy Bypass -Command "Invoke-WebRequest -Uri 'https://your-server/install.ps1' -OutFile 'install.ps1'; .\install.ps1 -AccessKey YOUR_ACCESS_KEY -SecretKey YOUR_SECRET_KEY -ServerHttpAddr your-server:443 -ServerEdgeAddr your-server:30012"
```

或者先下载脚本，然后运行：

```powershell
# 下载脚本
Invoke-WebRequest -Uri 'https://your-server/install.ps1' -OutFile 'install.ps1'

# 运行脚本
.\install.ps1 -AccessKey YOUR_ACCESS_KEY -SecretKey YOUR_SECRET_KEY -ServerHttpAddr your-server:443 -ServerEdgeAddr your-server:30012
```

**方法 2：使用批处理脚本**

```cmd
# 下载 install.bat 后直接运行
install.bat --access-key=YOUR_ACCESS_KEY --secret-key=YOUR_SECRET_KEY --server-http-addr=your-server:443 --server-edge-addr=your-server:30012
```

**方法 3：使用 Git Bash 或 WSL**

如果已安装 Git Bash 或 WSL，可以使用 `install.sh` 脚本：

```bash
# 在 Git Bash 中运行
curl -k -sSL https://your-server/install.sh | bash -s -- \
  --access-key=YOUR_ACCESS_KEY \
  --secret-key=YOUR_SECRET_KEY \
  --server-http-addr=your-server:443 \
  --server-edge-addr=your-server:30012
```

**注意：** 如果在 Windows CMD 或 PowerShell 中直接运行 `install.sh`，会提示 "bash不是内部命令或外部命令"。请使用 `install.ps1` 或通过 Git Bash/WSL 运行。

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

