<template>
  <div class="file-explorer">
    <div class="explorer-header">
      <span class="explorer-title">资源管理器</span>
      <NButton text circle size="tiny" @click="refreshFiles">
        <template #icon>
          <NIcon><RefreshOutline /></NIcon>
        </template>
      </NButton>
    </div>
    
    <NTree
      :data="treeData"
      :default-expand-all="true"
      block-line
      selectable
      @update:selected-keys="handleFileSelect"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { NButton, NIcon, NTree } from 'naive-ui'
import type { TreeOption } from 'naive-ui'
import { RefreshOutline } from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'

const editorStore = useEditorStore()

// 模拟文件树数据（实际应该从后端获取）
const treeData = ref<TreeOption[]>([
  {
    key: 'src',
    label: 'src',
    children: [
      {
        key: 'frontend',
        label: 'frontend',
        children: [
          { key: 'src/main.ts', label: 'main.ts' },
          { key: 'src/App.vue', label: 'App.vue' },
          { key: 'src/components/layout/TitleBar.vue', label: 'TitleBar.vue' },
          { key: 'src/stores/editor.ts', label: 'editor.ts' }
        ]
      },
      {
        key: 'backend',
        label: 'backend',
        children: [
          { key: 'main.go', label: 'main.go' },
          { key: 'app.go', label: 'app.go' }
        ]
      }
    ]
  },
  {
    key: 'package.json',
    label: 'package.json'
  },
  {
    key: 'go.mod',
    label: 'go.mod'
  },
  {
    key: 'README.md',
    label: 'README.md'
  }
])

function handleFileSelect(keys: string[]) {
  if (keys.length > 0) {
    const filePath = keys[0]
    // TODO: 读取文件内容
    editorStore.openFile(filePath, '// File content')
  }
}

function refreshFiles() {
  // TODO: 刷新文件树
  console.log('Refresh files')
}

onMounted(() => {
  // TODO: 加载文件树
})
</script>

<style scoped>
.file-explorer {
  padding: 8px;
}

.explorer-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 8px;
  padding-bottom: 8px;
  border-bottom: 1px solid #3E3E42;
}

.explorer-title {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  color: #BBBBBB;
}
</style>
