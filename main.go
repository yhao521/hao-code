package main

import (
	"context"
	"embed"
	goruntime "runtime" // 重命名标准库 runtime 避免冲突

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/mac"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"hao-code/backend"
)

//go:embed all:frontend/dist
var assets embed.FS

// 全局上下文变量，用于菜单回调
var appCtx context.Context

func createMenu(isMacOS bool) *menu.Menu {
	appMenu := menu.NewMenu()

	// macOS 系统菜单
	if isMacOS {
		// 添加 macOS 特有的应用菜单（Hao-Code Editor 菜单）
		appMenu.Append(menu.AppMenu())

		// 文件菜单
		fileMenu := appMenu.AddSubmenu("文件")

		// 新建相关
		fileMenu.AddText("新建文本文件", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:new-text-file")
		})
		fileMenu.AddText("新建文件...", keys.Combo("n", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:new-file")
		})
		fileMenu.AddSeparator()

		// 打开相关
		fileMenu.AddText("打开...", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-file")
		})
		fileMenu.AddText("打开文件夹...", keys.Combo("o", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-folder")
		})
		fileMenu.AddText("打开最近的文件", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-recent")
		})
		fileMenu.AddSeparator()

		// 保存相关
		fileMenu.AddText("保存", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:save")
		})
		fileMenu.AddText("另存为...", keys.Combo("s", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:save-as")
		})
		fileMenu.AddText("全部保存", keys.Combo("s", keys.CmdOrCtrlKey, keys.OptionOrAltKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:save-all")
		})
		fileMenu.AddSeparator()

		// 自动保存
		fileMenu.AddCheckbox("自动保存", true, nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:toggle-auto-save")
		})
		fileMenu.AddSeparator()

		// 关闭相关
		fileMenu.AddText("关闭编辑器", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:close-editor")
		})
		fileMenu.AddText("关闭文件夹", keys.Combo("k", keys.CmdOrCtrlKey, keys.ControlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:close-folder")
		})
		fileMenu.AddText("关闭窗口", keys.Combo("w", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.Quit(appCtx)
		})

		// 编辑菜单
		appMenu.Append(menu.EditMenu())

		// 帮助菜单
		helpMenu := appMenu.AddSubmenu("帮助")
		helpMenu.AddText("欢迎", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:welcome")
		})
		helpMenu.AddText("显示所有命令", keys.Combo("p", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:show-all-commands")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("文档", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:documentation")
		})
		helpMenu.AddText("视频教程", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:video-tutorials")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("键盘快捷方式参考", keys.Combo("k", keys.CmdOrCtrlKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:keyboard-shortcuts")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("搜索功能请求", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:search-feature-requests")
		})
		helpMenu.AddText("使用英文报告问题", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:report-issues")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("查看许可证", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:view-license")
		})
		helpMenu.AddText("隐私声明", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:privacy-statement")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("切换开发人员工具", keys.CmdOrCtrl("i"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:toggle-devtools")
		})
		helpMenu.AddText("打开进程资源管理器", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-process-explorer")
		})

		// 窗口菜单
		appMenu.Append(menu.WindowMenu())
	} else {
		// Windows/Linux 也需要菜单
		fileMenu := appMenu.AddSubmenu("文件")

		fileMenu.AddText("新建文本文件", keys.CmdOrCtrl("n"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:new-text-file")
		})
		fileMenu.AddText("新建文件...", keys.Combo("n", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:new-file")
		})
		fileMenu.AddSeparator()
		fileMenu.AddText("打开...", keys.CmdOrCtrl("o"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-file")
		})
		fileMenu.AddText("打开文件夹...", keys.Combo("o", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-folder")
		})
		fileMenu.AddText("打开最近的文件", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-recent")
		})
		fileMenu.AddSeparator()
		fileMenu.AddText("保存", keys.CmdOrCtrl("s"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:save")
		})
		fileMenu.AddText("另存为...", keys.Combo("s", keys.CmdOrCtrlKey, keys.ShiftKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:save-as")
		})
		fileMenu.AddSeparator()
		fileMenu.AddCheckbox("自动保存", true, nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:toggle-auto-save")
		})
		fileMenu.AddSeparator()
		fileMenu.AddText("关闭文件夹", keys.Combo("k", keys.CmdOrCtrlKey, keys.ControlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:close-folder")
		})
		fileMenu.AddText("关闭窗口", keys.CmdOrCtrl("w"), func(_ *menu.CallbackData) {
			runtime.Quit(appCtx)
		})

		appMenu.Append(menu.EditMenu())

		// 帮助菜单
		helpMenu := appMenu.AddSubmenu("帮助")
		helpMenu.AddText("欢迎", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:welcome")
		})
		helpMenu.AddText("显示所有命令", keys.Combo("p", keys.ShiftKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:show-all-commands")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("文档", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:documentation")
		})
		helpMenu.AddText("视频教程", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:video-tutorials")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("键盘快捷方式参考", keys.Combo("k", keys.CmdOrCtrlKey, keys.CmdOrCtrlKey), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:keyboard-shortcuts")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("搜索功能请求", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:search-feature-requests")
		})
		helpMenu.AddText("使用英文报告问题", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:report-issues")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("查看许可证", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:view-license")
		})
		helpMenu.AddText("隐私声明", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:privacy-statement")
		})
		helpMenu.AddSeparator()
		helpMenu.AddText("切换开发人员工具", keys.CmdOrCtrl("i"), func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:toggle-devtools")
		})
		helpMenu.AddText("打开进程资源管理器", nil, func(_ *menu.CallbackData) {
			runtime.EventsEmit(appCtx, "menu:open-process-explorer")
		})
	}

	return appMenu
}

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
		OnStartup: func(ctx context.Context) {
			// 保存上下文供菜单使用
			appCtx = ctx
			// 初始化 app
			app.Startup(ctx)
		},
		Menu: createMenu(isMacOS), // 设置系统菜单
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
