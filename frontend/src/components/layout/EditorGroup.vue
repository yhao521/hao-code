<template>
  <div class="editor-group" @click="activateGroup">
    <div class="tabs-header">
      <div
        v-for="tab in group.tabs"
        :key="tab.id"
        class="tab-item"
        :class="{ active: tab.id === group.activeTabId }"
        @click.stop="selectTab(tab.id)"
      >
        <span class="tab-name">{{ tab.name }}</span>
        <span class="close-icon" @click.stop="handleCloseTab(tab.id)">×</span>
      </div>
    </div>
    <div class="monaco-container" ref="monacoRef"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed } from "vue";
import * as monaco from "monaco-editor";
import { useEditorStore, type EditorGroup } from "@/stores/editor";

const props = defineProps<{
  group: EditorGroup;
}>();

const store = useEditorStore();
const monacoRef = ref<HTMLDivElement>();
let editorInstance: monaco.editor.IStandaloneCodeEditor | null = null;

const isActive = computed(() => store.activeGroupId === props.group.id);

function activateGroup() {
  store.setActiveGroup(props.group.id);
}

function selectTab(id: string) {
  props.group.activeTabId = id;
  activateGroup();
}

function handleCloseTab(id: string) {
  store.closeTab(id);
}

onMounted(() => {
  if (monacoRef.value) {
    editorInstance = monaco.editor.create(monacoRef.value, {
      value: "",
      language: "plaintext",
      theme: "vs-dark",
      automaticLayout: true,
      minimap: { enabled: false },
      scrollBeyondLastLine: false,
    });
  }
});

watch(
  () => props.group.activeTabId,
  (newId) => {
    const tab = props.group.tabs.find((t) => t.id === newId);
    if (tab && editorInstance) {
      editorInstance.setValue(tab.content || "");
      monaco.editor.setModelLanguage(
        editorInstance.getModel()!,
        tab.language || "plaintext",
      );
    }
  },
  { immediate: true },
);

watch(
  () => isActive.value,
  (active) => {
    if (active && editorInstance) {
      editorInstance.layout();
      editorInstance.focus();
    }
  },
);
</script>

<style scoped>
.editor-group {
  display: flex;
  flex-direction: column;
  height: 100%;
  border-right: 1px solid #2b2b2b;
}

.tabs-header {
  display: flex;
  background-color: #252526;
  overflow-x: auto;
}

.tab-item {
  padding: 8px 12px;
  font-size: 13px;
  cursor: pointer;
  background-color: #2d2d2d;
  border-right: 1px solid #252526;
  display: flex;
  align-items: center;
  min-width: 100px;
}

.tab-item.active {
  background-color: #1e1e1e;
  color: #ffffff;
  border-top: 1px solid #007acc;
}

.tab-name {
  margin-right: 8px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.close-icon {
  opacity: 0;
  transition: opacity 0.2s;
}

.tab-item:hover .close-icon {
  opacity: 1;
}

.monaco-container {
  flex: 1;
  width: 100%;
}
</style>
