# Liaison Password Verifier

这个工具用于验证Liaison数据库中的用户密码。

## 工具说明

### password-verifier - 密码验证器

用于验证用户密码是否正确，支持Argon2id哈希算法。

#### 构建
```bash
make password-verifier
```

#### 使用方法
```bash
# 验证密码
./bin/password-verifier <email> <password>

# 示例
./bin/password-verifier default@liaison.local mypassword
```

#### 功能
- 验证用户密码是否正确
- 支持Argon2id哈希算法
- 自动查找数据库路径
- 提供常见默认密码建议

## 快速开始

### 构建工具
```bash
make tools
```

### 验证密码
```bash
./bin/password-verifier default@liaison.local default123
```

## 数据库路径

工具会自动查找以下路径的数据库文件：
1. `/opt/liaison/data/liaison.db` (默认)
2. `./etc/liaison.db`
3. `./liaison.db`
4. `./data/liaison.db`

你也可以通过命令行参数指定数据库路径：
```bash
./bin/password-viewer /custom/path/to/database.db
```

## 示例输出

### password-viewer 输出示例
```
🔍 Liaison Database Password Viewer
Database: /opt/liaison/data/liaison.db
==================================================
✅ Found 1 user(s):

👤 User #1
   📧 Email: default@liaison.local
   🔑 Password Hash: $2a$10$abc123...
   📊 Status: 1 (Active)
   📅 Created: 2025-10-14 10:30:00
   🕒 Last Login: Never

📁 Password File Information:
   📍 Location: /Users/username/.liaison/default_password.txt
   ✅ Password file exists!
   📄 Content:
Liaison 默认用户账户信息
邮箱: default@liaison.local
密码: default123
请妥善保管此信息，首次登录后建议修改密码
```

## 故障排除

### 数据库连接失败
- 检查数据库文件是否存在
- 检查文件权限
- 确认数据库路径正确

### 没有找到用户
- 确认users表存在
- 检查表结构是否正确
- 确认数据已插入

### 密码文件不存在
- 这是正常的，如果默认用户还没有创建
- 密码文件会在创建默认用户时生成

## 安全注意事项

⚠️ **重要提醒**：
- 这些工具会显示密码哈希，请确保在安全环境中使用
- 不要在生产环境中运行这些工具
- 密码哈希是加密的，无法直接逆向得到原始密码
- 默认密码文件包含明文密码，请妥善保管

## 开发说明

### 添加新工具
1. 在 `tools/` 目录下创建新的 `.go` 文件
2. 在 `Makefile` 中添加构建规则
3. 更新此README文档

### 依赖
- Go 1.23+
- SQLite3
- CGO enabled
