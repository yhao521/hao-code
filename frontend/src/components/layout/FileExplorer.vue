<template>
  <div class="file-explorer" @contextmenu.prevent>
    <div class="explorer-header">
      <span class="explorer-title">{{ workspaceName }}</span>
      <div class="header-actions">
        <NButton 
          v-if="hasWorkspace" 
          text 
          circle 
          size="tiny" 
          @click="handleToolbarNewFile" 
          title="新建文件"
        >
          <template #icon>
            <NIcon><AddOutline /></NIcon>
          </template>
        </NButton>
        <NButton 
          v-if="hasWorkspace" 
          text 
          circle 
          size="tiny" 
          @click="handleToolbarNewFolder" 
          title="新建文件夹"
        >
          <template #icon>
            <NIcon><FolderOpenOutline /></NIcon>
          </template>
        </NButton>
        <NButton 
          v-if="hasWorkspace" 
          text 
          circle 
          size="tiny" 
          @click="refreshFiles" 
          title="刷新"
        >
          <template #icon>
            <NIcon><RefreshOutline /></NIcon>
          </template>
        </NButton>
      </div>
    </div>
    
    <!-- 未打开文件夹时显示提示 -->
    <div v-if="!hasWorkspace" class="no-workspace">
      <NEmpty description="未打开文件夹">
        <template #extra>
          <NButton size="small" type="primary" @click="openFolderDialog">
            <template #icon>
              <NIcon><FolderOpenOutline /></NIcon>
            </template>
            打开文件夹
          </NButton>
        </template>
      </NEmpty>
    </div>
    
    <!-- 有工作区时显示文件树 -->
    <NSpin v-else :show="loading">
      <NTree
        v-if="treeData.length > 0"
        :data="treeData"
        :default-expand-all="false"
        block-line
        selectable
        key-field="path"
        children-field="children"
        show-line
        :render-label="renderTreeNode"
        @update:selected-keys="handleFileSelect"
        @update:expanded-keys="handleExpand"
      />
      <div v-else class="empty-tree">
        <NEmpty description="文件夹为空" />
      </div>
    </NSpin>

    <!-- 新建文件对话框 -->
    <NModal v-model:show="showNewFileModal" preset="dialog" title="新建文件">
      <NForm :model="newFileForm" label-placement="left" label-width="80">
        <NFormItem label="文件名" path="name">
          <NInput 
            v-model:value="newFileForm.name" 
            placeholder="例如: index.ts"
            @keyup.enter="createNewFile"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showNewFileModal = false">取消</NButton>
          <NButton type="primary" @click="createNewFile">创建</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 新建文件夹对话框 -->
    <NModal v-model:show="showNewFolderModal" preset="dialog" title="新建文件夹">
      <NForm :model="newFolderForm" label-placement="left" label-width="80">
        <NFormItem label="文件夹名" path="name">
          <NInput 
            v-model:value="newFolderForm.name" 
            placeholder="例如: components"
            @keyup.enter="createNewFolder"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showNewFolderModal = false">取消</NButton>
          <NButton type="primary" @click="createNewFolder">创建</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 重命名对话框 -->
    <NModal v-model:show="showRenameModal" preset="dialog" title="重命名">
      <NForm :model="renameForm" label-placement="left" label-width="80">
        <NFormItem label="新名称" path="name">
          <NInput 
            v-model:value="renameForm.name" 
            placeholder="输入新名称"
            @keyup.enter="renameItem"
          />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showRenameModal = false">取消</NButton>
          <NButton type="primary" @click="renameItem">确定</NButton>
        </NSpace>
      </template>
    </NModal>

    <!-- 删除确认对话框 -->
    <NModal v-model:show="showDeleteModal" preset="dialog" title="确认删除">
      <p>确定要删除 "{{ deleteItemName }}" 吗？此操作不可撤销。</p>
      <template #footer>
        <NSpace>
          <NButton @click="showDeleteModal = false">取消</NButton>
          <NButton type="error" @click="confirmDelete">删除</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, h, computed } from 'vue'
import { 
  NButton, NIcon, NTree, NSpin, NEmpty, NModal, NForm, NFormItem, 
  NInput, NSpace, NDropdown, useMessage, type TreeOption, type DropdownOption 
} from 'naive-ui'
import { 
  RefreshOutline, 
  AddOutline, 
  FolderOpenOutline, 
  DocumentTextOutline,
  CodeSlashOutline,
  LogoVue,
  LogoJavascript,
  LogoCss3,
  LogoHtml5,
  FileTrayOutline,
  CreateOutline,
  TrashOutline,
  CopyOutline,
  PencilOutline,
  LogoMarkdown
} from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'
import { 
  GetProjectRoot, 
  ListDir, 
  ReadFile,
  CreateFile,
  CreateDirectory,
  DeleteFileOrDirectory,
  RenameFileOrDirectory,
  CopyFileOrDirectory,
  OpenFolderDialog
} from '@wails/go/backend/App'

// 扩展 TreeOption 类型以支持自定义属性
interface ExtendedTreeOption extends TreeOption {
  isDir?: boolean
  extension?: string
}

const editorStore = useEditorStore()
const message = useMessage()
const loading = ref(false)
const treeData = ref<ExtendedTreeOption[]>([])
const workspaceName = ref('Hao-Code')
const projectRoot = ref('')
const hasWorkspace = computed(() => !!editorStore.workspace)

// 模态框状态
const showNewFileModal = ref(false)
const showNewFolderModal = ref(false)
const showRenameModal = ref(false)
const showDeleteModal = ref(false)

// 表单数据
const newFileForm = ref({ name: '' })
const newFolderForm = ref({ name: '' })
const renameForm = ref({ name: '' })
const deleteItemName = ref('')

// 当前操作项
const currentTargetPath = ref('')
const currentParentPath = ref('')

// 获取文件扩展名
function getExtension(filename: string): string {
  const parts = filename.split('.')
  return parts.length > 1 ? parts.pop()?.toLowerCase() || '' : ''
}

// 根据文件类型获取图标
function getFileIcon(option: ExtendedTreeOption) {
  if (option.isDir) {
    return FolderOpenOutline
  }
  
  const ext = option.extension || getExtension(option.name as string)
  
  const iconMap: Record<string, any> = {
    'vue': LogoVue,
    'js': LogoJavascript,
    'jsx': LogoJavascript,
    'ts': CodeSlashOutline,
    'tsx': CodeSlashOutline,
    'css': LogoCss3,
    'html': LogoHtml5,
    'htm': LogoHtml5,
    'md': LogoMarkdown,
    'markdown': LogoMarkdown,
    'json': DocumentTextOutline,
    'yaml': DocumentTextOutline,
    'yml': DocumentTextOutline,
    'xml': DocumentTextOutline,
    'go': CodeSlashOutline,
    'py': CodeSlashOutline,
    'java': CodeSlashOutline,
    'cpp': CodeSlashOutline,
    'c': CodeSlashOutline,
    'h': CodeSlashOutline,
    'rs': CodeSlashOutline,
    'rb': CodeSlashOutline,
    'php': CodeSlashOutline,
    'sh': CodeSlashOutline,
    'bash': CodeSlashOutline,
    'txt': DocumentTextOutline,
    'log': DocumentTextOutline,
    'env': DocumentTextOutline,
    'gitignore': DocumentTextOutline,
    'dockerfile': DocumentTextOutline,
    'mod': DocumentTextOutline,
    'sum': DocumentTextOutline,
  }
  
  return iconMap[ext] || FileTrayOutline
}

// 根据文件类型获取图标颜色
function getFileIconColor(option: ExtendedTreeOption): string {
  if (option.isDir) {
    return '#D7BA7D' // 文件夹黄色
  }
  
  const ext = option.extension || getExtension(option.name as string)
  
  const colorMap: Record<string, string> = {
    'vue': '#42b883',
    'js': '#f7df1e',
    'jsx': '#61dafb',
    'ts': '#3178c6',
    'tsx': '#3178c6',
    'css': '#1572b6',
    'html': '#e34c26',
    'md': '#519aba',
    'go': '#00add8',
    'py': '#3776ab',
    'java': '#b07219',
    'json': '#519aba',
  }
  
  return colorMap[ext] || '#CCCCCC'
}

// 自定义树节点渲染
function renderTreeNode(data: any) {
  const option = data.option as ExtendedTreeOption
  
  return h('div', { class: 'tree-node-content' }, [
    h(NIcon, {
      component: getFileIcon(option),
      color: getFileIconColor(option),
      size: 16,
      style: { marginRight: '6px', flexShrink: '0' }
    }),
    h('span', { class: 'tree-node-label' }, [String(option.name)])
  ])
}

async function loadFiles(rootPath?: string) {
  loading.value = true
  try {
    // 如果没有指定路径，使用当前工作区
    const root = rootPath || await GetProjectRoot()
    projectRoot.value = root
    
    // 读取目录内容
    const files = await ListDir(root)
    treeData.value = convertToTree(files)
  } catch (error) {
    console.error('Failed to load files:', error)
    message.error('加载文件失败')
  } finally {
    loading.value = false
  }
}

function convertToTree(files: any[]): ExtendedTreeOption[] {
  // 排序：文件夹优先，然后按名称字母顺序
  const sortedFiles = files.sort((a, b) => {
    // 文件夹排在前面
    if (a.isDir && !b.isDir) return -1
    if (!a.isDir && b.isDir) return 1
    // 同类型按名称排序
    return a.name.localeCompare(b.name)
  })
  
  return sortedFiles.map(file => ({
    key: file.path,
    path: file.path,
    name: file.name,
    isLeaf: !file.isDir,
    disabled: false,
    extension: file.isDir ? undefined : getExtension(file.name),
    ...(file.isDir ? { 
      children: [],
      isDir: true
    } : {
      isDir: false
    })
  }))
}

async function handleFileSelect(keys: string[]) {
  if (keys.length === 0) return
  
  const filePath = keys[0]
  
  // 查找选中的节点
  const selectedNode = findNodeByKey(treeData.value, filePath) as ExtendedTreeOption
  
  // 如果是文件夹，不做任何操作（NTree 会自动处理展开/折叠）
  if (selectedNode?.isDir) {
    console.log('[FileExplorer] Clicked folder:', filePath)
    return
  }
  
  // 只有文件才打开到编辑器
  try {
    loading.value = true
    console.log('[FileExplorer] Reading file:', filePath)
    
    const content = await ReadFile(filePath)
    
    console.log('[FileExplorer] File content loaded, length:', content.length)
    
    editorStore.openFile(filePath, content)
    message.success(`已打开: ${filePath.split('/').pop()}`)
  } catch (error) {
    console.error('[FileExplorer] Failed to read file:', error)
    const errorMsg = error instanceof Error ? error.message : String(error)
    message.error(`读取文件失败: ${errorMsg}`)
  } finally {
    loading.value = false
  }
}

async function handleExpand(keys: string[]) {
  // 懒加载子目录
  for (const key of keys) {
    const node = findNodeByKey(treeData.value, key) as ExtendedTreeOption
    if (node && node.isDir && (!node.children || node.children.length === 0)) {
      try {
        const files = await ListDir(key)
        node.children = convertToTree(files)
      } catch (error) {
        console.error('Failed to load directory:', error)
        message.error('加载目录失败')
      }
    }
  }
}

function findNodeByKey(nodes: ExtendedTreeOption[], key: string): ExtendedTreeOption | null {
  for (const node of nodes) {
    if (node.key === key) return node
    if (node.children) {
      const found = findNodeByKey(node.children as ExtendedTreeOption[], key)
      if (found) return found
    }
  }
  return null
}

async function refreshFiles() {
  treeData.value = []
  await loadFiles()
  message.success('刷新成功')
}

// 创建新文件
async function createNewFile() {
  if (!newFileForm.value.name) {
    message.warning('请输入文件名')
    return
  }
  
  try {
    const filePath = pathJoin(currentParentPath.value, newFileForm.value.name)
    await CreateFile(filePath)
    message.success('文件创建成功')
    showNewFileModal.value = false
    newFileForm.value.name = ''
    await refreshFiles()
  } catch (error) {
    console.error('Failed to create file:', error)
    message.error(`创建文件失败: ${error}`)
  }
}

// 创建新文件夹
async function createNewFolder() {
  if (!newFolderForm.value.name) {
    message.warning('请输入文件夹名')
    return
  }
  
  try {
    const dirPath = pathJoin(currentParentPath.value, newFolderForm.value.name)
    await CreateDirectory(dirPath)
    message.success('文件夹创建成功')
    showNewFolderModal.value = false
    newFolderForm.value.name = ''
    await refreshFiles()
  } catch (error) {
    console.error('Failed to create folder:', error)
    message.error(`创建文件夹失败: ${error}`)
  }
}

// 重命名
async function renameItem() {
  if (!renameForm.value.name) {
    message.warning('请输入新名称')
    return
  }
  
  try {
    const oldPath = currentTargetPath.value
    const newPath = pathJoin(pathDirname(oldPath), renameForm.value.name)
    await RenameFileOrDirectory(oldPath, newPath)
    message.success('重命名成功')
    showRenameModal.value = false
    renameForm.value.name = ''
    await refreshFiles()
  } catch (error) {
    console.error('Failed to rename:', error)
    message.error(`重命名失败: ${error}`)
  }
}

// 复制
async function copyItem(sourcePath: string) {
  try {
    const ext = pathExtname(sourcePath)
    const baseName = pathBasename(sourcePath, ext)
    const dirName = pathDirname(sourcePath)
    const targetPath = pathJoin(dirName, `${baseName}-copy${ext}`)
    
    await CopyFileOrDirectory(sourcePath, targetPath)
    message.success('复制成功')
    await refreshFiles()
  } catch (error) {
    console.error('Failed to copy:', error)
    message.error(`复制失败: ${error}`)
  }
}

// 确认删除
async function confirmDelete() {
  try {
    await DeleteFileOrDirectory(currentTargetPath.value)
    message.success('删除成功')
    showDeleteModal.value = false
    await refreshFiles()
  } catch (error) {
    console.error('Failed to delete:', error)
    message.error(`删除失败: ${error}`)
  }
}

// 工具栏操作
function handleToolbarNewFile() {
  currentParentPath.value = projectRoot.value
  newFileForm.value.name = ''
  showNewFileModal.value = true
}

function handleToolbarNewFolder() {
  currentParentPath.value = projectRoot.value
  newFolderForm.value.name = ''
  showNewFolderModal.value = true
}

// 打开文件夹对话框
async function openFolderDialog() {
  try {
    message.loading('正在打开文件夹选择对话框...', { duration: 0 })
    
    const selectedPath = await OpenFolderDialog()
    
    message.destroyAll()
    
    if (!selectedPath) {
      message.info('已取消选择')
      return
    }
    
    message.loading('正在加载文件夹...', { duration: 0 })
    
    // 验证文件夹
    try {
      await ListDir(selectedPath)
    } catch (error) {
      message.destroyAll()
      message.error('无法访问该文件夹')
      return
    }
    
    // 设置工作区并加载文件
    editorStore.setWorkspace(selectedPath)
    await loadFiles(selectedPath)
    
    message.destroyAll()
    message.success(`已打开: ${selectedPath.split('/').pop()}`)
    
  } catch (error) {
    message.destroyAll()
    console.error('Failed to open folder:', error)
    const errorMsg = error instanceof Error ? error.message : String(error)
    if (!errorMsg.includes('cancelled')) {
      message.error(`打开文件夹失败: ${errorMsg}`)
    } else {
      message.info('已取消选择')
    }
  }
}

// 简单的路径处理函数
function pathJoin(...paths: string[]): string {
  return paths.join('/').replace(/\/+/g, '/')
}

function pathDirname(p: string): string {
  return p.split('/').slice(0, -1).join('/')
}

function pathBasename(p: string, ext?: string): string {
  const base = p.split('/').pop() || ''
  return ext ? base.replace(new RegExp(ext.replace('.', '\\.') + '$'), '') : base
}

function pathExtname(p: string): string {
  const base = p.split('/').pop() || ''
  const match = base.match(/\.[^.]+$/)
  return match ? match[0] : ''
}

onMounted(() => {
  // 如果已经有工作区，加载文件
  if (editorStore.workspace) {
    loadFiles(editorStore.workspace.path)
  } else {
    // 否则加载默认项目根目录
    loadFiles()
  }
})
</script>

<style scoped>
.file-explorer {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #252526;
}

.explorer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 4px 8px;
  min-height: 22px;
  background-color: #252526;
  border-bottom: 1px solid #3E3E42;
  user-select: none;
}

.header-actions {
  display: flex;
  gap: 2px;
}

.explorer-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  color: #BBBBBB;
  letter-spacing: 0.5px;
}

.empty-tree {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #858585;
}

.no-workspace {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 20px;
}

/* 树节点样式优化 */
:deep(.n-tree-node) {
  font-size: 13px;
  color: #CCCCCC;
  padding: 2px 0;
}

:deep(.n-tree-node:hover) {
  background-color: #2A2D2E;
}

:deep(.n-tree-node--selected) {
  background-color: #37373D !important;
}

:deep(.n-tree-node__content) {
  padding-left: 0 !important; /* 移除左侧内边距 */
  padding-right: 0 !important; /* 移除右侧内边距 */
  height: 22px;
  justify-content: flex-start !important; /* 左对齐 */
  text-align: left !important; /* 文本左对齐 */
}

:deep(.n-tree-node__content > .n-tree-node-switcher) {
  width: 16px;
  height: 22px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

:deep(.n-tree-node__content > .n-tree-node-checkbox) {
  margin-right: 4px;
}

:deep(.n-tree-node__content > .n-tree-node-content__text) {
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* 自定义树节点内容 */
.tree-node-content {
  display: flex;
  align-items: center;
  justify-content: flex-start; /* 左对齐 */
  flex: 1;
  overflow: hidden;
  width: 100%;
  min-width: 0;
}

.tree-node-label {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
  text-align: left; /* 文本左对齐 */
  min-width: 0; /* 允许文本收缩 */
}

/* 修复 Naive UI NTree 可能的居中问题 */
:deep(.n-tree .n-tree-node) {
  text-align: left;
}

:deep(.n-tree .n-tree-node-wrapper) {
  text-align: left;
}

:deep(.n-tree-node-content) {
  justify-content: flex-start !important;
  text-align: left !important;
  width: 100%;
}
</style>
