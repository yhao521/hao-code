<template>
  <div class="git-graph-container" ref="containerRef">
    <canvas ref="canvasRef"></canvas>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from "vue";

interface GitGraphNode {
  hash: string;
  shortHash: string;
  message: string;
  author: string;
  timestamp: number;
  branches: string[];
  parents: string[];
  color: string;
}

const props = defineProps<{
  nodes: GitGraphNode[];
}>();

const containerRef = ref<HTMLElement | null>(null);
const canvasRef = ref<HTMLCanvasElement | null>(null);

const NODE_RADIUS = 6;
const LINE_WIDTH = 2;
const VERTICAL_SPACING = 30;
const HORIZONTAL_SPACING = 20;

function drawGraph() {
  const canvas = canvasRef.value;
  const ctx = canvas?.getContext("2d");
  if (!canvas || !ctx || props.nodes.length === 0) return;

  // 计算画布大小
  const width = containerRef.value?.clientWidth || 300;
  const height = props.nodes.length * VERTICAL_SPACING + 20;

  canvas.width = width;
  canvas.height = height;

  ctx.clearRect(0, 0, width, height);
  ctx.font = "12px -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto";

  // 简单的分支轨道分配逻辑
  const branchTracks: Record<string, number> = {};
  let nextTrack = 0;

  props.nodes.forEach((node, index) => {
    const y = index * VERTICAL_SPACING + 15;

    // 确定当前节点的轨道
    let track = 0;
    if (node.branches.length > 0) {
      const branch = node.branches[0];
      if (branchTracks[branch] === undefined) {
        branchTracks[branch] = nextTrack++;
      }
      track = branchTracks[branch];
    }

    const x = 30 + track * HORIZONTAL_SPACING;

    // 绘制连线到父节点
    ctx.lineWidth = LINE_WIDTH;
    node.parents.forEach((parentHash: string) => {
      const parentIndex = props.nodes.findIndex((n) => n.hash === parentHash);
      if (parentIndex !== -1) {
        const parentY = parentIndex * VERTICAL_SPACING + 15;
        // 简单起见，这里假设父节点也在同一轨道或寻找父节点轨道
        // 实际 Idea 级图谱需要复杂的拓扑排序
        let parentTrack = 0;
        // 简化处理：直接连到正上方
        ctx.strokeStyle = node.color;
        ctx.beginPath();
        ctx.moveTo(x, y - NODE_RADIUS);
        ctx.lineTo(x, parentY + NODE_RADIUS);
        ctx.stroke();
      }
    });

    // 绘制节点
    ctx.beginPath();
    ctx.arc(x, y, NODE_RADIUS, 0, Math.PI * 2);
    ctx.fillStyle = node.color;
    ctx.fill();
    ctx.strokeStyle = "#1e1e1e";
    ctx.stroke();

    // 绘制提交信息
    ctx.fillStyle = "#cccccc";
    ctx.fillText(`${node.shortHash} - ${node.message}`, x + 15, y + 4);
  });
}

onMounted(() => {
  drawGraph();
});

watch(() => props.nodes, drawGraph, { deep: true });
</script>

<style scoped>
.git-graph-container {
  width: 100%;
  overflow-x: auto;
  background-color: #1e1e1e;
}

canvas {
  display: block;
}
</style>
