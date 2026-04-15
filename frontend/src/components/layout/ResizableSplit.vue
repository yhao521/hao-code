<template>
  <div
    class="resizable-split"
    :style="{ flexDirection: horizontal ? 'row' : 'column' }"
  >
    <!-- 左侧/上侧面板 -->
    <div class="panel panel-1" :style="panel1Style">
      <slot name="1"></slot>
    </div>

    <!-- 拖拽条 -->
    <div
      class="resize-handle"
      :class="{
        'resize-handle--horizontal': horizontal,
        'resize-handle--hover': isHovering,
        'resize-handle--dragging': isDragging,
      }"
      @mousedown="startResize"
      @mouseenter="isHovering = true"
      @mouseleave="isHovering = false"
    >
      <div class="resize-handle-line"></div>
    </div>

    <!-- 右侧/下侧面板 -->
    <div class="panel panel-2" :style="panel2Style">
      <slot name="2"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch } from "vue";

interface Props {
  min?: number;
  max?: number;
  horizontal?: boolean;
  defaultSize?: number;
  direction?: "ltr" | "rtl"; // left-to-right or right-to-left
}

const props = withDefaults(defineProps<Props>(), {
  min: 150,
  max: 600,
  horizontal: true,
  defaultSize: 240,
  direction: "ltr",
});

const emit = defineEmits<{
  "update:size": [value: number];
}>();

const containerRef = ref<HTMLElement | null>(null);
const size = ref(props.defaultSize);
const isDragging = ref(false);
const isHovering = ref(false);
const startPos = ref(0);
const startSize = ref(0);

// 监听 defaultSize 变化
watch(
  () => props.defaultSize,
  (newSize) => {
    size.value = newSize;
  },
);

const panel1Style = computed(() => {
  // 根据 direction 决定面板顺序
  if (props.direction === "rtl") {
    return {
      [props.horizontal ? "width" : "height"]: `${size.value}px`,
      flexShrink: 0,
      flexGrow: 0,
      order: 2,
    };
  }
  return {
    [props.horizontal ? "width" : "height"]: `${size.value}px`,
    flexShrink: 0,
    flexGrow: 0,
    order: 1,
  };
});

const panel2Style = computed(() => {
  if (props.direction === "rtl") {
    return {
      flex: 1,
      minWidth: 0,
      width: 0,
      overflow: "hidden",
      order: 1,
    };
  }
  return {
    flex: 1,
    minWidth: 0,
    width: 0,
    overflow: "hidden",
    order: 2,
  };
});

function startResize(e: MouseEvent) {
  isDragging.value = true;
  startPos.value = props.horizontal ? e.clientX : e.clientY;
  startSize.value = size.value;

  document.addEventListener("mousemove", handleResize);
  document.addEventListener("mouseup", stopResize);
  document.body.style.cursor = props.horizontal ? "col-resize" : "row-resize";
  document.body.style.userSelect = "none";
}

function handleResize(e: MouseEvent) {
  if (!isDragging.value) return;

  const currentPos = props.horizontal ? e.clientX : e.clientY;
  let delta = currentPos - startPos.value;

  // RTL 模式下反向计算
  if (props.direction === "rtl") {
    delta = -delta;
  }

  let newSize = startSize.value + delta;

  // 限制范围
  newSize = Math.max(props.min, Math.min(props.max, newSize));

  size.value = newSize;
  emit("update:size", newSize);
}

function stopResize() {
  isDragging.value = false;
  document.removeEventListener("mousemove", handleResize);
  document.removeEventListener("mouseup", stopResize);
  document.body.style.cursor = "";
  document.body.style.userSelect = "";
}

onUnmounted(() => {
  document.removeEventListener("mousemove", handleResize);
  document.removeEventListener("mouseup", stopResize);
});

defineExpose({
  size,
});
</script>

<style scoped>
.resizable-split {
  display: flex;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.panel {
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.panel-1 {
  flex-shrink: 0;
  flex-grow: 0;
}

.panel-2 {
  flex: 1;
  min-width: 0;
  min-height: 0;
  overflow: hidden;
}

.resize-handle {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
  transition: background-color 0.2s;
}

.resize-handle--horizontal {
  width: 4px;
  cursor: col-resize;
}

.resize-handle--hover,
.resize-handle--dragging {
  background-color: rgba(14, 99, 156, 0.2);
}

.resize-handle--dragging {
  background-color: rgba(14, 99, 156, 0.3);
}

.resize-handle-line {
  position: absolute;
  background-color: #3e3e42;
  transition:
    background-color 0.2s,
    width 0.2s;
}

.resize-handle--horizontal .resize-handle-line {
  width: 1px;
  height: 40px;
  left: 50%;
  transform: translateX(-50%);
}

.resize-handle--hover .resize-handle-line,
.resize-handle--dragging .resize-handle-line {
  background-color: #0e639c;
  width: 2px;
}
</style>
