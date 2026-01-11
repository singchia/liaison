# Liaison

网络马上通达，Liaison 让你可以轻松连接分布在不同位置的设备与应用。无论你的设备在哪里（家、办公室或机房），都能通过 Liaison 建立稳定、安全的网络连接。

---

## 为什么选择 Liaison？

- **简单易用** - 通过 Web 界面即可完成所有操作，无需复杂的命令行配置
- **跨平台支持** - 支持 Linux、macOS、Windows 等多种操作系统
- **自动发现** - 自动发现设备上的应用和服务，无需手动配置
- **安全可靠** - 内置身份认证机制，保障连接安全
- **实时监控** - 实时查看设备状态和资源使用情况

## 主要功能

### 设备连接

自动发现并连接你的所有设备，实时查看设备状态、CPU、内存、磁盘等资源使用情况。无论设备在哪里，都能一目了然。

![设备管理](docs/pages/device.png)

### 连接器

创建多个连接器，每个连接器可以部署在不同的设备上。支持一键安装，几分钟内即可完成部署。

![连接器管理](docs/pages/connector.png)

### 应用发现

自动扫描并发现设备上运行的应用和服务，无需手动配置端口信息。轻松掌控你的应用状态。

![应用管理](docs/pages/application.png)

### 流量转发

灵活配置网络代理规则，实现流量转发和负载均衡。让你的网络流量按需分配，提高资源利用率。

![代理管理](docs/pages/proxy.png)

## 快速开始

### 安装服务端

1. 下载安装包 `liaison-v1.1.0-linux-amd64.tar.gz` 并解压
2. 进入解压后的目录，运行安装脚本：
   ```bash
   tar -xzf liaison-v1.1.0-linux-amd64.tar.gz
   cd liaison-v1.1.0-linux-amd64
   sudo ./install.sh
   ```
3. 访问 Web 控制台（默认地址：https://localhost）

### 添加设备

1. 在 Web 控制台中创建连接器，获取 Access Key 和 Secret Key
2. 在目标设备上运行安装命令：
   ```bash
   curl -sSL https://your-server/install.sh | bash -s -- \
     --access-key=YOUR_ACCESS_KEY \
     --secret-key=YOUR_SECRET_KEY \
     ...
   ```
3. 等待几秒钟，设备会自动出现在控制台中

就这么简单！现在你可以开始使用你的设备了。

## 使用场景

- **远程办公/开发** - 连接办公室和家中的设备，随时随地访问
- **多机房部署** - 统一连接分布在不同机房的服务器和设备
- **边缘计算** - 连接和监控边缘设备上的应用和服务
- **内网穿透** - 让内网服务可以通过公网安全访问

## 系统要求

- **服务端**: Linux 系统
- **连接器**: Linux、macOS、Windows（支持 x86_64 和 ARM64 架构）

## 许可证

Apache License 2.0

详见 [LICENSE](LICENSE) 文件。

## 版本

当前版本: v1.1.0

---

如有问题或建议，欢迎提交 Issue 或 PR。
