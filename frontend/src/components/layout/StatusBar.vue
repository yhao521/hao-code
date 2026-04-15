<template>
  <div class="statusbar">
    <div class="statusbar-left">
      <span class="status-item" @click="showGitInfo">
        <NIcon><GitBranchOutline /></NIcon>
        {{ currentBranch || "No Repository" }}
      </span>
      <span class="status-item" v-if="hasChanges">
        {{ changeCount }} changes
      </span>
    </div>
    <div class="statusbar-right">
      <span class="status-item" v-if="activeEditor">
        Ln {{ cursorLine }}, Col {{ cursorCol }}
      </span>
      <span class="status-item" v-if="activeEditor"> Spaces: 2 </span>
      <span class="status-item" v-if="activeEditor"> UTF-8 </span>
      <span class="status-item" v-if="activeEditor">
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
import { GitBranchOutline, NotificationsOutline } from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import { useGitStore } from "@/stores/git";

const editorStore = useEditorStore();
const gitStore = useGitStore();

const activeEditor = computed(() => editorStore.activeTab);
const currentBranch = computed(() => gitStore.currentBranch);
const hasChanges = computed(() => gitStore.changes.length > 0);
const changeCount = computed(() => gitStore.changes.length);

// 模拟光标位置（实际应该从 Monaco Editor 获取）
const cursorLine = ref(1);
const cursorCol = ref(1);

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

.status-item .n-icon {
  font-size: 14px;
}
</style>
