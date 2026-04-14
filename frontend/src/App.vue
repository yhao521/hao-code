<script lang="ts" setup>
import { ref } from 'vue'
import { 
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  darkTheme,
  zhCN,
  dateZhCN
} from 'naive-ui'

// 导入组件
import TitleBar from './components/layout/TitleBar.vue'
import SideBar from './components/layout/SideBar.vue'
import EditorArea from './components/editor/EditorArea.vue'
import StatusBar from './components/layout/StatusBar.vue'
import ResizableSplit from './components/layout/ResizableSplit.vue'

// 主题配置 - VSCode 风格深色主题
const theme = ref(darkTheme)
const themeOverrides = {
  common: {
    primaryColor: '#0E639C',
    primaryColorHover: '#1177BB',
    bodyColor: '#1E1E1E',
    cardColor: '#252526',
    borderColor: '#3E3E42',
    textColor: '#CCCCCC'
  },
  Menu: {
    itemColor: '#CCCCCC',
    itemColorActive: '#FFFFFF',
    itemColorHover: '#2A2D2E'
  },
  Tree: {
    nodeColor: '#CCCCCC',
    nodeColorHover: '#2A2D2E'
  }
}
</script>

<template>
  <NConfigProvider
    :theme="theme"
    :theme-overrides="themeOverrides"
    :locale="zhCN"
    :date-locale="dateZhCN"
  >
    <NMessageProvider>
      <NNotificationProvider>
        <NDialogProvider>
          <div class="app-container">
            <!-- 标题栏（所有平台显示） -->
            <TitleBar />
            
            <!-- 主内容区 - 可拖拽分割 -->
            <div class="main-content">
              <ResizableSplit 
                :min="180"
                :max="500"
                :horizontal="true"
                class="main-split"
              >
                <template #1>
                  <div class="sidebar-container">
                    <SideBar />
                  </div>
                </template>
                <template #2>
                  <div class="editor-container">
                    <EditorArea />
                  </div>
                </template>
              </ResizableSplit>
            </div>
            
            <!-- 状态栏 -->
            <StatusBar />
          </div>
        </NDialogProvider>
      </NNotificationProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif;
  overflow: hidden;
}

#app {
  width: 100vw;
  height: 100vh;
}

.app-container {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background-color: #1E1E1E;
}

.main-content {
  flex: 1;
  overflow: hidden;
  background-color: #1E1E1E;
}

.main-split {
  width: 100%;
  height: 100%;
}

.sidebar-container {
  height: 100%;
  background-color: #252526;
  border-right: 1px solid #3E3E42;
  overflow: hidden;
}

.editor-container {
  height: 100%;
  background-color: #1E1E1E;
  overflow: hidden;
}
</style>