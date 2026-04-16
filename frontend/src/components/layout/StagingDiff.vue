<template>
  <div class="staging-diff">
    <div class="diff-header">
      <span class="file-name">{{ diff.path }}</span>
      <div class="actions">
        <NButton size="tiny" @click="stageAll">暂存所有更改</NButton>
        <NButton
          size="tiny"
          type="primary"
          :disabled="selectedLines.length === 0"
          @click="stageSelected"
        >
          暂存选中行 ({{ selectedLines.length }})
        </NButton>
      </div>
    </div>

    <div class="diff-content">
      <div
        v-for="(line, index) in diffLines"
        :key="index"
        class="diff-line"
        :class="{ selected: selectedLines.includes(index) }"
        @click="toggleLine(index)"
      >
        <div class="line-checkbox">
          <input
            type="checkbox"
            :checked="selectedLines.includes(index)"
            @click.stop
          />
        </div>
        <div class="line-number">{{ line.oldNum || "-" }}</div>
        <div class="line-number">{{ line.newNum || "-" }}</div>
        <div class="line-text" :class="line.type">{{ line.content }}</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from "vue";
import { NButton, NCheckbox } from "naive-ui";
import { StageSelectedRanges, GetProjectRoot } from "@wails/backend/appservice";

const props = defineProps<{
  diff: { path: string; oldContent: string; newContent: string };
}>();

const selectedLines = ref<number[]>([]);

// 简单的 Diff 解析逻辑（实际应使用 diff-match-patch 等库）
const diffLines = computed(() => {
  const oldLines = props.diff.oldContent.split("\n");
  const newLines = props.diff.newContent.split("\n");
  const lines: any[] = [];

  // 这里简化处理：仅展示新内容并标记为添加/修改
  // 真正的逐行暂存需要复杂的 Diff 算法来映射行号
  newLines.forEach((content, i) => {
    lines.push({
      oldNum: i + 1,
      newNum: i + 1,
      content: content,
      type: "new", // 简化为全部视为新行以便演示 UI
    });
  });

  return lines;
});

function toggleLine(index: number, checked?: boolean) {
  if (checked === undefined) {
    const idx = selectedLines.value.indexOf(index);
    if (idx > -1) selectedLines.value.splice(idx, 1);
    else selectedLines.value.push(index);
  } else {
    if (checked && !selectedLines.value.includes(index))
      selectedLines.value.push(index);
    if (!checked)
      selectedLines.value = selectedLines.value.filter((i) => i !== index);
  }
}

async function stageSelected() {
  if (selectedLines.value.length === 0) return;

  const root = await GetProjectRoot();
  const ranges = selectedLines.value.map((lineIndex) => ({
    start: lineIndex + 1,
    end: lineIndex + 1,
  }));

  try {
    await StageSelectedRanges(root, props.diff.path, ranges);
    selectedLines.value = [];
    // 触发刷新事件
    window.dispatchEvent(new CustomEvent("git:refresh"));
  } catch (e) {
    console.error("Stage failed", e);
  }
}

async function stageAll() {
  const root = await GetProjectRoot();
  // 暂存所有行
  const ranges = diffLines.value.map((_, i) => ({ start: i + 1, end: i + 1 }));
  await StageSelectedRanges(root, props.diff.path, ranges);
  window.dispatchEvent(new CustomEvent("git:refresh"));
}
</script>

<style scoped>
.staging-diff {
  border: 1px solid #3e3e42;
  background: #1e1e1e;
}

.diff-header {
  padding: 8px;
  background: #252526;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.file-name {
  font-size: 12px;
  color: #cccccc;
}

.actions {
  display: flex;
  gap: 8px;
}

.diff-content {
  max-height: 400px;
  overflow-y: auto;
}

.diff-line {
  display: flex;
  align-items: center;
  font-family: monospace;
  font-size: 12px;
  cursor: pointer;
}

.diff-line.selected {
  background-color: rgba(46, 160, 67, 0.2);
}

.diff-line:hover {
  background-color: #2a2d2e;
}

.line-checkbox {
  width: 24px;
  display: flex;
  justify-content: center;
  align-items: center;
}

.line-checkbox input {
  cursor: pointer;
}

.line-number {
  width: 40px;
  text-align: right;
  padding-right: 8px;
  color: #858585;
  user-select: none;
}

.line-text {
  flex: 1;
  white-space: pre;
  padding-left: 8px;
}

.line-text.added {
  color: #81b88b;
}
.line-text.deleted {
  color: #c74e39;
}
</style>
