<template>
  <div class="file-explorer">
    <div class="explorer-header">
      <span class="explorer-title">{{ workspaceName }}</span>
      <NButton text circle size="tiny" @click="refreshFiles">
        <template #icon>
          <NIcon><RefreshOutline /></NIcon>
        </template>
      </NButton>
    </div>
    
    <NSpin :show="loading">
      <NTree
        v-if="treeData.length > 0"
        :data="treeData"
        :default-expand-all="false"
        block-line
        selectable
        key-field="path"
        label-field="name"
        children-field="children"
        @update:selected-keys="handleFileSelect"
      />
      <div v-else class="empty-tree">
        <NEmpty description="未找到文件" />
      </div>
    </NSpin>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { NButton, NIcon, NTree, NSpin, NEmpty } from 'naive-ui'
import type { TreeOption } from 'naive-ui'
import { RefreshOutline } from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'
import { GetProjectRoot, ListDir, ReadFile } from '@wails/go/backend/App'

const editorStore = useEditorStore()
const loading = ref(false)
const treeData = ref<TreeOption[]>([])
const workspaceName = ref('Hao-Code')
const projectRoot = ref('')

async function loadFiles() {
  loading.value = true
  try {
    // 获取项目根目录
    const root = await GetProjectRoot()
    projectRoot.value = root
    
    // 读取目录内容
    const files = await ListDir(root)
    treeData.value = convertToTree(files)
  } catch (error) {
    console.error('Failed to load files:', error)
  } finally {
    loading.value = false
  }
}

function convertToTree(files: any[]): TreeOption[] {
  return files.map(file => ({
    key: file.path,
    name: file.name,
    isLeaf: !file.isDir,
    ...(file.isDir ? { 
      children: [] // 懒加载子目录
    } : {})
  }))
}

async function handleFileSelect(keys: string[]) {
  if (keys.length === 0) return
  
  const filePath = keys[0]
  
  try {
    loading.value = true
    const content = await ReadFile(filePath)
    editorStore.openFile(filePath, content)
  } catch (error) {
    console.error('Failed to read file:', error)
  } finally {
    loading.value = false
  }
}

async function refreshFiles() {
  await loadFiles()
}

onMounted(() => {
  loadFiles()
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
  padding-left: 8px !important;
}
</style>
