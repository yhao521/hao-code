#!/bin/bash

echo "🚀 Starting Hao-Code Editor in development mode..."
echo ""

# 进入项目目录
cd "$(dirname "$0")"

# 检查依赖
if ! command -v wails &> /dev/null; then
    echo "❌ Wails is not installed. Installing..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
fi

# 安装前端依赖
echo "📦 Checking frontend dependencies..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi
cd ..

echo ""
echo "🔨 Building and running Hao-Code Editor..."
echo ""

# 启动开发模式
wails dev
