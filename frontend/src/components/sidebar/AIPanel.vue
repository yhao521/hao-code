<template>
  <div class="ai-panel">
    <div class="chat-header">
      <span class="title">AI Assistant</span>
      <NButton text circle size="tiny" @click="clearChat" title="清空对话">
        <template #icon
          ><NIcon><TrashOutline /></NIcon
        ></template>
      </NButton>
    </div>

    <div class="chat-history" ref="historyRef">
      <div
        v-for="(msg, index) in messages"
        :key="index"
        :class="['message', msg.role]"
      >
        <div class="avatar">{{ msg.role === "user" ? "U" : "AI" }}</div>
        <div class="content">
          <pre v-if="msg.code">{{ msg.code }}</pre>
          <p>{{ msg.text }}</p>
        </div>
      </div>
      <div v-if="isLoading" class="message ai">
        <div class="avatar">AI</div>
        <div class="content typing">Thinking...</div>
      </div>
    </div>

    <div class="input-area">
      <div v-if="selectedCode" class="code-reference">
        <span>引用: {{ selectedCode.path }}:{{ selectedCode.line }}</span>
        <NIcon class="close-ref" @click="selectedCode = null"
          ><CloseOutline
        /></NIcon>
      </div>
      <NInput
        v-model:value="inputValue"
        type="textarea"
        placeholder="Ask AI anything..."
        :autosize="{ minRows: 3, maxRows: 6 }"
        @keydown.enter.ctrl="sendMessage"
      />
      <NButton
        block
        type="primary"
        class="send-btn"
        @click="sendMessage"
        :loading="isLoading"
      >
        Send
      </NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick } from "vue";
import { NButton, NIcon, NInput } from "naive-ui";
import { TrashOutline, CloseOutline } from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import { ChatWithAI } from "@wails/backend/appservice";

interface Message {
  role: "user" | "ai";
  text: string;
  code?: string;
}

const editorStore = useEditorStore();
const messages = ref<Message[]>([]);
const inputValue = ref("");
const isLoading = ref(false);
const historyRef = ref<HTMLElement>();
const selectedCode = ref<{
  path: string;
  line: number;
  content: string;
} | null>(null);

// 发送消息
async function sendMessage() {
  if (!inputValue.value.trim()) return;

  const userMsg = inputValue.value;
  messages.value.push({ role: "user", text: userMsg });
  inputValue.value = "";
  isLoading.value = true;

  try {
    // 准备上下文：如果选中了代码，则加入上下文
    let context = "";
    if (selectedCode.value) {
      context = selectedCode.value.content;
    }

    // 转换消息格式以适配后端
    const apiMessages = messages.value.map((m) => ({
      role: m.role === "user" ? "user" : "assistant",
      content: m.text + (m.code ? `\n\n\`\`\`\n${m.code}\n\`\`\`` : ""),
    }));

    const result = await ChatWithAI(apiMessages, context);

    if (result && result.reply) {
      messages.value.push({
        role: "ai",
        text: result.reply,
      });
    }
  } catch (error) {
    console.error("Chat error:", error);
    messages.value.push({
      role: "ai",
      text: "抱歉，AI 服务暂时不可用。请检查 API Key 配置。",
    });
  } finally {
    isLoading.value = false;
    scrollToBottom();
  }
}

function clearChat() {
  messages.value = [];
}

function scrollToBottom() {
  nextTick(() => {
    if (historyRef.value) {
      historyRef.value.scrollTop = historyRef.value.scrollHeight;
    }
  });
}
</script>

<style scoped>
.ai-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.chat-header {
  padding: 8px 12px;
  border-bottom: 1px solid #333;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.title {
  font-size: 12px;
  font-weight: bold;
  color: #cccccc;
}

.chat-history {
  flex: 1;
  overflow-y: auto;
  padding: 12px;
}

.message {
  display: flex;
  margin-bottom: 16px;
}

.message.user {
  flex-direction: row-reverse;
}

.avatar {
  width: 24px;
  height: 24px;
  border-radius: 4px;
  background: #007acc;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 10px;
  margin: 0 8px;
}

.content {
  background: #2d2d2d;
  padding: 8px 12px;
  border-radius: 6px;
  max-width: 85%;
  font-size: 13px;
  color: #d4d4d4;
}

.message.user .content {
  background: #0e639c;
}

pre {
  background: #1e1e1e;
  padding: 8px;
  border-radius: 4px;
  overflow-x: auto;
  margin-top: 4px;
  font-family: "Consolas", monospace;
}

.input-area {
  padding: 12px;
  border-top: 1px solid #333;
}

.code-reference {
  background: #2d2d2d;
  padding: 4px 8px;
  font-size: 11px;
  color: #858585;
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
  border-radius: 4px;
}

.close-ref {
  cursor: pointer;
}

.send-btn {
  margin-top: 8px;
}
</style>
