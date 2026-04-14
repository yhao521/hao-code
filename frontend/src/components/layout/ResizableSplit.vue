<template>
  <div class="resizable-split" :style="{ flexDirection: horizontal ? 'row' : 'column' }">
    <!-- 左侧/上侧面板 -->
    <div 
      class="panel panel-1"
      :style="panel1Style"
    >
      <slot name="1"></slot>
    </div>
    
    <!-- 拖拽条 -->
    <div 
      class="resize-handle"
      :class="{ 
        'resize-handle--horizontal': horizontal,
        'resize-handle--hover': isHovering,
        'resize-handle--dragging': isDragging
      }"
      @mousedown="startResize"
      @mouseenter="isHovering = true"
      @mouseleave="isHovering = false"
    >
      <div class="resize-handle-line"></div>
    </div>
    
    <!-- 右侧/下侧面板 -->
    <div 
      class="panel panel-2"
      :style="panel2Style"
    >
      <slot name="2"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'

interface Props {
  min?: number
  max?: number
  horizontal?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  min: 150,
  max: 600,
  horizontal: true
})

const emit = defineEmits<{
  'update:size': [value: number]
}>()

const containerRef = ref<HTMLElement | null>(null)
const size = ref(240)
const isDragging = ref(false)
const isHovering = ref(false)
const startPos = ref(0)
const startSize = ref(0)

const panel1Style = computed(() => {
  return {
    [props.horizontal ? 'width' : 'height']: `${size.value}px`,
    flexShrink: 0,
    flexGrow: 0
  }
})

const panel2Style = computed(() => {
  return {
    flex: 1,
    minWidth: 0,
    overflow: 'hidden'
  }
})

function startResize(e: MouseEvent) {
  isDragging.value = true
  startPos.value = props.horizontal ? e.clientX : e.clientY
  startSize.value = size.value
  
  document.addEventListener('mousemove', handleResize)
  document.addEventListener('mouseup', stopResize)
  document.body.style.cursor = props.horizontal ? 'col-resize' : 'row-resize'
  document.body.style.userSelect = 'none'
}

function handleResize(e: MouseEvent) {
  if (!isDragging.value) return
  
  const currentPos = props.horizontal ? e.clientX : e.clientY
  const delta = currentPos - startPos.value
  let newSize = startSize.value + delta
  
  // 限制范围
  newSize = Math.max(props.min, Math.min(props.max, newSize))
  
  size.value = newSize
  emit('update:size', newSize)
}

function stopResize() {
  isDragging.value = false
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
  document.body.style.cursor = ''
  document.body.style.userSelect = ''
}

onUnmounted(() => {
  document.removeEventListener('mousemove', handleResize)
  document.removeEventListener('mouseup', stopResize)
})

defineExpose({
  size
})
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
  background-color: #3E3E42;
  transition: background-color 0.2s, width 0.2s;
}

.resize-handle--horizontal .resize-handle-line {
  width: 1px;
  height: 40px;
  left: 50%;
  transform: translateX(-50%);
}

.resize-handle--hover .resize-handle-line,
.resize-handle--dragging .resize-handle-line {
  background-color: #0E639C;
  width: 2px;
}
</style>