<template>
  <div class="statusbar">
    <div class="statusbar-left">
      <span class="status-item" @click="showGitInfo" v-if="currentBranch">
        <NIcon><GitBranchOutline /></NIcon>
        {{ currentBranch }}
      </span>
      <span class="status-item" v-if="hasChanges">
        <NIcon><CreateOutline /></NIcon>
        {{ changeCount }}
      </span>
      <span
        class="status-item"
        v-if="diagnosticsCount > 0"
        :class="{ 'has-errors': hasErrors }"
      >
        <NIcon><BugOutline /></NIcon>
        {{ diagnosticsCount }}
      </span>
    </div>
    <div class="statusbar-right">
      <span class="status-item" v-if="activeEditor">
        Ln {{ cursorLine }}, Col {{ cursorCol }}
      </span>
      <span class="status-item" v-if="activeEditor"> Spaces: 2 </span>
      <span class="status-item" v-if="activeEditor"> UTF-8 </span>
      <span class="status-item" v-if="activeEditor" @click="changeLanguage">
        {{ languageName }}
      </span>
      <span class="status-item">
        <NIcon><NotificationsOutline /></NIcon>
        0
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from "vue";
import { NIcon } from "naive-ui";
import {
  GitBranchOutline,
  NotificationsOutline,
  CreateOutline,
  BugOutline,
} from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import { useGitStore } from "@/stores/git";
import { useDiagnosticsStore } from "@/stores/diagnostics";

const editorStore = useEditorStore();
const gitStore = useGitStore();
const diagnosticsStore = useDiagnosticsStore();

const activeEditor = computed(() => editorStore.activeTab);
const currentBranch = computed(() => gitStore.currentBranch);
const hasChanges = computed(() => gitStore.changes.length > 0);
const changeCount = computed(() => gitStore.changes.length);

// 诊断信息同步
const diagnosticsCount = computed(() => {
  if (!activeEditor.value) return 0;
  const markers = diagnosticsStore.markers[activeEditor.value.path] || [];
  return markers.length;
});

const hasErrors = computed(() => {
  if (!activeEditor.value) return false;
  const markers = diagnosticsStore.markers[activeEditor.value.path] || [];
  return markers.some((m: any) => m.severity === 8); // monaco.MarkerSeverity.Error
});

// 模拟光标位置（实际应该从 Monaco Editor 获取）
const cursorLine = ref(1);
const cursorCol = ref(1);

// 监听编辑器光标变化（通过自定义事件或 Store）
window.addEventListener("editor:cursor-change", (e: any) => {
  if (e.detail) {
    cursorLine.value = e.detail.line;
    cursorCol.value = e.detail.col;
  }
});

const languageName = computed(() => {
  if (!activeEditor.value) return "Plain Text";
  const langMap: Record<string, string> = {
    typescript: "TypeScript",
    javascript: "JavaScript",
    python: "Python",
    go: "Go",
    java: "Java",
    html: "HTML",
    css: "CSS",
    json: "JSON",
    markdown: "Markdown",
  };
  return langMap[activeEditor.value.language || ""] || "Plain Text";
});

function showGitInfo() {
  // TODO: 显示 Git 详细信息
  console.log("Show Git info");
}

function changeLanguage() {
  console.log("Change language mode");
}
</script>

<style scoped>
.statusbar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  height: 22px;
  background-color: #007acc;
  color: white;
  font-size: 12px;
  padding: 0 10px;
  user-select: none;
}

.statusbar-left,
.statusbar-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.status-item {
  display: flex;
  align-items: center;
  gap: 4px;
  cursor: pointer;
  padding: 2px 6px;
  border-radius: 3px;
  transition: background-color 0.15s;
  height: 22px;
}

.status-item:hover {
  background-color: rgba(255, 255, 255, 0.12);
}

.status-item.has-errors {
  color: #f48771;
}

.status-item .n-icon {
  font-size: 14px;
}
</style>
