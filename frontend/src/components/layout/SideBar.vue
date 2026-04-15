<template>
  <div class="sidebar">
    <!-- 活动栏（垂直图标） -->
    <div class="activity-bar">
      <div 
        v-for="item in activityItems" 
        :key="item.id"
        class="activity-item"
        :class="{ active: activeView === item.id }"
        @click="handleViewChange(item.id)"
        :title="item.title"
      >
        <NIcon :component="item.icon" :size="24" />
      </div>
    </div>
    
    <!-- 侧边栏内容 -->
    <div class="sidebar-content">
      <FileExplorer v-if="activeView === 'explorer'" />
      <SearchPanel v-else-if="activeView === 'search'" />
      <GitPanel v-else-if="activeView === 'git'" />
      <ExtensionsPanel v-else-if="activeView === 'extensions'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NIcon } from 'naive-ui'
import type { Component } from 'vue'
import { 
  FolderOutline, 
  SearchOutline, 
  GitBranchOutline,
  CubeOutline
} from '@vicons/ionicons5'
import FileExplorer from './FileExplorer.vue'
import SearchPanel from '../SearchPanel.vue'
import GitPanel from './GitPanel.vue'
import ExtensionsPanel from './ExtensionsPanel.vue'

interface ActivityItem {
  id: string
  title: string
  icon: Component
}

const activeView = ref('explorer')

const activityItems: ActivityItem[] = [
  {
    id: 'explorer',
    title: '资源管理器',
    icon: FolderOutline
  },
  {
    id: 'search',
    title: '搜索',
    icon: SearchOutline
  },
  {
    id: 'git',
    title: '源代码管理',
    icon: GitBranchOutline
  },
  {
    id: 'extensions',
    title: '扩展',
    icon: CubeOutline
  }
]

function handleViewChange(id: string) {
  activeView.value = activeView.value === id ? '' : id
}
</script>

<style scoped>
.sidebar {
  display: flex;
  height: 100%;
  background-color: #2C2C2C;
}

/* 活动栏样式 */
.activity-bar {
  width: 48px;
  min-width: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 4px 0;
  background-color: #333333;
  border-right: 1px solid #3E3E42;
}

.activity-item {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #858585;
  transition: all 0.2s;
  position: relative;
}

.activity-item:hover {
  color: #CCCCCC;
  background-color: #2A2D2E;
}

.activity-item.active {
  color: #FFFFFF;
}

.activity-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: #FFFFFF;
}

/* 侧边栏内容 */
.sidebar-content {
  flex: 1;
  width: 240px;
  min-width: 170px;
  overflow-y: auto;
  background-color: #252526;
}
</style>
