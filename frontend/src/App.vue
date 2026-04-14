<script lang="ts" setup>
import { ref } from 'vue'
import { 
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  NLayout,
  NLayoutSider,
  NLayoutContent,
  darkTheme,
  zhCN,
  dateZhCN
} from 'naive-ui'

// 导入组件（稍后创建）
import TitleBar from './components/layout/TitleBar.vue'
import SideBar from './components/layout/SideBar.vue'
import EditorArea from './components/editor/EditorArea.vue'
import StatusBar from './components/layout/StatusBar.vue'

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
            <TitleBar />
            <NLayout has-sider class="main-layout">
              <!-- 侧边栏 -->
              <NLayoutSider
                bordered
                collapse-mode="width"
                :collapsed-width="60"
                :width="240"
                :native-scrollbar="false"
              >
                <SideBar />
              </NLayoutSider>
              
              <!-- 主内容区 -->
              <NLayoutContent>
                <EditorArea />
              </NLayoutContent>
            </NLayout>
            
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

.main-layout {
  flex: 1;
  overflow: hidden;
}
</style>
