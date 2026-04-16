<template>
  <div class="virtual-tree" ref="containerRef" @scroll="handleScroll">
    <div :style="{ height: totalHeight + 'px' }" class="virtual-list-phantom">
      <div
        v-for="item in visibleData"
        :key="item.key"
        class="virtual-tree-node"
        :style="{
          transform: `translateY(${item.top}px)`,
          paddingLeft: `${item.level * 16 + 8}px`,
        }"
        @click="handleNodeClick(item)"
      >
        <span
          class="toggle-icon"
          v-if="item.isDir"
          @click.stop="toggleExpand(item)"
        >
          {{ expandedKeys.has(item.key) ? "▼" : "▶" }}
        </span>
        <span v-else class="toggle-placeholder"></span>

        <NIcon
          :component="getFileIcon(item)"
          :color="getFileIconColor(item)"
          size="14"
        />
        <span class="node-name">{{ item.name }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from "vue";
import { NIcon } from "naive-ui";
import { FolderOpenOutline, FileTrayOutline } from "@vicons/ionicons5";

const props = defineProps<{
  data: any[];
  expandedKeys: Set<string>;
}>();

const emit = defineEmits(["update:expandedKeys", "select"]);

const containerRef = ref<HTMLElement>();
const itemHeight = 24; // 每个节点的高度
const visibleCount = 20; // 可视区域显示的节点数
const startIndex = ref(0);

// 将树形结构扁平化以便虚拟滚动
const flatData = computed(() => {
  const result: any[] = [];
  const flatten = (nodes: any[], level: number) => {
    for (const node of nodes) {
      result.push({ ...node, level });
      if (node.isDir && props.expandedKeys.has(node.key) && node.children) {
        flatten(node.children, level + 1);
      }
    }
  };
  flatten(props.data, 0);
  return result;
});

const totalHeight = computed(() => flatData.value.length * itemHeight);

const visibleData = computed(() => {
  return flatData.value.slice(
    startIndex.value,
    startIndex.value + visibleCount,
  );
});

function handleScroll() {
  if (!containerRef.value) return;
  const scrollTop = containerRef.value.scrollTop;
  startIndex.value = Math.floor(scrollTop / itemHeight);
}

function toggleExpand(item: any) {
  const newKeys = new Set(props.expandedKeys);
  if (newKeys.has(item.key)) {
    newKeys.delete(item.key);
  } else {
    newKeys.add(item.key);
  }
  emit("update:expandedKeys", newKeys);
}

function handleNodeClick(item: any) {
  emit("select", item);
}

function getFileIcon(item: any) {
  return item.isDir ? FolderOpenOutline : FileTrayOutline;
}

function getFileIconColor(item: any) {
  return item.isDir ? "#D7BA7D" : "#CCCCCC";
}
</script>

<style scoped>
.virtual-tree {
  height: 100%;
  overflow-y: auto;
  position: relative;
}

.virtual-list-phantom {
  position: relative;
  width: 100%;
}

.virtual-tree-node {
  position: absolute;
  left: 0;
  right: 0;
  height: 24px;
  display: flex;
  align-items: center;
  cursor: pointer;
  font-size: 13px;
  color: #cccccc;
}

.virtual-tree-node:hover {
  background-color: #2a2d2e;
}

.toggle-icon {
  width: 16px;
  font-size: 10px;
  text-align: center;
  margin-right: 4px;
  color: #858585;
}

.toggle-placeholder {
  width: 20px;
}

.node-name {
  margin-left: 6px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
