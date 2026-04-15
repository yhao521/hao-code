# Wails v3 对话框功能实现指南

## ✅ 已完成的功能

已成功实现以下三个对话框功能：

1. **OpenFolderDialog** - 打开文件夹选择对话框
2. **OpenFileDialog** - 打开文件选择对话框
3. **SaveFileDialog** - 保存文件对话框

---

## 📝 实现细节

### 1. OpenFolderDialog（打开文件夹）

```go
func (f *FileSystemService) OpenFolderDialog() (string, error) {
    // 使用 Wails v3 的全局应用实例
    app := application.Get()
    if app == nil {
        return "", fmt.Errorf("application not initialized")
    }

    // 使用 Wails v3 的对话框 API
    dialog := app.Dialog.OpenFile()
    dialog.SetTitle("选择项目文件夹")
    dialog.CanChooseFiles(false)      // 不允许选择文件
    dialog.CanChooseDirectories(true)  // 允许选择目录
    dialog.CanCreateDirectories(true)  // 允许创建新目录

    result, err := dialog.PromptForSingleSelection()
    if err != nil {
        return "", err
    }

    if result == "" {
        return "", fmt.Errorf("user cancelled")
    }

    return result, nil
}
```

**特性**：

- ✅ 只能选择文件夹，不能选择文件
- ✅ 支持创建新文件夹
- ✅ 用户取消时返回错误

---

### 2. OpenFileDialog（打开文件）

```go
func (f *FileSystemService) OpenFileDialog() (string, error) {
    app := application.Get()
    if app == nil {
        return "", fmt.Errorf("application not initialized")
    }

    dialog := app.Dialog.OpenFile()
    dialog.SetTitle("选择文件")
    dialog.CanChooseFiles(true)        // 允许选择文件
    dialog.CanChooseDirectories(false) // 不允许选择目录
    dialog.CanCreateDirectories(true)  // 允许创建新目录

    result, err := dialog.PromptForSingleSelection()
    if err != nil {
        return "", err
    }

    if result == "" {
        return "", fmt.Errorf("user cancelled")
    }

    return result, nil
}
```

**特性**：

- ✅ 只能选择文件，不能选择文件夹
- ✅ 支持创建新文件夹（用于导航）
- ✅ 单选模式

---

### 3. SaveFileDialog（保存文件）

```go
func (f *FileSystemService) SaveFileDialog() (string, error) {
    app := application.Get()
    if app == nil {
        return "", fmt.Errorf("application not initialized")
    }

    dialog := app.Dialog.SaveFile()
    dialog.SetMessage("保存文件")
    dialog.SetButtonText("保存")
    dialog.CanCreateDirectories(true)

    result, err := dialog.PromptForSingleSelection()
    if err != nil {
        return "", err
    }

    if result == "" {
        return "", fmt.Errorf("user cancelled")
    }

    return result, nil
}
```

**特性**：

- ✅ 保存文件对话框
- ✅ 自定义按钮文本
- ✅ 支持创建新文件夹

---

## 🔧 技术要点

### 1. 使用全局应用实例

Wails v3 提供了 `application.Get()` 方法来获取全局应用实例：

```go
import "github.com/wailsapp/wails/v3/pkg/application"

app := application.Get()
if app == nil {
    return "", fmt.Errorf("application not initialized")
}
```

**优势**：

- 无需在 service 中存储 app 引用
- 简化了依赖注入
- 符合 Wails v3 的设计理念

### 2. 对话框 API 变化

| Wails v2                                 | Wails v3                                           | 说明         |
| ---------------------------------------- | -------------------------------------------------- | ------------ |
| `runtime.OpenDirectoryDialog(ctx, opts)` | `app.Dialog.OpenFile().CanChooseDirectories(true)` | 更灵活的配置 |
| `runtime.OpenFileDialog(ctx, opts)`      | `app.Dialog.OpenFile()`                            | 链式调用     |
| `runtime.SaveFileDialog(ctx, opts)`      | `app.Dialog.SaveFile()`                            | 简化的 API   |

### 3. 错误处理

所有对话框方法都遵循统一的错误处理模式：

```go
result, err := dialog.PromptForSingleSelection()
if err != nil {
    return "", err  // 系统错误
}

if result == "" {
    return "", fmt.Errorf("user cancelled")  // 用户取消
}

return result, nil  // 成功
```

---

## 🚀 使用方法

### 前端调用示例

```typescript
// 打开文件夹
async function openFolder() {
  try {
    const path = await OpenFolderDialog();
    console.log("选择的文件夹:", path);
    // 处理文件夹路径
  } catch (error) {
    if (error === "user cancelled") {
      console.log("用户取消了操作");
    } else {
      console.error("打开文件夹失败:", error);
    }
  }
}

// 打开文件
async function openFile() {
  try {
    const path = await OpenFileDialog();
    console.log("选择的文件:", path);
    // 读取文件内容
  } catch (error) {
    console.error("打开文件失败:", error);
  }
}

// 保存文件
async function saveFile() {
  try {
    const path = await SaveFileDialog();
    console.log("保存位置:", path);
    // 写入文件内容
  } catch (error) {
    console.error("保存文件失败:", error);
  }
}
```

---

## 🎨 高级用法

### 1. 添加文件过滤器

```go
dialog := app.Dialog.OpenFile()
dialog.AddFilter("Go 文件", "*.go")
dialog.AddFilter("所有文件", "*")
```

### 2. 设置默认目录

```go
dialog := app.Dialog.OpenFile()
dialog.SetDirectory("/Users/username/Documents")
```

### 3. 多选文件

```go
dialog := app.Dialog.OpenFile()
selections, err := dialog.PromptForMultipleSelection()
if err == nil {
    for _, file := range selections {
        fmt.Println(file)
    }
}
```

### 4. 自定义消息对话框

```go
// 信息对话框
app.Dialog.Info().
    SetTitle("提示").
    SetMessage("操作成功！").
    Show()

// 警告对话框
app.Dialog.Warning().
    SetTitle("警告").
    SetMessage("确定要删除吗？").
    AddButton("确定").OnClick(func() {
        // 执行删除
    }).
    AddButton("取消").SetAsCancel().
    Show()
```

---

## ⚠️ 注意事项

1. **必须在应用启动后使用**
   - 对话框需要访问全局应用实例
   - 确保在 `application.Run()` 之后调用

2. **主线程限制**
   - 对话框必须在主线程调用
   - Wails v3 会自动处理线程问题

3. **平台差异**
   - macOS、Windows、Linux 的对话框样式可能不同
   - 某些选项在某些平台上可能被忽略

4. **用户取消处理**
   - 用户取消时返回空字符串
   - 建议检查返回值并适当处理

---

## 📊 性能对比

| 功能       | Wails v2     | Wails v3     | 改进        |
| ---------- | ------------ | ------------ | ----------- |
| API 简洁度 | 需要 context | 无需 context | ✅ 更简洁   |
| 配置灵活性 | 固定选项     | 链式调用     | ✅ 更灵活   |
| 类型安全   | 一般         | 强类型       | ✅ 更安全   |
| 代码量     | ~15 行       | ~10 行       | ✅ 减少 33% |

---

## 🐛 常见问题

### Q1: 为什么 `application.Get()` 返回 nil？

**A**: 应用尚未初始化。确保在 `application.Run()` 之后调用对话框。

### Q2: 如何添加多个文件过滤器？

**A**: 多次调用 `AddFilter()` 方法：

```go
dialog.AddFilter("图片", "*.jpg;*.png;*.gif")
dialog.AddFilter("文档", "*.pdf;*.doc;*.txt")
dialog.AddFilter("所有文件", "*")
```

### Q3: 对话框不显示怎么办？

**A**: 检查以下几点：

1. 确认应用已正确初始化
2. 检查是否有编译错误
3. 查看控制台输出是否有错误信息
4. 确保在主线程调用

---

## 📚 相关文档

- [Wails v3 官方文档 - Dialogs](https://v3.wails.io/docs/reference/dialogs)
- [WAILS_V3_UPGRADE.md](./WAILS_V3_UPGRADE.md) - 完整升级报告
- [backend/file_service.go](./backend/file_service.go) - 实现源码

---

**实现日期**: 2026-04-15  
**Wails 版本**: v3.0.0-alpha.74  
**状态**: ✅ 已完成并测试通过
