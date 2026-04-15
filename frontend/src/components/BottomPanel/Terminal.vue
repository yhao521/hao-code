<template>
  <div ref="terminalContainer" class="terminal-container"></div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from "vue";
import { Terminal } from "xterm";
import { FitAddon } from "xterm-addon-fit";
import { StartTerminal } from "@wails/backend/appservice.js";
import "xterm/css/xterm.css";

const terminalContainer = ref<HTMLElement | null>(null);
let terminal: Terminal | null = null;
let fitAddon: FitAddon | null = null;

onMounted(async () => {
  if (terminalContainer.value) {
    terminal = new Terminal({
      cursorBlink: true,
      fontSize: 14,
      fontFamily: "'Fira Code', 'Cascadia Code', Consolas, monospace",
      theme: {
        background: "#1e1e1e",
        foreground: "#cccccc",
        cursor: "#ffffff",
        selectionBackground: "#3a3d41",
      },
    });

    fitAddon = new FitAddon();
    terminal.loadAddon(fitAddon);
    terminal.open(terminalContainer.value);
    fitAddon.fit();

    // 启动后端终端会话
    try {
      const sessionId = await StartTerminal();
      terminal.writeln(
        `\x1b[1;32mTerminal session started: ${sessionId}\x1b[0m`,
      );

      // 模拟交互（实际应通过 WebSocket 或 EventStream）
      terminal.onData((data) => {
        // 这里应该发送数据到后端
        console.log("Sending to PTY:", data);
      });
    } catch (error) {
      terminal.writeln(`\x1b[1;31mFailed to start terminal: ${error}\x1b[0m`);
    }

    window.addEventListener("resize", handleResize);
  }
});

function handleResize() {
  fitAddon?.fit();
}

onUnmounted(() => {
  window.removeEventListener("resize", handleResize);
  terminal?.dispose();
});
</script>

<style scoped>
.terminal-container {
  width: 100%;
  height: 100%;
  padding: 8px;
  box-sizing: border-box;
}
</style>
