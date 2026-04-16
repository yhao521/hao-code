<template>
  <div v-if="visible" class="command-palette-overlay" @click.self="close">
    <div class="command-palette">
      <input
        ref="inputRef"
        v-model="query"
        type="text"
        placeholder="输入命令..."
        class="palette-input"
        @keydown="handleKeydown"
      />
      <div class="palette-list">
        <div
          v-for="(cmd, index) in filteredCommands"
          :key="cmd.id"
          class="palette-item"
          :class="{ active: selectedIndex === index }"
          @click="execute(cmd)"
        >
          <span class="item-label">{{ cmd.label }}</span>
          <span class="item-shortcut" v-if="cmd.shortcut">{{
            cmd.shortcut
          }}</span>
        </div>
        <div v-if="filteredCommands.length === 0" class="no-results">
          无匹配结果
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick } from "vue";
import { lspManager } from "@/utils/lspManager";
import { useEditorStore } from "@/stores/editor";

const visible = ref(false);
const query = ref("");
const selectedIndex = ref(0);
const inputRef = ref<HTMLInputElement | null>(null);

interface Command {
  id: string;
  label: string;
  shortcut?: string;
  action: () => void;
}

const commands: Command[] = [
  {
    id: "file.new",
    label: "新建文件",
    shortcut: "Ctrl+N",
    action: () => console.log("New File"),
  },
  {
    id: "file.open",
    label: "打开文件...",
    shortcut: "Ctrl+O",
    action: () => window.dispatchEvent(new CustomEvent("menu:open-file")),
  },
  {
    id: "view.toggleSidebar",
    label: "切换侧边栏",
    shortcut: "Ctrl+B",
    action: () => window.dispatchEvent(new CustomEvent("view:toggle-sidebar")),
  },
  {
    id: "help.about",
    label: "关于 Hao-Code",
    action: () => alert("Hao-Code v1.0"),
  },
];

const editorStore = useEditorStore();
const workspaceSymbols = ref<Command[]>([]);

// 监听输入，如果以 @ 开头则触发工作区符号搜索
watch(query, async (newQuery) => {
  if (newQuery.startsWith("@") && newQuery.length > 1) {
    const searchQuery = newQuery.substring(1);
    if (searchQuery.length >= 2) {
      try {
        const symbols = await lspManager.getWorkspaceSymbols(searchQuery);
        workspaceSymbols.value = symbols.map((s: any) => ({
          id: `symbol-${s.name}-${s.location.uri}`,
          label: `${s.name} (${s.kind})`,
          action: () => {
            // 跳转到符号位置
            const uri = s.location.uri;
            const range = s.location.range;
            window.dispatchEvent(
              new CustomEvent("editor:jump-to-location", {
                detail: {
                  uri,
                  line: range.start.line + 1,
                  col: range.start.character + 1,
                },
              }),
            );
            close();
          },
        }));
      } catch (e) {
        console.error("Failed to fetch workspace symbols", e);
      }
    }
  } else {
    workspaceSymbols.value = [];
  }
});

const filteredCommands = computed(() => {
  if (workspaceSymbols.value.length > 0) return workspaceSymbols.value;
  if (!query.value) return commands;
  return commands.filter((cmd) =>
    cmd.label.toLowerCase().includes(query.value.toLowerCase()),
  );
});

function open() {
  visible.value = true;
  query.value = "";
  selectedIndex.value = 0;
  nextTick(() => inputRef.value?.focus());
}

function close() {
  visible.value = false;
}

function handleKeydown(e: KeyboardEvent) {
  if (e.key === "ArrowDown") {
    e.preventDefault();
    selectedIndex.value =
      (selectedIndex.value + 1) % filteredCommands.value.length;
  } else if (e.key === "ArrowUp") {
    e.preventDefault();
    selectedIndex.value =
      selectedIndex.value <= 0
        ? filteredCommands.value.length - 1
        : selectedIndex.value - 1;
  } else if (e.key === "Enter") {
    execute(filteredCommands.value[selectedIndex.value]);
  } else if (e.key === "Escape") {
    close();
  }
}

function execute(cmd?: Command) {
  if (cmd) {
    cmd.action();
    close();
  }
}

// 监听全局快捷键 Ctrl+Shift+P
window.addEventListener("keydown", (e) => {
  if ((e.ctrlKey || e.metaKey) && e.shiftKey && e.key === "p") {
    e.preventDefault();
    open();
  }
});

defineExpose({ open, close });
</script>

<style scoped>
.command-palette-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.4);
  z-index: 9999;
  display: flex;
  justify-content: center;
  padding-top: 10vh;
}

.command-palette {
  width: 600px;
  max-width: 90%;
  background-color: #252526;
  border-radius: 6px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.palette-input {
  width: 100%;
  padding: 12px 16px;
  background-color: #3c3c3c;
  border: none;
  color: #cccccc;
  font-size: 14px;
  outline: none;
}

.palette-list {
  max-height: 300px;
  overflow-y: auto;
}

.palette-item {
  padding: 8px 16px;
  display: flex;
  justify-content: space-between;
  cursor: pointer;
  color: #cccccc;
}

.palette-item.active {
  background-color: #094771;
  color: #ffffff;
}

.item-shortcut {
  color: #858585;
  font-size: 12px;
}

.no-results {
  padding: 12px 16px;
  color: #858585;
  text-align: center;
}
</style>
