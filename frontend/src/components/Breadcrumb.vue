<template>
  <div class="breadcrumb">
    <span
      v-for="(segment, index) in segments"
      :key="index"
      class="breadcrumb-segment"
    >
      <span
        @click="navigateTo(segment.path)"
        class="breadcrumb-link"
        :class="{ clickable: segment.clickable }"
      >
        {{ segment.name }}
      </span>
      <span v-if="index < segments.length - 1" class="separator">›</span>
    </span>
  </div>
</template>

<script setup lang="ts">
import { computed } from "vue";
import { useEditorStore } from "@/stores/editor";

const editorStore = useEditorStore();

const props = defineProps<{
  path: string;
}>();

interface Segment {
  name: string;
  path: string;
  clickable: boolean;
}

const segments = computed(() => {
  const parts = props.path.split("/");
  const result: Segment[] = [];

  parts.forEach((part, index) => {
    if (part) {
      result.push({
        name: part,
        path: parts.slice(0, index + 1).join("/"),
        clickable: index < parts.length - 1, // 最后一部分不可点击
      });
    }
  });

  return result;
});

function navigateTo(path: string) {
  // 尝试在文件树中展开并选中该路径
  window.dispatchEvent(
    new CustomEvent("breadcrumb:navigate", { detail: { path } }),
  );
}
</script>

<style scoped>
.breadcrumb {
  display: flex;
  align-items: center;
  padding: 4px 12px;
  background-color: #1e1e1e;
  border-bottom: 1px solid #2b2b2b;
  font-size: 12px;
  min-height: 22px;
  overflow-x: auto;
  white-space: nowrap;
  flex-shrink: 0;
}

.breadcrumb-segment {
  display: flex;
  align-items: center;
}

.breadcrumb-link {
  color: #969696;
  transition: color 0.2s;
}

.breadcrumb-link.clickable {
  cursor: pointer;
}

.breadcrumb-link.clickable:hover {
  color: #ffffff;
  text-decoration: underline;
}

.separator {
  margin: 0 4px;
  color: #555555;
  font-size: 14px;
}
</style>
