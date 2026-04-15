#!/bin/bash
# Hao-Code Editor 开发模式启动脚本
# 由于 Wails v3 CLI 尚未正式发布，使用此脚本实现类似功能

echo "🚀 Starting Hao-Code Editor in development mode..."
echo ""

# 进入项目目录
cd "$(dirname "$0")"

# 检查前端依赖
if [ ! -d "frontend/node_modules" ]; then
    echo "📦 Installing frontend dependencies..."
    cd frontend && npm install && cd ..
fi

# 清理之前的构建
echo "🧹 Cleaning previous builds..."
rm -rf frontend/dist

# 构建前端
echo "🔨 Building frontend..."
cd frontend
npm run build
cd ..

echo ""
echo "✅ Frontend build complete!"
echo ""
echo "🏗️  Building and running backend..."
echo ""

# 编译并运行 Go 后端
go run main.go

echo ""
echo "👋 Application closed."
