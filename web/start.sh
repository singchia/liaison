#!/bin/bash

echo "🚀 启动 Liaison 产品管理系统..."

# 设置正确的Node.js版本
export PATH="/usr/local/opt/node@20/bin:$PATH"

# 检查Node.js版本
echo "📦 当前Node.js版本: $(node --version)"

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖..."
    npm install
fi

# 启动开发服务器
echo "🔥 启动开发服务器..."
echo "🌐 访问地址: http://localhost:3000"
npm run dev 