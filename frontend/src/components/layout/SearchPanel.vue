<template>
  <div class="search-panel">
    <div class="search-header">
      <span class="search-title">搜索</span>
    </div>
    <NInput
      v-model:value="searchQuery"
      placeholder="Search files..."
      size="small"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { NInput } from 'naive-ui'

const searchQuery = ref('')
</script>

<style scoped>
.search-panel {
  padding: 8px;
}

.search-header {
  margin-bottom: 8px;
}

.search-title {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  color: #BBBBBB;
}
</style>
