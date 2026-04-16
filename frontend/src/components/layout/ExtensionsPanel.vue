<template>
  <div class="extensions-panel">
    <div class="panel-header">
      <span>扩展</span>
      <n-input
        v-model:value="searchQuery"
        placeholder="搜索扩展..."
        size="small"
        clearable
      />
    </div>

    <div class="extensions-content">
      <div v-if="loading" class="loading-state">加载中...</div>
      <div v-else-if="filteredPlugins.length === 0" class="empty-state">
        <p>没有找到匹配的扩展</p>
      </div>

      <div
        v-for="plugin in filteredPlugins"
        :key="plugin.name"
        class="extension-item"
      >
        <div class="extension-icon">🧩</div>
        <div class="extension-info">
          <div class="extension-name">{{ plugin.name }}</div>
          <div class="extension-desc">
            {{ plugin.description || "暂无描述" }}
          </div>
          <div class="extension-meta">
            <span class="version">v{{ plugin.version }}</span>
            <span class="author" v-if="plugin.author">{{ plugin.author }}</span>
          </div>
        </div>
        <div class="extension-actions">
          <n-button
            size="tiny"
            type="primary"
            ghost
            @click="handleActivate(plugin.name)"
          >
            启用
          </n-button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { NInput, NButton } from "naive-ui";
import {
  GetInstalledPlugins,
  ActivatePlugin,
} from "@wails/backend/appservice.js";

interface PluginManifest {
  name: string;
  version: string;
  description: string;
  main: string;
  author: string;
  license: string;
}

const plugins = ref<PluginManifest[]>([]);
const searchQuery = ref("");
const loading = ref(false);

const filteredPlugins = computed(() => {
  if (!searchQuery.value) return plugins.value;
  const query = searchQuery.value.toLowerCase();
  return plugins.value.filter(
    (p) =>
      p.name.toLowerCase().includes(query) ||
      (p.description && p.description.toLowerCase().includes(query)),
  );
});

async function loadPlugins() {
  loading.value = true;
  try {
    const result = await GetInstalledPlugins();
    plugins.value = result || [];
  } catch (error) {
    console.error("Failed to load plugins:", error);
  } finally {
    loading.value = false;
  }
}

async function handleActivate(name: string) {
  try {
    await ActivatePlugin(name);
    // TODO: 显示成功提示
    console.log(`Plugin ${name} activated`);
  } catch (error) {
    console.error(`Failed to activate plugin ${name}:`, error);
  }
}

onMounted(() => {
  loadPlugins();
});
</script>

<style scoped>
.extensions-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #252526;
  color: #cccccc;
}

.panel-header {
  padding: 8px 12px;
  border-bottom: 1px solid #3c3c3c;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.panel-header span {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
}

.extensions-content {
  flex: 1;
  overflow-y: auto;
  padding: 4px 0;
}

.extension-item {
  display: flex;
  padding: 8px 12px;
  gap: 10px;
  cursor: pointer;
  transition: background-color 0.15s;
}

.extension-item:hover {
  background-color: #2a2d2e;
}

.extension-icon {
  width: 40px;
  height: 40px;
  background-color: #3c3c3c;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  border-radius: 4px;
}

.extension-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.extension-name {
  font-weight: 600;
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.extension-desc {
  font-size: 12px;
  color: #969696;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.extension-meta {
  display: flex;
  gap: 8px;
  font-size: 11px;
  color: #858585;
  margin-top: 2px;
}

.extension-actions {
  display: flex;
  align-items: center;
}

.loading-state,
.empty-state {
  padding: 20px;
  text-align: center;
  color: #858585;
  font-size: 13px;
}
</style>
