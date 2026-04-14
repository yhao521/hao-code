<template>
  <div class="sidebar">
    <NMenu
      :options="menuOptions"
      :value="activeView"
      @update:value="handleViewChange"
    />
    
    <div class="sidebar-content">
      <FileExplorer v-if="activeView === 'explorer'" />
      <SearchPanel v-else-if="activeView === 'search'" />
      <GitPanel v-else-if="activeView === 'git'" />
      <ExtensionsPanel v-else-if="activeView === 'extensions'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { h, ref } from 'vue'
import { NMenu } from 'naive-ui'
import type { MenuOption } from 'naive-ui'
import { 
  FolderOutline, 
  SearchOutline, 
  GitBranchOutline,
  CubeOutline
} from '@vicons/ionicons5'
import FileExplorer from './FileExplorer.vue'
import SearchPanel from './SearchPanel.vue'
import GitPanel from './GitPanel.vue'
import ExtensionsPanel from './ExtensionsPanel.vue'

const activeView = ref('explorer')

const menuOptions: MenuOption[] = [
  {
    label: '资源管理器',
    key: 'explorer',
    icon: () => h('div', null, [h(FolderOutline)])
  },
  {
    label: '搜索',
    key: 'search',
    icon: () => h('div', null, [h(SearchOutline)])
  },
  {
    label: '源代码管理',
    key: 'git',
    icon: () => h('div', null, [h(GitBranchOutline)])
  },
  {
    label: '扩展',
    key: 'extensions',
    icon: () => h('div', null, [h(CubeOutline)])
  }
]

function handleViewChange(key: string) {
  activeView.value = key
}
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #252526;
}

.sidebar-content {
  flex: 1;
  overflow-y: auto;
}
</style>
