<template>
  <div class="type-node" :style="{ paddingLeft: level * 12 + 'px' }">
    <div class="node-content" @dblclick="handleNavigate">
      <span class="icon" v-if="hasChildren" @click.stop="toggle">
        {{ expanded ? '▼' : '▶' }}
      </span>
      <span class="icon" v-else style="opacity: 0.3">●</span>
      <span class="name">{{ item.name }}</span>
      <span class="kind" v-if="item.kind">{{ getKindLabel(item.kind) }}</span>
    </div>
    
    <div class="children" v-if="expanded && hasChildren">
      <TypeHierarchyNode 
        v-for="(child, index) in childrenItems" 
        :key="index"
        :item="child" 
        :level="level + 1"
        @navigate="$emit('navigate', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { lspManager } from '@/utils/lspManager';

const props = defineProps<{
  item: any;
  level: number;
}>();

const emit = defineEmits(['navigate']);

const expanded = ref(true);
const childrenItems = ref<any[]>([]);

const hasChildren = computed(() => {
  return props.item.children && props.item.children.length > 0;
});

async function toggle() {
  if (!expanded.value && !childrenItems.value.length) {
    // 这里可以添加懒加载逻辑，如果 LSP 支持的话
    childrenItems.value = props.item.children || [];
  }
  expanded.value = !expanded.value;
}

function handleNavigate() {
  emit('navigate', props.item);
}

function getKindLabel(kind: number): string {
  const kinds = ['File', 'Module', 'Namespace', 'Package', 'Class', 'Method', 'Property', 'Field', 'Constructor', 'Enum', 'Interface', 'Function', 'Variable', 'Constant', 'String', 'Number', 'Boolean', 'Array', 'Object', 'Key', 'Null', 'EnumMember', 'Struct', 'Event', 'Operator', 'TypeParameter'];
  return kinds[kind] || '';
}
</script>

<style scoped>
.type-node {
  user-select: none;
}

.node-content {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  cursor: pointer;
  font-size: 13px;
}

.node-content:hover {
  background-color: #2a2d2e;
}

.icon {
  width: 16px;
  text-align: center;
  margin-right: 4px;
  font-size: 10px;
  color: #cccccc;
}

.name {
  color: #dcdcaa;
  margin-right: 6px;
}

.kind {
  font-size: 11px;
  color: #858585;
  background-color: #3c3c3c;
  padding: 1px 4px;
  border-radius: 3px;
}

.children {
  border-left: 1px solid #3c3c3c;
  margin-left: 7px;
}
</style>
