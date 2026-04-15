<template>
  <div class="bottom-panel">
    <!-- 下面板标签页 -->
    <div class="panel-tabs">
      <div
        v-for="tab in layoutStore.bottomPanelTabs"
        :key="tab.id"
        class="panel-tab"
        :class="{ active: layoutStore.activeBottomPanel === tab.id }"
        @click="layoutStore.setActiveBottomPanel(tab.id)"
      >
        <span class="tab-icon">{{ tab.icon }}</span>
        <span class="tab-label">{{ tab.label }}</span>
      </div>
      <div class="panel-actions">
        <NButton
          size="tiny"
          quaternary
          circle
          @click="layoutStore.bottomPanelVisible = false"
          title="关闭面板"
        >
          <template #icon>
            <span style="font-size: 14px">✕</span>
          </template>
        </NButton>
      </div>
    </div>

    <!-- 面板内容 -->
    <div class="panel-content">
      <!-- 终端面板 -->
      <div
        v-if="layoutStore.activeBottomPanel === 'terminal'"
        class="panel-section"
      >
        <Terminal />
      </div>

      <!-- 输出面板 -->
      <div
        v-else-if="layoutStore.activeBottomPanel === 'output'"
        class="panel-section"
      >
        <div class="output-container">
          <div class="output-header">
            <span class="output-title">📄 输出</span>
          </div>
          <div class="output-body">
            <div
              class="output-line"
              v-for="(line, index) in outputLines"
              :key="index"
            >
              <span class="output-time">{{ line.time }}</span>
              <span class="output-message">{{ line.message }}</span>
            </div>
            <div v-if="outputLines.length === 0" class="empty-state">
              <p>暂无输出信息</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 问题面板 -->
      <div
        v-else-if="layoutStore.activeBottomPanel === 'problems'"
        class="panel-section"
      >
        <div class="problems-container">
          <div class="problems-header">
            <span class="problems-title">⚠️ 问题</span>
            <span class="problems-count">{{ problems.length }}</span>
          </div>
          <div class="problems-body">
            <div
              class="problem-item"
              v-for="(problem, index) in problems"
              :key="index"
            >
              <span class="problem-icon">{{
                problem.type === "error" ? "❌" : "⚠️"
              }}</span>
              <span class="problem-message">{{ problem.message }}</span>
              <span class="problem-location">{{ problem.location }}</span>
            </div>
            <div v-if="problems.length === 0" class="empty-state">
              <p>✨ 没有检测到问题</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 调试控制台面板 -->
      <div
        v-else-if="layoutStore.activeBottomPanel === 'debug'"
        class="panel-section"
      >
        <div class="debug-container">
          <div class="debug-header">
            <span class="debug-title">🐛 调试控制台</span>
          </div>
          <div class="debug-body">
            <div class="empty-state">
              <p>调试会话未启动</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { useLayoutStore } from "@/stores/layout";
import Terminal from "../BottomPanel/Terminal.vue";

const layoutStore = useLayoutStore();

// 终端相关
function clearTerminal() {
  console.log("清空终端");
}

function addNewTerminal() {
  console.log("新建终端");
}

// 输出相关
interface OutputLine {
  time: string;
  message: string;
}

const outputLines = ref<OutputLine[]>([
  { time: "10:30:15", message: "应用启动成功" },
  {
    time: "10:30:16",
    message: "加载工作区: /Users/yanghao/Desktop/wails3-vue3-naviteui-test",
  },
  { time: "10:30:17", message: "初始化编辑器完成" },
]);

// 问题相关
interface Problem {
  type: "error" | "warning";
  message: string;
  location: string;
}

const problems = ref<Problem[]>([
  {
    type: "warning",
    message: "未使用的变量 'temp'",
    location: "main.go:15",
  },
]);
</script>

<style scoped>
.bottom-panel {
  display: flex;
  flex-direction: column;
  height: 100%;
  background-color: #252526;
  border-top: 1px solid #3e3e42;
}

.panel-tabs {
  display: flex;
  align-items: center;
  background-color: #252526;
  border-bottom: 1px solid #3e3e42;
  min-height: 35px;
}

.panel-tab {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  cursor: pointer;
  color: #969696;
  transition: all 0.2s;
  border-bottom: 2px solid transparent;
  font-size: 12px;
}

.panel-tab:hover {
  background-color: #2a2d2e;
  color: #cccccc;
}

.panel-tab.active {
  background-color: #1e1e1e;
  color: #ffffff;
  border-bottom-color: #007acc;
}

.tab-icon {
  font-size: 12px;
}

.panel-actions {
  margin-left: auto;
  display: flex;
  align-items: center;
  gap: 4px;
  padding-right: 8px;
}

.panel-content {
  flex: 1;
  overflow: hidden;
  background-color: #1e1e1e;
}

.panel-section {
  display: flex;
  flex-direction: column;
  height: 100%;
}

/* 终端样式 */
.terminal-container,
.output-container,
.problems-container,
.debug-container {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.terminal-header,
.output-header,
.problems-header,
.debug-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 8px 12px;
  background-color: #252526;
  border-bottom: 1px solid #3e3e42;
}

.terminal-title,
.output-title,
.problems-title,
.debug-title {
  font-size: 12px;
  font-weight: 600;
  color: #cccccc;
}

.terminal-actions {
  display: flex;
  gap: 4px;
}

.terminal-body,
.output-body,
.problems-body,
.debug-body {
  flex: 1;
  padding: 12px;
  overflow-y: auto;
  font-family: "Consolas", "Courier New", monospace;
}

.terminal-line {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #cccccc;
  font-size: 13px;
}

.terminal-prompt {
  color: #4ec9b0;
  font-weight: bold;
}

.terminal-path {
  color: #569cd6;
}

.terminal-cursor {
  animation: blink 1s infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* 输出样式 */
.output-line {
  display: flex;
  gap: 12px;
  padding: 4px 0;
  font-size: 12px;
  color: #cccccc;
}

.output-time {
  color: #858585;
  min-width: 70px;
}

.output-message {
  flex: 1;
}

/* 问题样式 */
.problems-count {
  margin-left: 8px;
  padding: 2px 6px;
  background-color: #0e639c;
  color: #ffffff;
  border-radius: 10px;
  font-size: 11px;
}

.problem-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  margin-bottom: 4px;
  background-color: #2a2d2e;
  border-radius: 4px;
  font-size: 12px;
  color: #cccccc;
}

.problem-icon {
  font-size: 14px;
}

.problem-message {
  flex: 1;
}

.problem-location {
  color: #858585;
  font-size: 11px;
}

.empty-state {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  color: #858585;
  font-size: 13px;
}

:deep(.n-button) {
  padding: 4px !important;
  height: 24px !important;
  min-width: 24px !important;
}
</style>
