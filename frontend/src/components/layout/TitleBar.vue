<template>
  <div class="titlebar">
    <div class="titlebar-left">
      <span class="app-title">Hao-Code Editor</span>
    </div>
    <div class="titlebar-center">
      <span class="workspace-name">{{ workspaceName }}</span>
    </div>
    <div class="titlebar-right">
      <NButton quaternary circle size="small" @click="minimizeWindow">
        <template #icon>
          <NIcon><RemoveOutline /></NIcon>
        </template>
      </NButton>
      <NButton quaternary circle size="small" @click="maximizeWindow">
        <template #icon>
          <NIcon><SquareOutline /></NIcon>
        </template>
      </NButton>
      <NButton quaternary circle size="small" @click="closeWindow" class="close-btn">
        <template #icon>
          <NIcon><CloseOutline /></NIcon>
        </template>
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { NButton, NIcon } from 'naive-ui'
import { RemoveOutline, SquareOutline, CloseOutline } from '@vicons/ionicons5'
import { useEditorStore } from '@/stores/editor'

const editorStore = useEditorStore()

const workspaceName = computed(() => 'Hao-Code')

function minimizeWindow() {
  // TODO: 调用 Wails API 最小化窗口
  console.log('Minimize window')
}

function maximizeWindow() {
  // TODO: 调用 Wails API 最大化窗口
  console.log('Maximize window')
}

function closeWindow() {
  // TODO: 调用 Wails API 关闭窗口
  console.log('Close window')
}
</script>

<style scoped>
.titlebar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 30px;
  background-color: #323233;
  color: #CCCCCC;
  user-select: none;
  -webkit-app-region: drag;
  border-bottom: 1px solid #252526;
}

.titlebar-left,
.titlebar-center,
.titlebar-right {
  display: flex;
  align-items: center;
  gap: 8px;
}

.app-title {
  font-size: 12px;
  font-weight: 500;
  padding-left: 12px;
}

.workspace-name {
  font-size: 12px;
  opacity: 0.9;
}

.titlebar-right {
  -webkit-app-region: no-drag;
  padding-right: 8px;
}

:deep(.n-button) {
  width: 46px !important;
  height: 30px !important;
  border-radius: 0 !important;
}

:deep(.n-button:hover) {
  background-color: rgba(255, 255, 255, 0.1) !important;
}

:deep(.close-btn:hover) {
  background-color: #E81123 !important;
  color: white !important;
}

:deep(.n-button .n-icon) {
  font-size: 16px;
}
</style>
