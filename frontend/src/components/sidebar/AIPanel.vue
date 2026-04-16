<template>
  <div class="ai-panel">
    <div class="chat-header">
      <span class="title">AI Assistant</span>
      <div class="header-actions">
        <NButton
          text
          circle
          size="tiny"
          @click="showSettings = true"
          title="AI 设置"
        >
          <template #icon
            ><NIcon><SettingsOutline /></NIcon
          ></template>
        </NButton>
        <NButton text circle size="tiny" @click="clearChat" title="清空对话">
          <template #icon
            ><NIcon><TrashOutline /></NIcon
          ></template>
        </NButton>
      </div>
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

      <!-- 文件建议列表 -->
      <div v-if="showFileSuggestions" class="file-suggestions">
        <div
          v-for="file in fileSuggestions"
          :key="file"
          class="suggestion-item"
          @click="selectFile(file)"
        >
          {{ file }}
        </div>
      </div>
    </div>
    <!-- 设置弹窗 -->
    <NModal v-model:show="showSettings" preset="dialog" title="AI 配置">
      <NForm :model="aiConfig" label-placement="left" label-width="80">
        <NFormItem label="API Key">
          <NInput
            v-model:value="aiConfig.apiKey"
            type="password"
            placeholder="sk-..."
          />
        </NFormItem>
        <NFormItem label="Base URL">
          <NInput
            v-model:value="aiConfig.baseURL"
            placeholder="https://api.openai.com/v1"
          />
        </NFormItem>
        <NFormItem label="Model">
          <NInput v-model:value="aiConfig.model" placeholder="gpt-3.5-turbo" />
        </NFormItem>
      </NForm>
      <template #footer>
        <NSpace>
          <NButton @click="showSettings = false">取消</NButton>
          <NButton type="primary" @click="saveAIConfig">保存</NButton>
        </NSpace>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, onMounted, watch } from "vue";
import {
  NButton,
  NIcon,
  NInput,
  NModal,
  NForm,
  NFormItem,
  NSpace,
  useMessage,
} from "naive-ui";
import { TrashOutline, CloseOutline, SettingsOutline } from "@vicons/ionicons5";
import { useEditorStore } from "@/stores/editor";
import {
  ChatWithAI,
  SetAIConfig,
  GetAIConfig,
  GetProjectRoot,
  ListDir,
  GetAIContextFromFiles,
} from "@wails/backend/appservice";

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
const showSettings = ref(false);
const message = useMessage();
const fileSuggestions = ref<string[]>([]);
const showFileSuggestions = ref(false);
const referencedFiles = ref<string[]>([]);

const aiConfig = ref({
  apiKey: "",
  baseURL: "https://api.openai.com/v1",
  model: "gpt-3.5-turbo",
});

// 监听输入框中的 @ 符号以触发文件建议
watch(inputValue, (val) => {
  const lastChar = val.slice(-1);
  if (lastChar === "@") {
    loadFileSuggestions();
    showFileSuggestions.value = true;
  } else if (!val.includes("@")) {
    showFileSuggestions.value = false;
  }
});

async function loadFileSuggestions() {
  try {
    const root = await GetProjectRoot();
    const files = await ListDir(root);
    fileSuggestions.value = files.map((f) => f.name);
  } catch (e) {
    console.error("Failed to load files", e);
  }
}

function selectFile(file: string) {
  inputValue.value = inputValue.value.slice(0, -1) + `@${file} `;
  if (!referencedFiles.value.includes(file)) {
    referencedFiles.value.push(file);
  }
  showFileSuggestions.value = false;
}

onMounted(async () => {
  try {
    const cfg = await GetAIConfig();
    aiConfig.value.baseURL = cfg.baseURL || aiConfig.value.baseURL;
    aiConfig.value.model = cfg.model || aiConfig.value.model;
  } catch (e) {
    console.error("Failed to load AI config", e);
  }
});

async function saveAIConfig() {
  try {
    await SetAIConfig(
      aiConfig.value.apiKey,
      aiConfig.value.baseURL,
      aiConfig.value.model,
    );
    message.success("AI 配置已保存");
    showSettings.value = false;
  } catch (error) {
    message.error("保存失败");
  }
}

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

    // 如果有 @file 引用，获取文件内容
    if (referencedFiles.value.length > 0) {
      const root = await GetProjectRoot();
      const fileContext = await GetAIContextFromFiles(
        root,
        referencedFiles.value,
      );
      context += "\n" + fileContext;
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

.file-suggestions {
  position: absolute;
  bottom: 100%;
  left: 0;
  right: 0;
  background: #252526;
  border: 1px solid #3e3e42;
  max-height: 150px;
  overflow-y: auto;
  z-index: 10;
}

.suggestion-item {
  padding: 4px 8px;
  cursor: pointer;
  font-size: 12px;
  color: #cccccc;
}

.suggestion-item:hover {
  background: #094771;
}
</style>
