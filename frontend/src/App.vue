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
          <AppContent />
        </NDialogProvider>
      </NNotificationProvider>
    </NMessageProvider>
  </NConfigProvider>
</template>

<script lang="ts" setup>
import { ref } from "vue";
import {
  NConfigProvider,
  NMessageProvider,
  NNotificationProvider,
  NDialogProvider,
  darkTheme,
  zhCN,
  dateZhCN,
} from "naive-ui";

// 导入实际的应用内容组件
import AppContent from "./AppContent.vue";

// 主题配置 - VSCode 风格深色主题
const theme = ref(darkTheme);
const themeOverrides = {
  common: {
    primaryColor: "#0E639C",
    primaryColorHover: "#1177BB",
    bodyColor: "#1E1E1E",
    cardColor: "#252526",
    borderColor: "#3E3E42",
    textColor: "#CCCCCC",
  },
  Menu: {
    itemColor: "#CCCCCC",
    itemColorActive: "#FFFFFF",
    itemColorHover: "#2A2D2E",
  },
  Tree: {
    nodeColor: "#CCCCCC",
    nodeColorHover: "#2A2D2E",
  },
};
</script>

<style>
* {
  margin: 0;
  padding: 0;
  box-sizing: border-box;
}

body {
  font-family:
    -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue",
    Arial, sans-serif;
  overflow: hidden;
}

#app {
  width: 100vw;
  height: 100vh;
}
</style>
