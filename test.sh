#!/bin/bash

echo "🧪 测试 Hao-Code Editor 核心功能"
echo ""

# 检查 Go 依赖
echo "📦 检查 Go 依赖..."
cd "$(dirname "$0")"
go mod tidy

if [ $? -ne 0 ]; then
    echo "❌ Go 依赖检查失败"
    exit 1
fi

echo "✅ Go 依赖检查通过"
echo ""

# 检查前端依赖
echo "📦 检查前端依赖..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "安装前端依赖..."
    npm install
fi

echo "✅ 前端依赖检查通过"
echo ""

# 构建前端
echo "🔨 构建前端..."
npm run build

if [ $? -ne 0 ]; then
    echo "❌ 前端构建失败"
    exit 1
fi

echo "✅ 前端构建成功"
echo ""

cd ..

# 启动 Wails
echo "🚀 启动 Hao-Code Editor..."
echo ""
echo "应用即将启动，请观察窗口是否正常打开..."
echo "按 Ctrl+C 可以停止应用"
echo ""

wails dev
