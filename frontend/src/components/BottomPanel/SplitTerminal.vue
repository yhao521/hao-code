<template>
  <div class="split-terminal-wrapper">
    <!-- 终端工具栏 -->
    <div class="terminal-toolbar">
      <div class="toolbar-left">
        <span class="title">TERMINAL</span>
      </div>
      <div class="toolbar-right">
        <NTooltip trigger="hover">
          <template #trigger>
            <NIcon class="tool-icon" @click="store.toggleSplitDirection()">
              <SplitOutline v-if="store.splitDirection === 'horizontal'" />
              <AppsOutline v-else />
            </NIcon>
          </template>
          Toggle Split Direction
        </NTooltip>
        <NTooltip trigger="hover">
          <template #trigger>
            <NIcon class="tool-icon" @click="handleAdd">
              <AddOutline />
            </NIcon>
          </template>
          New Terminal
        </NTooltip>
      </div>
    </div>

    <div class="split-terminal" :style="layoutStyle">
      <div
        v-for="term in store.instances"
        :key="term.id"
        class="terminal-instance"
        :class="{ active: term.id === store.activeId }"
        @click="store.setActive(term.id)"
      >
        <div class="term-header">
          <span>{{ term.name }}</span>
          <NIcon class="close-term" @click.stop="handleClose(term.id)"
            ><CloseOutline
          /></NIcon>
        </div>
        <div
          :ref="(el) => setTermRef(el, term.id)"
          class="xterm-container"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, nextTick } from "vue";
import { NIcon, NTooltip } from "naive-ui";
import {
  CloseOutline,
  AddOutline,
  SplitOutline,
  AppsOutline,
} from "@vicons/ionicons5";
import { useTerminalStore } from "@/stores/terminal";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import "xterm/css/xterm.css";

const store = useTerminalStore();
const termRefs = new Map<string, HTMLElement>();
const terminals = new Map<
  string,
  { xterm: Terminal; fit: FitAddon; id: string; ws: WebSocket }
>();

const layoutStyle = computed(() => ({
  flexDirection:
    store.splitDirection === "vertical"
      ? ("row" as const)
      : ("column" as const),
}));

function setTermRef(el: any, id: string) {
  if (el) termRefs.set(id, el);
}

async function initTerminalInstance(id: string) {
  const el = termRefs.get(id);
  if (!el || terminals.has(id)) return;

  const xterm = new Terminal({
    cursorBlink: true,
    fontSize: 13,
    fontFamily: "'Fira Code', Consolas, monospace",
    theme: { background: "#1e1e1e", foreground: "#cccccc" },
  });

  const fit = new FitAddon();
  xterm.loadAddon(fit);
  xterm.open(el);
  fit.fit();

  // 建立 WebSocket 连接
  const protocol = window.location.protocol === "https:" ? "wss:" : "ws:";
  const wsUrl = `${protocol}//${window.location.host}/ws/terminal`;
  const ws = new WebSocket(wsUrl);

  ws.onopen = () => {
    console.log(`Terminal ${id} connected`);
  };

  ws.onmessage = (event) => {
    xterm.write(new Uint8Array(event.data));
  };

  ws.onclose = () => {
    console.log(`Terminal ${id} disconnected`);
  };

  xterm.onData((data) => {
    if (ws.readyState === WebSocket.OPEN) {
      ws.send(data);
    }
  });

  terminals.set(id, { xterm, fit, id, ws });
}

async function handleAdd() {
  const id = await CreateTerminal();
  store.addInstance(id);
  nextTick(() => initTerminalInstance(id));
}

async function handleClose(id: string) {
  const t = terminals.get(id);
  if (t) {
    t.ws.close();
    t.xterm.dispose();
    terminals.delete(id);
  }
  store.removeInstance(id);
}

// 暴露给父组件用于初始化第一个终端
defineExpose({ handleAdd, initTerminalInstance });

onMounted(() => {
  if (store.instances.length === 0) {
    store.addInstance("term-1");
  }
  // 监听任务运行事件
  window.addEventListener("terminal:run", (e: any) => {
    const command = e.detail;
    // 确保有活跃的终端实例
    if (store.instances.length === 0) {
      store.addInstance("term-1");
    }

    // 找到当前活跃的终端并发送命令
    const activeId = store.activeId || store.instances[0].id;
    const termData = terminals.get(activeId);

    if (termData && termData.ws.readyState === WebSocket.OPEN) {
      // 在命令后添加换行符以模拟用户按下回车
      termData.ws.send(command + "\r");
      // 切换到该终端标签页
      store.setActive(activeId);
    }
  });
});
</script>

<style scoped>
.split-terminal-wrapper {
  display: flex;
  flex-direction: column;
  height: 100%;
  width: 100%;
}

.terminal-toolbar {
  height: 30px;
  background: #252526;
  border-bottom: 1px solid #1e1e1e;
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 8px;
}

.toolbar-left .title {
  font-size: 11px;
  font-weight: bold;
  color: #bbbbbb;
}

.toolbar-right {
  display: flex;
  gap: 8px;
}

.tool-icon {
  cursor: pointer;
  color: #cccccc;
  font-size: 16px;
}

.tool-icon:hover {
  color: #ffffff;
}

.split-terminal {
  display: flex;
  flex: 1;
  overflow: hidden;
}

.terminal-instance {
  flex: 1;
  display: flex;
  flex-direction: column;
  border: 1px solid #2b2b2b;
}

.term-header {
  background: #252526;
  padding: 4px 8px;
  font-size: 11px;
  display: flex;
  justify-content: space-between;
  align-items: center;
  cursor: pointer;
}

.close-term {
  opacity: 0;
  transition: opacity 0.2s;
}

.terminal-instance:hover .close-term {
  opacity: 1;
}

.xterm-container {
  flex: 1;
  padding: 8px;
}
</style>
