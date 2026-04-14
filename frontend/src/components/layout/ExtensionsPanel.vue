<template>
  <div class="extensions-panel">
    <div class="extensions-header">
      <span class="extensions-title">扩展</span>
    </div>
    <NInput
      v-model:value="searchQuery"
      placeholder="Search extensions..."
      size="small"
    />
    <div class="extensions-list">
      <p class="placeholder">暂无扩展</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NInput } from 'naive-ui'

const searchQuery = ref('')
</script>

<style scoped>
.extensions-panel {
  padding: 8px;
}

.extensions-header {
  margin-bottom: 8px;
}

.extensions-title {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  color: #BBBBBB;
}

.placeholder {
  text-align: center;
  color: #666;
  margin-top: 40px;
}
</style>
