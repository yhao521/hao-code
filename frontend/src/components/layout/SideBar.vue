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
        <component :is="item.icon" class="activity-icon" />
      </div>
    </div>

    <!-- 侧边栏内容 -->
    <div class="sidebar-content">
      <WelcomeView v-if="!editorStore.workspace && activeView === 'explorer'" />
      <FileExplorer v-else-if="activeView === 'explorer'" />
      <SearchPanel v-else-if="activeView === 'search'" />
      <GitPanel v-else-if="activeView === 'git'" />
      <OutlinePanel v-else-if="activeView === 'outline'" />
      <CallHierarchyPanel v-else-if="activeView === 'call-hierarchy'" />
      <TypeHierarchyPanel v-else-if="activeView === 'type-hierarchy'" />
      <TasksPanel v-else-if="activeView === 'tasks'" />
      <ExtensionsPanel v-else-if="activeView === 'extensions'" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NIcon } from "naive-ui";
import type { Component } from "vue";
import {
  FolderOutline,
  SearchOutline,
  GitBranchOutline,
  ListOutline,
  CubeOutline,
  CallOutline,
  LayersOutline,
  PlayOutline,
} from "@vicons/ionicons5";
import FileExplorer from "./FileExplorer.vue";
import SearchPanel from "../SearchPanel.vue";
import GitPanel from "./GitPanel.vue";
import OutlinePanel from "./OutlinePanel.vue";
import ExtensionsPanel from "./ExtensionsPanel.vue";
import CallHierarchyPanel from "../sidebar/CallHierarchyPanel.vue";
import TypeHierarchyPanel from "../sidebar/TypeHierarchyPanel.vue";
import TasksPanel from "../sidebar/TasksPanel.vue";
import WelcomeView from "../sidebar/WelcomeView.vue";

interface ActivityItem {
  id: string;
  title: string;
  icon: Component;
}

import { useEditorStore } from "@/stores/editor";

const editorStore = useEditorStore();
const activeView = ref("explorer");

const activityItems: ActivityItem[] = [
  {
    id: "explorer",
    title: "资源管理器",
    icon: FolderOutline,
  },
  {
    id: "search",
    title: "搜索",
    icon: SearchOutline,
  },
  {
    id: "git",
    title: "源代码管理",
    icon: GitBranchOutline,
  },
  {
    id: "outline",
    title: "大纲",
    icon: ListOutline,
  },
  {
    id: "call-hierarchy",
    title: "调用层级",
    icon: CallOutline,
  },
  {
    id: "type-hierarchy",
    title: "类型层次结构",
    icon: LayersOutline,
  },
  {
    id: "tasks",
    title: "任务",
    icon: PlayOutline,
  },
  {
    id: "extensions",
    title: "扩展",
    icon: CubeOutline,
  },
];

function handleViewChange(id: string) {
  activeView.value = activeView.value === id ? "" : id;
}
</script>

<style scoped>
.sidebar {
  display: flex;
  height: 100%;
  background-color: #252526;
}

/* Activity Bar - 左侧图标栏 */
.activity-bar {
  width: 48px;
  min-width: 48px;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding-top: 10px;
  background-color: #333333;
  border-right: 1px solid #1e1e1e;
  flex-shrink: 0;
}

.activity-item {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  color: #858585;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
}

.activity-icon {
  width: 26px;
  height: 26px;
  fill: currentColor;
  filter: drop-shadow(0 1px 2px rgba(0, 0, 0, 0.3));
}

.activity-item:hover {
  color: #ffffff;
  transform: scale(1.05);
}

.activity-item.active {
  color: #ffffff;
}

.activity-item.active::before {
  content: "";
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
  background-color: #007acc;
  box-shadow: 0 0 8px rgba(0, 122, 204, 0.5);
}

/* 侧边栏内容区域 */
.sidebar-content {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  min-width: 0; /* 防止内容溢出 */
}
</style>
