<template>
  <NModal
    v-model:show="visible"
    preset="card"
    style="width: 500px"
    title="提交详情"
  >
    <div class="blame-detail" v-if="blameData">
      <div class="author-section">
        <div class="avatar">{{ getInitials(blameData.author) }}</div>
        <div class="info">
          <div class="name">{{ blameData.author }}</div>
          <div class="time">{{ formatTime(blameData.timestamp) }}</div>
        </div>
      </div>

      <div class="commit-message">
        <h4>提交信息</h4>
        <p>{{ blameData.message }}</p>
      </div>

      <div class="hash-section">
        <span class="label">Commit Hash:</span>
        <code>{{ blameData.hash || "N/A" }}</code>
      </div>
    </div>
  </NModal>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NModal } from "naive-ui";

const visible = ref(false);
const blameData = ref<any>(null);

function show(data: any) {
  blameData.value = data;
  visible.value = true;
}

function getInitials(name: string): string {
  return name
    .split(" ")
    .map((n) => n[0])
    .join("")
    .toUpperCase()
    .slice(0, 2);
}

function formatTime(timestamp: string | number): string {
  if (!timestamp) return "";
  const date = new Date(
    typeof timestamp === "string" ? parseInt(timestamp) : timestamp * 1000,
  );
  return date.toLocaleString();
}

defineExpose({ show });
</script>

<style scoped>
.blame-detail {
  padding: 8px 0;
}

.author-section {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  background-color: #007acc;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: bold;
  font-size: 18px;
}

.info .name {
  font-size: 16px;
  font-weight: 500;
  color: #cccccc;
}

.info .time {
  font-size: 12px;
  color: #858585;
  margin-top: 4px;
}

.commit-message h4 {
  margin: 0 0 8px 0;
  font-size: 13px;
  color: #bbbbbb;
  text-transform: uppercase;
}

.commit-message p {
  margin: 0;
  font-size: 14px;
  color: #cccccc;
  line-height: 1.5;
  background: #252526;
  padding: 12px;
  border-radius: 4px;
}

.hash-section {
  margin-top: 16px;
  font-size: 12px;
  color: #858585;
}

.hash-section code {
  display: block;
  margin-top: 4px;
  padding: 4px 8px;
  background: #1e1e1e;
  border-radius: 3px;
  color: #569cd6;
  font-family: monospace;
}
</style>
