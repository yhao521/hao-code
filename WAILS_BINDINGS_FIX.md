# Wails v3 绑定问题修复指南

**日期**: 2026年4月15日  
**问题**: `window.go.backend` is undefined  
**状态**: ✅ 已修复

---

## 🐛 问题描述

在浏览器控制台中出现以下错误：

```javascript
[Error] Failed to open folder: – TypeError: undefined is not an object (evaluating 'window.go.backend')
```

**原因**: Wails v3 alpha.74 的绑定生成机制发生了变化，导致前端无法正确访问后端方法。

---

## ✅ 解决方案

### 步骤 1: 重新生成绑定

```bash
cd /Users/yanghao/storage/code_projects/goProjects/hao-code
wails3 generate bindings -clean=true
```

这会生成新的绑定到 `frontend/bindings/` 目录。

### 步骤 2: 复制绑定到兼容位置

```bash
# 将新绑定复制到旧的 wailsjs 目录
cp -r frontend/bindings/hao-code/backend/* frontend/wailsjs/go/backend/
```

### 步骤 3: 创建兼容层

已创建以下兼容文件：
- `frontend/wailsjs/go/backend/App.js` - JavaScript 兼容层
- `frontend/wailsjs/go/backend/App.d.ts` - TypeScript 声明

这些文件导出所有新方法，保持与旧代码的兼容性。

### 步骤 4: 重启应用

**重要**: 必须完全重启应用才能加载新的绑定！

```bash
# 停止当前运行的应用（Ctrl+C）

# 重新启动
task dev
# 或
./dev.sh
```

---

## 🔍 技术细节

### Wails v3 绑定变化

#### 旧版本 (alpha.74 之前)
```javascript
// 直接通过 window.go 访问
export function OpenFolderDialog() {
  return window['go']['backend']['App']['OpenFolderDialog']();
}
```

#### 新版本 (alpha.74+)
```javascript
// 使用 @wailsio/runtime
import { Call as $Call } from "@wailsio/runtime";

export function OpenFolderDialog() {
  return $Call.ByID(123456789);
}
```

### 兼容层设计

为了保持现有代码不变，创建了兼容层：

```javascript
// frontend/wailsjs/go/backend/App.js
export * from './appservice.js';

export {
    BranchInfo,
    Change,
    // ... 其他模型
} from './models.js';
```

这样所有现有的导入语句仍然有效：
```typescript
import { OpenFolderDialog } from '@wails/go/backend/App'
```

---

## 📋 验证修复

### 1. 检查绑定文件

```bash
ls -la frontend/wailsjs/go/backend/
```

应该看到：
- `App.js` - 兼容层
- `App.d.ts` - TypeScript 声明
- `appservice.js` - 新生成的服务
- `models.js` - 数据模型
- `index.js` - 索引文件

### 2. 检查 SearchInFiles 方法

```bash
grep "SearchInFiles" frontend/wailsjs/go/backend/appservice.js
```

应该看到方法定义。

### 3. 测试功能

在应用中：
1. 点击"打开文件夹"按钮
2. 应该能正常打开对话框
3. 控制台不应有错误

---

## 🚨 常见问题

### Q1: 重新生成绑定后仍然报错？

**A**: 确保完全重启了应用：
```bash
# 停止应用（Ctrl+C）
# 等待几秒
# 重新启动
task dev
```

### Q2: 浏览器缓存问题？

**A**: 强制刷新浏览器：
- macOS: `Cmd + Shift + R`
- Windows/Linux: `Ctrl + Shift + R`

或者清除缓存后重新加载。

### Q3: 绑定生成失败？

**A**: 检查 Go 代码是否有编译错误：
```bash
go build
```

如果有错误，先修复 Go 代码，然后重新生成绑定。

### Q4: TypeScript 类型错误？

**A**: 确保 `.d.ts` 文件存在且正确：
```bash
ls -la frontend/wailsjs/go/backend/*.d.ts
```

---

## 🔄 何时需要重新生成绑定

以下情况需要重新生成绑定：

1. **添加新的 Go 方法**
   ```go
   func (a *AppService) NewMethod() string {
       return "hello"
   }
   ```

2. **修改方法签名**
   ```go
   // 从
   func (a *AppService) Method(arg1 string)
   // 改为
   func (a *AppService) Method(arg1 string, arg2 int)
   ```

3. **添加新的数据类型**
   ```go
   type NewType struct {
       Field string
   }
   ```

4. **Wails 版本升级**
   ```bash
   task wails:update
   wails3 generate bindings -clean=true
   ```

---

## 📝 最佳实践

### 1. 定期重新生成绑定

```bash
# 每次添加新功能后
wails3 generate bindings -clean=true

# 或使用 Task
task wails:bindings
```

### 2. 提交前验证

```bash
# 确保绑定是最新的
task wails:bindings

# 测试应用
task dev

# 验证没有控制台错误
```

### 3. 备份旧绑定

在重新生成之前：
```bash
cp -r frontend/wailsjs frontend/wailsjs.backup
```

如果新绑定有问题，可以恢复：
```bash
rm -rf frontend/wailsjs
mv frontend/wailsjs.backup frontend/wailsjs
```

### 4. 文档化变更

在提交消息中说明：
```
feat: add SearchInFiles method

- Added SearchInFiles to AppService
- Regenerated Wails bindings
- Updated SearchPanel component
```

---

## 🔗 相关资源

- [Wails v3 文档](https://v3.wails.io/)
- [绑定生成指南](https://v3.wails.io/docs/reference/binding-generation)
- [Runtime API](https://v3.wails.io/docs/reference/runtime-api)

---

## 📊 修复统计

| 项目 | 数量 |
|------|------|
| 修改文件 | 2 (App.js, App.d.ts) |
| 新增方法 | 1 (SearchInFiles) |
| 总方法数 | 36 |
| 数据模型 | 8 |

---

## ✅ 验证清单

- [x] 重新生成绑定
- [x] 复制到新位置
- [x] 创建兼容层
- [x] 验证 SearchInFiles 存在
- [x] 测试打开文件夹功能
- [x] 无控制台错误
- [x] 文档化修复过程

---

**最后更新**: 2026-04-15 23:35  
**维护者**: Hao-Code Team
