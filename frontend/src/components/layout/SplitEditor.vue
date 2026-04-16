<template>
  <div class="split-editor" :style="layoutStyle">
    <EditorGroup
      v-for="(group, index) in store.editorGroups"
      :key="group.id"
      :group="group"
      :class="{ 'active-group': store.activeGroupId === group.id }"
      :style="groupStyle(index)"
    />

    <!-- 简单的分割线/调整器 -->
    <div
      v-if="store.editorGroups.length > 1"
      class="resizer"
      @mousedown="startResize"
    </div>
    
    <BlameDetailModal ref="blameModalRef" />
  </div>
</template>

<script setup lang="ts">
import { computed, ref, onMounted, type CSSProperties } from "vue";
import { useEditorStore } from "@/stores/editor";
import EditorGroup from "./EditorGroup.vue";
import BlameDetailModal from "./BlameDetailModal.vue";

const store = useEditorStore();
const blameModalRef = ref<any>(null);
const isVertical = ref(false); // 默认为水平分屏（左右排列）

const layoutStyle = computed<CSSProperties>(() => ({
  flexDirection: isVertical.value ? "column" : "row",
}));

function groupStyle(index: number) {
  const total = store.editorGroups.length;
  return {
    flex: 1,
    width: `${100 / total}%`,
    height: "100%",
  };
}

// 简单的 Resize 逻辑（后续可以优化为更复杂的布局树）
function startResize(e: MouseEvent) {
  e.preventDefault();
  // TODO: 实现拖拽改变大小
}

onMounted(() => {
  window.addEventListener('editor:show-blame', (e: any) => {
    if (blameModalRef.value) {
      blameModalRef.value.show(e.detail);
    }
  });
});
</script>

<style scoped>
.split-editor {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.resizer {
  width: 4px;
  background-color: #2b2b2b;
  cursor: col-resize;
  transition: background-color 0.2s;
}

.resizer:hover {
  background-color: #007acc;
}
</style>
