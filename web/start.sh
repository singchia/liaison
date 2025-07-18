#!/bin/bash

echo "🚀 启动 Liaison 产品管理系统..."

# 检查是否已安装依赖
if [ ! -d "node_modules" ]; then
    echo "📦 安装依赖..."
    npm install
fi

# 启动开发服务器
echo "🔥 启动开发服务器..."
npm run dev 