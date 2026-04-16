<template>
  <div class="call-hierarchy-node" :style="{ paddingLeft: level * 12 + 'px' }">
    <div class="node-item" @click="toggleExpand" @dblclick="handleJump">
      <n-icon v-if="hasChildren" class="expand-icon">
        <ChevronDownOutline v-if="expanded" />
        <ChevronForwardOutline v-else />
      </n-icon>
      <span class="node-name">{{ item.name }}</span>
      <span class="node-detail">{{ item.detail }}</span>
    </div>

    <div v-if="expanded && hasChildren" class="children-container">
      <CallHierarchyNode
        v-for="(child, index) in children"
        :key="index"
        :item="child"
        :level="level + 1"
        @jump-to-location="$emit('jumpToLocation', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from "vue";
import { NIcon } from "naive-ui";
import { ChevronDownOutline, ChevronForwardOutline } from "@vicons/ionicons5";
import { lspManager } from "@/utils/lspManager";
import { useEditorStore } from "@/stores/editor";

const props = defineProps<{
  item: any;
  level: number;
}>();

const emit = defineEmits(["jumpToLocation"]);

const editorStore = useEditorStore();
const expanded = ref(false);
const children = ref<any[]>([]);
const loading = ref(false);

const hasChildren = computed(() => {
  return children.value.length > 0 || !expanded.value; // 允许尝试展开
});

const toggleExpand = async () => {
  if (!expanded.value && children.value.length === 0) {
    await loadChildren();
  }
  expanded.value = !expanded.value;
};

const loadChildren = async () => {
  if (!editorStore.activeTab) return;

  loading.value = true;
  try {
    const langId = getLanguage(editorStore.activeTab.path);
    const calls = await lspManager.getIncomingCalls(langId, props.item);
    if (calls) {
      children.value = calls.map((c: any) => c.from);
    }
  } catch (e) {
    console.error("Failed to load incoming calls:", e);
  } finally {
    loading.value = false;
  }
};

const handleJump = () => {
  if (props.item.uri) {
    emit("jumpToLocation", {
      uri: props.item.uri,
      range: props.item.range,
    });
  }
};

function getLanguage(path: string): string {
  const ext = path.split(".").pop()?.toLowerCase();
  const map: Record<string, string> = {
    go: "go",
    ts: "typescript",
    js: "javascript",
    py: "python",
    java: "java",
  };
  return map[ext || ""] || "plaintext";
}
</script>

<style scoped>
.call-hierarchy-node {
  user-select: none;
}

.node-item {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 13px;
}

.node-item:hover {
  background-color: #2a2d2e;
}

.expand-icon {
  margin-right: 4px;
  font-size: 12px;
  color: #cccccc;
}

.node-name {
  color: #dcdcaa;
  margin-right: 6px;
}

.node-detail {
  color: #888;
  font-size: 11px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.children-container {
  border-left: 1px solid #3e3e42;
  margin-left: 6px;
}
</style>
