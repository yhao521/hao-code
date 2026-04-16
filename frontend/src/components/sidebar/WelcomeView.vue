<template>
  <div class="welcome-view">
    <div class="welcome-header">
      <h1>Hao-Code</h1>
      <p class="subtitle">Editing evolved</p>
    </div>

    <div class="welcome-section">
      <h2>Start</h2>
      <ul class="command-list">
        <li @click="handleOpenFolder">
          <NIcon><FolderOpenOutline /></NIcon>
          <span>Open Folder...</span>
        </li>
        <li @click="handleNewFile">
          <NIcon><DocumentTextOutline /></NIcon>
          <span>New Text File</span>
        </li>
      </ul>
    </div>

    <div class="welcome-section" v-if="recentFolders.length > 0">
      <h2>Recent Folders</h2>
      <ul class="recent-list">
        <li
          v-for="(folder, index) in recentFolders"
          :key="index"
          @click="openRecentFolder(folder)"
        >
          <NIcon><FolderOutline /></NIcon>
          <span class="path">{{ folder }}</span>
        </li>
      </ul>
    </div>

    <div class="welcome-section">
      <h2>Help</h2>
      <ul class="command-list">
        <li @click="showShortcuts">
          <NIcon><KeyOutline /></NIcon>
          <span>Show All Commands (Ctrl+Shift+P)</span>
        </li>
        <li @click="openDocs">
          <NIcon><BookOutline /></NIcon>
          <span>Documentation</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import { NIcon, useMessage } from "naive-ui";
import {
  FolderOpenOutline,
  DocumentTextOutline,
  KeyOutline,
  BookOutline,
  FolderOutline,
} from "@vicons/ionicons5";
import { GetRecentFolders } from "@wails/backend/appservice";
import { Events } from "@wailsio/runtime";

const message = useMessage();
const recentFolders = ref<string[]>([]);

onMounted(async () => {
  try {
    recentFolders.value = await GetRecentFolders();
  } catch (error) {
    console.error("Failed to load recent folders:", error);
  }
});

function handleOpenFolder() {
  Events.Emit("menu:open-folder");
}

function handleNewFile() {
  Events.Emit("menu:new-text-file");
}

async function openRecentFolder(path: string) {
  try {
    message.loading("Opening folder...", { duration: 0 });
    // Logic to open folder would go here, likely via store or event
    message.destroyAll();
    message.success(`Opened: ${path.split("/").pop()}`);
  } catch (error) {
    message.destroyAll();
    message.error("Failed to open folder");
  }
}

function showShortcuts() {
  Events.Emit("menu:show-all-commands");
}

function openDocs() {
  window.open("https://github.com/your-repo/hao-code", "_blank");
}
</script>

<style scoped>
.welcome-view {
  padding: 40px;
  color: #cccccc;
  height: 100%;
  overflow-y: auto;
}

.welcome-header h1 {
  font-size: 3em;
  margin: 0;
  font-weight: 300;
  color: #007acc;
}

.subtitle {
  font-size: 1.2em;
  color: #858585;
  margin-top: 5px;
}

.welcome-section {
  margin-top: 40px;
}

.welcome-section h2 {
  font-size: 1.1em;
  font-weight: 600;
  margin-bottom: 15px;
  color: #e7e7e7;
}

.command-list,
.recent-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.command-list li,
.recent-list li {
  display: flex;
  align-items: center;
  padding: 6px 0;
  cursor: pointer;
  transition: color 0.2s;
}

.command-list li:hover,
.recent-list li:hover {
  color: #ffffff;
}

.command-list li .n-icon,
.recent-list li .n-icon {
  margin-right: 10px;
  font-size: 16px;
}

.path {
  font-family: monospace;
  opacity: 0.8;
}
</style>
