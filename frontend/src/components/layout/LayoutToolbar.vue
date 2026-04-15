<template>
  <div class="layout-toolbar">
    <div class="toolbar-left">
      <!-- 布局切换按钮组 -->
      <div class="layout-buttons">
        <NButton
          :type="layoutStore.layoutMode === 'default' ? 'primary' : 'default'"
          size="tiny"
          quaternary
          @click="layoutStore.setLayoutMode('default')"
          title="默认布局"
        >
          <template #icon>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <rect x="1" y="1" width="4" height="14" rx="1" opacity="0.5"/>
              <rect x="6" y="1" width="9" height="14" rx="1"/>
            </svg>
          </template>
        </NButton>
        
        <NButton
          :type="layoutStore.layoutMode === 'side-by-side' ? 'primary' : 'default'"
          size="tiny"
          quaternary
          @click="layoutStore.setLayoutMode('side-by-side')"
          title="左右分栏"
        >
          <template #icon>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <rect x="1" y="1" width="4" height="14" rx="1" opacity="0.5"/>
              <rect x="6" y="1" width="4" height="14" rx="1"/>
              <rect x="11" y="1" width="4" height="14" rx="1" opacity="0.5"/>
            </svg>
          </template>
        </NButton>
        
        <NButton
          :type="layoutStore.layoutMode === 'bottom-focused' ? 'primary' : 'default'"
          size="tiny"
          quaternary
          @click="layoutStore.setLayoutMode('bottom-focused')"
          title="底部面板"
        >
          <template #icon>
            <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
              <rect x="1" y="1" width="14" height="6" rx="1"/>
              <rect x="1" y="9" width="14" height="6" rx="1" opacity="0.5"/>
            </svg>
          </template>
        </NButton>
      </div>
    </div>
    
    <div class="toolbar-center">
      <!-- 工作区名称 -->
      <span class="workspace-name">{{ workspaceName }}</span>
    </div>
    
    <div class="toolbar-right">
      <!-- 面板切换按钮 -->
      <NButton
        :type="layoutStore.bottomPanelVisible ? 'primary' : 'default'"
        size="tiny"
        quaternary
        @click="layoutStore.toggleBottomPanel"
        title="切换终端面板"
      >
        <template #icon>
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <path d="M2 4h12v8H2V4zm1 1v6h10V5H3zm1 2l2 1.5L4 10V7zm3 0h4v1H7V7z"/>
          </svg>
        </template>
      </NButton>
      
      <NButton
        :type="layoutStore.rightPanelVisible ? 'primary' : 'default'"
        size="tiny"
        quaternary
        @click="layoutStore.toggleRightPanel"
        title="切换右侧面板"
      >
        <template #icon>
          <svg width="16" height="16" viewBox="0 0 16 16" fill="currentColor">
            <rect x="1" y="1" width="10" height="14" rx="1" opacity="0.5"/>
            <rect x="12" y="1" width="3" height="14" rx="1"/>
          </svg>
        </template>
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NButton } from 'naive-ui'
import { useLayoutStore } from '@/stores/layout'
import { useEditorStore } from '@/stores/editor'

const layoutStore = useLayoutStore()
const editorStore = useEditorStore()

const workspaceName = computed(() => {
  if (!editorStore.workspace) return 'hao-code'
  return editorStore.workspace.path.split('/').pop() || 'hao-code'
})
</script>

<style scoped>
.layout-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 35px;
  background-color: #252526;
  border-bottom: 1px solid #3E3E42;
  padding: 0 8px;
  user-select: none;
}

.toolbar-left,
.toolbar-right {
  display: flex;
  align-items: center;
  gap: 4px;
}

.toolbar-center {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.workspace-name {
  font-size: 13px;
  color: #CCCCCC;
  font-weight: 500;
}

.layout-buttons {
  display: flex;
  align-items: center;
  gap: 2px;
}

:deep(.n-button) {
  padding: 4px !important;
  height: 28px !important;
  min-width: 28px !important;
}

:deep(.n-button .n-button__icon) {
  font-size: 16px;
}

:deep(.n-button--primary-type) {
  background-color: #0E639C !important;
  color: #FFFFFF !important;
}

:deep(.n-button--primary-type:hover) {
  background-color: #1177BB !important;
}
</style>
