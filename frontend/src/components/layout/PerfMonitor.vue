<template>
  <div class="perf-monitor" v-if="visible">
    <div class="perf-item" :class="{ 'perf-warn': fps < 30 }">
      <span class="label">FPS</span>
      <span class="value">{{ fps }}</span>
    </div>
    <div class="perf-item">
      <span class="label">Memory</span>
      <span class="value">{{ memoryUsage }} MB</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";

const visible = ref(true);
const fps = ref(60);
const memoryUsage = ref(0);

let lastTime = performance.now();
let frames = 0;
let animationFrameId: number;

function updateFPS() {
  const now = performance.now();
  frames++;

  if (now - lastTime >= 1000) {
    fps.value = Math.round((frames * 1000) / (now - lastTime));
    frames = 0;
    lastTime = now;
  }

  // 获取内存占用 (仅在 Chrome/Edge 等支持 performance.memory 的浏览器中有效)
  if ((performance as any).memory) {
    memoryUsage.value = Number(
      ((performance as any).memory.usedJSHeapSize / 1024 / 1024).toFixed(1),
    );
  }

  animationFrameId = requestAnimationFrame(updateFPS);
}

onMounted(() => {
  animationFrameId = requestAnimationFrame(updateFPS);
});

onUnmounted(() => {
  cancelAnimationFrame(animationFrameId);
});
</script>

<style scoped>
.perf-monitor {
  position: fixed;
  top: 10px;
  right: 10px;
  background-color: rgba(0, 0, 0, 0.8);
  color: #fff;
  padding: 8px 12px;
  border-radius: 4px;
  font-size: 12px;
  font-family: monospace;
  z-index: 9999;
  display: flex;
  gap: 15px;
  pointer-events: none; /* 允许点击穿透 */
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.perf-item {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.label {
  color: #888;
  font-size: 10px;
  margin-bottom: 2px;
}

.value {
  font-weight: bold;
}

.perf-warn .value {
  color: #ff4d4f;
}
</style>
