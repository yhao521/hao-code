package main

import (
	"embed"
	goruntime "runtime" // 重命名标准库 runtime 避免冲突

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"

	"hao-code/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Create an instance of the app structure
	// App 内部使用了分层架构和接口解耦设计，便于未来迁移到 Wails v3
	app := backend.NewApp()

	// 根据操作系统配置不同的选项
	isMacOS := goruntime.GOOS == "darwin"

	// Create application with options
	err := wails.Run(&options.App{
		Title:  "Hao-Code Editor",
		Width:  1280,
		Height: 800,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 30, G: 30, B: 30, A: 1}, // VSCode 深色背景
		OnStartup:        app.Startup,
		Bind: []any{
			app, // 绑定 App 实例，暴露方法给前端
		},
		// macOS: 有边框窗口 + 系统菜单栏 + 系统交通灯按钮
		Mac: &mac.Options{
			TitleBar: mac.TitleBarDefault(), // 使用默认标题栏（有边框）
			About: &mac.AboutInfo{
				Title:   "Hao-Code Editor",
				Message: "© 2026 Hao-Code Team\n基于 Wails v2 + Vue 3 构建的跨平台代码编辑器",
			},
		},
		// Windows: 无边框窗口 + 系统菜单栏 + 前端自定义窗口控制按钮
		Windows: &windows.Options{
			WebviewIsTransparent:              false,
			WindowIsTranslucent:               false,
			DisableFramelessWindowDecorations: true, // 无边框窗口
		},
		// macOS: false（有边框）, Windows: true（无边框）
		Frameless: !isMacOS,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}