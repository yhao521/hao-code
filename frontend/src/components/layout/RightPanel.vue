<template>
  <div class="right-panel">
    <!-- 右侧面板标签页 -->
    <div class="panel-tabs">
      <div
        v-for="tab in layoutStore.rightPanelTabs"
        :key="tab.id"
        class="panel-tab"
        :class="{ active: layoutStore.activeRightPanel === tab.id }"
        @click="layoutStore.setActiveRightPanel(tab.id)"
      >
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-label">{{ tab.label }}</span>
      </div>
    </div>
    
    <!-- 面板内容 -->
    <div class="panel-content">
      <!-- 属性面板 -->
      <div v-if="layoutStore.activeRightPanel === 'properties'" class="panel-section">
        <div class="section-header">
          <h3>📋 文件属性</h3>
        </div>
        <div class="section-body">
          <div v-if="activeFile" class="property-list">
            <div class="property-item">
              <span class="property-label">文件名:</span>
              <span class="property-value">{{ activeFile.name }}</span>
            </div>
            <div class="property-item">
              <span class="property-label">路径:</span>
              <span class="property-value">{{ activeFile.path }}</span>
            </div>
            <div class="property-item">
              <span class="property-label">语言:</span>
              <span class="property-value">{{ activeFile.language || 'plaintext' }}</span>
            </div>
            <div class="property-item">
              <span class="property-label">大小:</span>
              <span class="property-value">{{ fileSize }}</span>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>没有打开的文件</p>
          </div>
        </div>
      </div>
      
      <!-- 预览面板 -->
      <div v-else-if="layoutStore.activeRightPanel === 'preview'" class="panel-section">
        <div class="section-header">
          <h3>👁️ 实时预览</h3>
        </div>
        <div class="section-body">
          <div class="empty-state">
            <p>预览功能开发中...</p>
          </div>
        </div>
      </div>
      
      <!-- 大纲面板 -->
      <div v-else-if="layoutStore.activeRightPanel === 'outline'" class="panel-section">
        <div class="section-header">
          <h3>📑 文档大纲</h3>
        </div>
        <div class="section-body">
          <div class="empty-state">
            <p>大纲功能开发中...</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useLayoutStore } from '@/stores/layout'
import { useEditorStore } from '@/stores/editor'

const layoutStore = useLayoutStore()
const editorStore = useEditorStore()

const activeFile = computed(() => {
  if (!editorStore.activeTab) return null
  return {
    name: editorStore.activeTab.name,
    path: editorStore.activeTab.path,
    language: editorStore.activeTab.language
  }
})

const fileSize = computed(() => {
  if (!activeFile.value?.path) return 'N/A'
  const content = editorStore.activeTab?.content || ''
  const bytes = new Blob([content]).size
  if (bytes < 1024) return `${bytes} B`
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(2)} KB`
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`
})
</script>

<style scoped>
.right-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #252526;
  border-left: 1px solid #3E3E42;
}

.panel-tabs {
  display: flex;
  flex-direction: column;
  background-color: #252526;
  border-bottom: 1px solid #3E3E42;
}

.panel-tab {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 8px 12px;
  cursor: pointer;
  color: #969696;
  transition: all 0.2s;
  border-left: 2px solid transparent;
}

.panel-tab:hover {
  background-color: #2A2D2E;
  color: #CCCCCC;
}

.panel-tab.active {
  background-color: #1E1E1E;
  color: #FFFFFF;
  border-left-color: #007ACC;
}

.tab-icon {
  font-size: 14px;
}

.tab-label {
  font-size: 12px;
}

.panel-content {
  flex: 1;
  overflow-y: auto;
  background-color: #1E1E1E;
}

.panel-section {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.section-header {
  padding: 12px;
  border-bottom: 1px solid #3E3E42;
  background-color: #252526;
}

.section-header h3 {
  margin: 0;
  font-size: 13px;
  font-weight: 600;
  color: #CCCCCC;
}

.section-body {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
}

.property-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.property-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.property-label {
  font-size: 11px;
  color: #858585;
  font-weight: 500;
}

.property-value {
  font-size: 12px;
  color: #CCCCCC;
  word-break: break-all;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #858585;
  font-size: 13px;
}
</style>
