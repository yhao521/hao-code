package main

import (
	"embed"
	goruntime "runtime" // 重命名标准库 runtime 避免冲突

	"github.com/wailsapp/wails/v3/pkg/application"
	"github.com/wailsapp/wails/v3/pkg/events"

	"hao-code/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

// 全局应用实例，用于菜单回调
var wailsApp *application.App
var mainWindow *application.WebviewWindow
var services *backend.ServiceContainer

// updateRecentMenu 更新最近文件子菜单
func updateRecentMenu(recentMenu *application.Menu) {
	// 清空旧菜单项
	// 注意：Wails v3 菜单 API 可能不支持清空，我们需要重建

	// 获取最近文件列表
	recentFiles := services.App.GetRecentFiles()
	for i, file := range recentFiles {
		if i >= 10 { // 最多显示 10 个文件
			break
		}
		fileCopy := file // 创建副本避免闭包问题
		recentMenu.Add(fileCopy.Name).OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-recent-file", fileCopy.Path)
		})
	}

	if len(recentFiles) > 0 {
		recentMenu.AddSeparator()
	}

	// 获取最近文件夹列表
	recentFolders := services.App.GetRecentFolders()
	for i, folder := range recentFolders {
		if i >= 5 { // 最多显示 5 个文件夹
			break
		}
		folderCopy := folder // 创建副本避免闭包问题
		recentMenu.Add("📁 " + folderCopy.Name).OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-recent-folder", folderCopy.Path)
		})
	}

	if len(recentFiles) == 0 && len(recentFolders) == 0 {
		recentMenu.Add("空").SetEnabled(false)
	}
}

func createMenu(app *application.App, isMacOS bool) {
	// 创建应用菜单
	appMenu := app.NewMenu()

	if isMacOS {
		// macOS 系统菜单结构
		// 应用名称菜单（macOS 特有）
		appInfoMenu := appMenu.AddSubmenu("Hao-Code Editor")
		appInfoMenu.Add("关于 Hao-Code Editor").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:about")
		})
		appInfoMenu.AddSeparator()
		appInfoMenu.Add("偏好设置").SetAccelerator("CmdOrCtrl+,")
		appInfoMenu.AddSeparator()
		appInfoMenu.Add("退出 Hao-Code Editor").SetAccelerator("CmdOrCtrl+Q").OnClick(func(ctx *application.Context) {
			app.Quit()
		})

		// 文件菜单
		fileMenu := appMenu.AddSubmenu("文件")

		// 新建相关
		fileMenu.Add("新建文本文件").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:new-text-file")
		})
		fileMenu.Add("新建文件...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:new-file")
		})
		fileMenu.AddSeparator()

		// 打开相关
		fileMenu.Add("打开...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-file")
		})
		fileMenu.Add("打开文件夹...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-folder")
		})

		// 打开最近的文件 - 改为子菜单
		recentMenu := fileMenu.AddSubmenu("打开最近的文件")

		// 动态添加最近文件和文件夹
		updateRecentMenu(recentMenu)

		fileMenu.AddSeparator()

		// 保存相关
		fileMenu.Add("保存").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:save")
		})
		fileMenu.Add("另存为...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:save-as")
		})
		fileMenu.Add("全部保存").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:save-all")
		})
		fileMenu.AddSeparator()

		// 自动保存
		fileMenu.AddCheckbox("自动保存", true).OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:toggle-auto-save")
		})
		fileMenu.AddSeparator()

		// 关闭相关
		fileMenu.Add("关闭编辑器").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:close-editor")
		})
		fileMenu.Add("关闭文件夹").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:close-folder")
		})
		fileMenu.Add("关闭窗口").OnClick(func(ctx *application.Context) {
			app.Quit()
		})

		// 编辑菜单（使用系统默认）
		editMenu := appMenu.AddSubmenu("编辑")
		editMenu.Add("撤销").SetAccelerator("CmdOrCtrl+Z")
		editMenu.Add("重做").SetAccelerator("CmdOrCtrl+Shift+Z")
		editMenu.AddSeparator()
		editMenu.Add("剪切").SetAccelerator("CmdOrCtrl+X")
		editMenu.Add("复制").SetAccelerator("CmdOrCtrl+C")
		editMenu.Add("粘贴").SetAccelerator("CmdOrCtrl+V")

		// 帮助菜单
		helpMenu := appMenu.AddSubmenu("帮助")
		helpMenu.Add("欢迎").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:welcome")
		})
		helpMenu.Add("显示所有命令").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:show-all-commands")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("文档").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:documentation")
		})
		helpMenu.Add("视频教程").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:video-tutorials")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("键盘快捷方式参考").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:keyboard-shortcuts")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("搜索功能请求").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:search-feature-requests")
		})
		helpMenu.Add("使用英文报告问题").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:report-issues")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("查看许可证").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:view-license")
		})
		helpMenu.Add("隐私声明").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:privacy-statement")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("切换开发人员工具").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:toggle-devtools")
		})
		helpMenu.Add("打开进程资源管理器").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-process-explorer")
		})

		// 窗口菜单
		windowMenu := appMenu.AddSubmenu("窗口")
		windowMenu.Add("最小化").SetAccelerator("CmdOrCtrl+M")
		windowMenu.Add("缩放")
		windowMenu.AddSeparator()
		windowMenu.Add("前置全部窗口")
	} else {
		// Windows/Linux 菜单
		fileMenu := appMenu.AddSubmenu("文件")

		fileMenu.Add("新建文本文件").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:new-text-file")
		})
		fileMenu.Add("新建文件...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:new-file")
		})
		fileMenu.AddSeparator()
		fileMenu.Add("打开...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-file")
		})
		fileMenu.Add("打开文件夹...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-folder")
		})

		// 打开最近的文件 - 改为子菜单
		recentMenu := fileMenu.AddSubmenu("打开最近的文件")
		updateRecentMenu(recentMenu)

		fileMenu.AddSeparator()
		fileMenu.Add("保存").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:save")
		})
		fileMenu.Add("另存为...").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:save-as")
		})
		fileMenu.AddSeparator()
		fileMenu.AddCheckbox("自动保存", true).OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:toggle-auto-save")
		})
		fileMenu.AddSeparator()
		fileMenu.Add("关闭文件夹").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:close-folder")
		})
		fileMenu.Add("关闭窗口").OnClick(func(ctx *application.Context) {
			app.Quit()
		})

		// 编辑菜单
		editMenu := appMenu.AddSubmenu("编辑")
		editMenu.Add("撤销").SetAccelerator("Ctrl+Z")
		editMenu.Add("重做").SetAccelerator("Ctrl+Shift+Z")
		editMenu.AddSeparator()
		editMenu.Add("剪切").SetAccelerator("Ctrl+X")
		editMenu.Add("复制").SetAccelerator("Ctrl+C")
		editMenu.Add("粘贴").SetAccelerator("Ctrl+V")

		// 帮助菜单
		helpMenu := appMenu.AddSubmenu("帮助")
		helpMenu.Add("欢迎").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:welcome")
		})
		helpMenu.Add("显示所有命令").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:show-all-commands")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("文档").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:documentation")
		})
		helpMenu.Add("视频教程").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:video-tutorials")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("键盘快捷方式参考").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:keyboard-shortcuts")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("搜索功能请求").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:search-feature-requests")
		})
		helpMenu.Add("使用英文报告问题").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:report-issues")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("查看许可证").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:view-license")
		})
		helpMenu.Add("隐私声明").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:privacy-statement")
		})
		helpMenu.AddSeparator()
		helpMenu.Add("切换开发人员工具").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:toggle-devtools")
		})
		helpMenu.Add("打开进程资源管理器").OnClick(func(ctx *application.Context) {
			mainWindow.EmitEvent("menu:open-process-explorer")
		})
	}

	// 设置应用菜单
	wailsApp.Menu.Set(appMenu)
}

func main() {
	// 根据操作系统配置
	isMacOS := goruntime.GOOS == "darwin"

	// 创建 Wails v3 应用实例
	// 创建服务容器并直接使用 AppService
	services := backend.NewServiceContainer()

	wailsApp = application.New(application.Options{
		Name:        "Hao-Code Editor",
		Description: "基于 Wails v3 + Vue 3 构建的跨平台代码编辑器",
		Services: []application.Service{
			// 直接注册 AppService，无需适配器层
			// 需要将接口转换为具体类型
			application.NewService(services.App.(*backend.AppService)),
		},
		Assets: application.AssetOptions{
			Handler: application.BundledAssetFileServer(assets),
		},
		Mac: application.MacOptions{
			ApplicationShouldTerminateAfterLastWindowClosed: true,
		},
	})

	// 创建主窗口
	mainWindow = wailsApp.Window.NewWithOptions(application.WebviewWindowOptions{
		Title:            "Hao-Code Editor",
		URL:              "/",
		Width:            1280,
		Height:           800,
		MinWidth:         800,
		MinHeight:        600,
		DevToolsEnabled:  true,
		Frameless:        !isMacOS,                             // macOS 使用有边框窗口，Windows 使用无边框
		BackgroundColour: application.NewRGBA(30, 30, 30, 255), // VSCode 深色背景
	})

	// 注册 WebSocket 路由用于终端通信
	wailsApp.Router.GET("/ws/terminal", backend.TerminalWebSocketHandler)

	// 监听窗口运行时就绪事件
	mainWindow.RegisterHook(events.Common.WindowRuntimeReady, func(e *application.WindowEvent) {
		// Wails v3 中服务会自动初始化，无需手动调用 Startup
	})

	// 创建系统菜单
	createMenu(wailsApp, isMacOS)

	// 运行应用
	err := wailsApp.Run()
	if err != nil {
		println("Error:", err.Error())
	}
}
