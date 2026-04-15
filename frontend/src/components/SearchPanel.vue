<template>
  <div class="search-panel">
    <!-- 搜索输入区域 -->
    <div class="search-input-section">
      <NInput
        v-model:value="searchText"
        placeholder="在文件中搜索..."
        size="small"
        @keyup.enter="handleSearch"
      >
        <template #prefix>
          <n-icon :component="SearchIcon" />
        </template>
      </NInput>

      <div class="search-options">
        <NCheckbox v-model:checked="caseSensitive" size="small">
          区分大小写
        </NCheckbox>
        <NCheckbox v-model:checked="useRegex" size="small">
          使用正则
        </NCheckbox>
      </div>

      <NButton
        size="small"
        type="primary"
        @click="handleSearch"
        :loading="searching"
        style="width: 100%; margin-top: 8px"
      >
        搜索
      </NButton>
    </div>

    <!-- 搜索结果列表 -->
    <div class="search-results" v-if="results.length > 0">
      <div class="results-header">
        <span>找到 {{ results.length }} 个结果</span>
        <NButton size="tiny" text @click="clearResults">清除</NButton>
      </div>

      <div class="results-list">
        <div
          v-for="(result, index) in results"
          :key="`${result.filePath}:${result.lineNumber}`"
          class="result-item"
          @click="openFile(result)"
        >
          <div class="file-info">
            <n-icon :component="FileIcon" size="14" />
            <span class="file-path">{{
              getRelativePath(result.filePath)
            }}</span>
          </div>
          <div class="line-info">
            <span class="line-number">第 {{ result.lineNumber }} 行</span>
            <span class="line-content">{{
              truncateLine(result.lineContent)
            }}</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 无结果提示 -->
    <div class="no-results" v-else-if="searched && !searching">
      <n-icon :component="EmptyIcon" size="48" color="#858585" />
      <p>未找到匹配的结果</p>
    </div>

    <!-- 初始状态 -->
    <div class="search-placeholder" v-else>
      <n-icon :component="SearchIcon" size="48" color="#858585" />
      <p>输入关键词开始搜索</p>
      <p class="hint">支持在整个工作区中搜索文件内容</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NInput, NCheckbox, NButton, NIcon } from "naive-ui";
import {
  Search as SearchIcon,
  Document as FileIcon,
  Sad as EmptyIcon,
} from "@vicons/ionicons5";
import { SearchInFiles, type SearchOptions } from "@wails/go/backend/App";
import { useEditorStore } from "@/stores/editor";
import { useMessage } from "naive-ui";

const editorStore = useEditorStore();
const message = useMessage();

// 搜索状态
const searchText = ref("");
const caseSensitive = ref(false);
const useRegex = ref(false);
const searching = ref(false);
const searched = ref(false);

interface SearchResult {
  filePath: string;
  lineNumber: number;
  lineContent: string;
}

const results = ref<SearchResult[]>([]);

// 执行搜索
async function handleSearch() {
  if (!searchText.value.trim()) {
    message.warning("请输入搜索关键词");
    return;
  }

  if (!editorStore.workspace) {
    message.warning("请先打开一个文件夹");
    return;
  }

  searching.value = true;
  searched.value = true;

  try {
    const workspacePath = editorStore.workspace.path;
    const opts: SearchOptions = {
      rootPath: workspacePath,
      query: searchText.value,
      caseSensitive: caseSensitive.value,
      matchWholeWord: false,
      useRegex: useRegex.value,
      exclude: "node_modules,.git",
    };
    const searchResults = await SearchInFiles(opts);
    results.value = searchResults || [];

    if (results.value.length === 0) {
      message.info("未找到匹配的结果");
    } else {
      message.success(`找到 ${results.value.length} 个结果`);
    }
  } catch (error) {
    console.error("Search failed:", error);
    message.error(
      "搜索失败: " + (error instanceof Error ? error.message : String(error)),
    );
    results.value = [];
  } finally {
    searching.value = false;
  }
}

// 清除结果
function clearResults() {
  results.value = [];
  searched.value = false;
  searchText.value = "";
}

// 打开文件并跳转到指定行
function openFile(result: SearchResult) {
  editorStore.openFile(result.filePath, "");
  // 触发跳转事件，由 EditorArea 监听处理
  window.dispatchEvent(
    new CustomEvent("editor:jump-to-line", {
      detail: { path: result.filePath, line: result.lineNumber },
    }),
  );
}

// 获取相对路径
function getRelativePath(filePath: string): string {
  if (!editorStore.workspace) return filePath;

  const workspacePath = editorStore.workspace.path;
  if (filePath.startsWith(workspacePath)) {
    return filePath.substring(workspacePath.length + 1);
  }
  return filePath;
}

// 获取文件名
function getFileName(filePath: string): string {
  return filePath.split("/").pop() || filePath;
}

// 截断行内容
function truncateLine(content: string, maxLength: number = 100): string {
  if (content.length <= maxLength) return content;
  return content.substring(0, maxLength) + "...";
}
</script>

<style scoped>
.search-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #252526;
  color: #cccccc;
}

.search-input-section {
  padding: 12px;
  border-bottom: 1px solid #3e3e42;
}

.search-options {
  display: flex;
  gap: 12px;
  margin-top: 8px;
  font-size: 12px;
}

.search-results {
  flex: 1;
  overflow-y: auto;
  padding: 8px;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px;
  font-size: 12px;
  color: #858585;
  border-bottom: 1px solid #3e3e42;
}

.results-list {
  padding: 8px 0;
}

.result-item {
  padding: 8px 12px;
  cursor: pointer;
  border-radius: 4px;
  margin-bottom: 4px;
  transition: background-color 0.2s;
}

.result-item:hover {
  background-color: #2a2d2e;
}

.file-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 4px;
  font-size: 12px;
}

.file-path {
  color: #4ec9b0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.line-info {
  display: flex;
  gap: 8px;
  font-size: 11px;
}

.line-number {
  color: #858585;
  min-width: 60px;
}

.line-content {
  color: #cccccc;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex: 1;
}

.no-results,
.search-placeholder {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  color: #858585;
  text-align: center;
  padding: 20px;
}

.no-results p,
.search-placeholder p {
  margin: 8px 0;
  font-size: 13px;
}

.hint {
  font-size: 11px !important;
  color: #666666;
}
</style>
