<template>
  <div class="api-tester-panel">
    <!-- 请求配置区 -->
    <div class="request-section">
      <div class="url-bar">
        <NSelect
          v-model:value="method"
          :options="methodOptions"
          style="width: 100px"
        />
        <NInput
          v-model:value="url"
          placeholder="Enter request URL"
          class="url-input"
        />
        <NButton type="primary" @click="sendRequest" :loading="isLoading"
          >Send</NButton
        >
        <NButton text @click="showEnvModal = true" title="管理环境变量">
          <template #icon
            ><NIcon><SettingsOutline /></NIcon
          ></template>
        </NButton>
      </div>

      <NTabs type="line" animated>
        <NTabPane name="headers" tab="Headers">
          <NInput
            v-model:value="headersText"
            type="textarea"
            placeholder="Key: Value (one per line)"
            :rows="4"
          />
        </NTabPane>
        <NTabPane name="body" tab="Body">
          <NInput
            v-model:value="body"
            type="textarea"
            placeholder="JSON Body"
            :rows="6"
          />
        </NTabPane>
      </NTabs>
    </div>

    <!-- 历史记录区 -->
    <div class="history-section" v-if="history.length > 0">
      <div class="section-title">最近请求</div>
      <NList hoverable size="small" style="background: transparent">
        <NListItem v-for="item in history" :key="item.id" class="history-item">
          <div class="history-content" @click="loadFromHistory(item)">
            <span class="history-method" :class="item.method.toLowerCase()">{{
              item.method
            }}</span>
            <span class="history-url">{{ item.url }}</span>
          </div>
          <template #suffix>
            <NIcon class="delete-icon" @click.stop="deleteHistoryItem(item.id)"
              ><CloseOutline
            /></NIcon>
          </template>
        </NListItem>
      </NList>
    </div>

    <!-- 响应展示区 -->
    <div class="response-section" v-if="response">
      <div class="response-header">
        <span
          :class="[
            'status-badge',
            response.status >= 400 ? 'error' : 'success',
          ]"
        >
          {{ response.status }} {{ response.statusText }}
        </span>
        <span class="duration">{{ response.duration }}ms</span>
      </div>
      <div class="response-body">
        <pre>{{ response.body }}</pre>
      </div>
    </div>

    <!-- 环境变量管理模态框 -->
    <NModal
      v-model:show="showEnvModal"
      preset="card"
      style="width: 600px"
      title="环境变量管理"
    >
      <div class="env-manager">
        <div v-for="(val, key) in envVars" :key="key" class="env-row">
          <NInput
            v-model:value="envVars[key]"
            :placeholder="key"
            size="small"
          />
          <NButton text type="error" @click="delete envVars[key]"
            ><NIcon><CloseOutline /></NIcon
          ></NButton>
        </div>
        <div class="add-env">
          <NInput
            v-model:value="newEnvKey"
            placeholder="变量名 (如 base_url)"
            size="small"
            style="width: 40%"
          />
          <NInput
            v-model:value="newEnvVal"
            placeholder="变量值"
            size="small"
            style="width: 50%"
          />
          <NButton size="small" @click="addEnvVar">添加</NButton>
        </div>
      </div>
      <template #footer>
        <NButton type="primary" @click="saveEnvVars">保存配置</NButton>
      </template>
    </NModal>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from "vue";
import {
  NSelect,
  NInput,
  NButton,
  NTabs,
  NTabPane,
  useMessage,
  NList,
  NListItem,
  NIcon,
} from "naive-ui";
import {
  SendHTTPRequest,
  GetApiHistory,
  DeleteApiHistory,
  GetEnvVariables,
  SaveEnvVariables,
} from "@wails/backend/appservice";
import { CloseOutline, SettingsOutline } from "@vicons/ionicons5";

const method = ref("GET");
const url = ref("");
const headersText = ref("");
const body = ref("");
const isLoading = ref(false);
const response = ref<any>(null);
const history = ref<any[]>([]);
const message = useMessage();
const envVars = ref<Record<string, string>>({});
const showEnvModal = ref(false);
const newEnvKey = ref("");
const newEnvVal = ref("");

const methodOptions = [
  { label: "GET", value: "GET" },
  { label: "POST", value: "POST" },
  { label: "PUT", value: "PUT" },
  { label: "DELETE", value: "DELETE" },
];

async function sendRequest() {
  if (!url.value) return;

  isLoading.value = true;
  try {
    // 1. 获取环境变量
    const vars = await GetEnvVariables();

    // 2. 替换 URL 中的变量
    let finalUrl = replaceVariables(url.value, vars);

    // 3. 替换 Headers 中的变量
    const headers: Record<string, string> = {};
    headersText.value.split("\n").forEach((line) => {
      const parts = line.split(":");
      if (parts.length >= 2) {
        const key = parts[0].trim();
        const val = replaceVariables(parts.slice(1).join(":").trim(), vars);
        headers[key] = val;
      }
    });

    // 4. 替换 Body 中的变量
    const finalBody = replaceVariables(body.value, vars);

    const req = {
      method: method.value,
      url: finalUrl,
      headers: headers,
      body: finalBody,
    };

    response.value = await SendHTTPRequest(req);
    message.success("请求发送成功");
  } catch (error) {
    console.error("Request failed:", error);
    message.error("请求失败");
  } finally {
    isLoading.value = false;
  }
}

function replaceVariables(str: string, vars: Record<string, string>): string {
  return str.replace(/\{\{(.*?)\}\}/g, (match, key) => {
    return vars[key.trim()] || match;
  });
}

async function loadHistory() {
  try {
    history.value = await GetApiHistory();
  } catch (e) {
    console.error("Failed to load history", e);
  }
}

function loadFromHistory(item: any) {
  method.value = item.method;
  url.value = item.url;
  body.value = item.body || "";

  // 转换 Headers 为文本格式
  const headerLines = [];
  if (item.headers) {
    for (const [key, value] of Object.entries(item.headers)) {
      headerLines.push(`${key}: ${value}`);
    }
  }
  headersText.value = headerLines.join("\n");
  message.info("已加载历史记录");
}

async function deleteHistoryItem(id: string) {
  await DeleteApiHistory(id);
  loadHistory();
}

onMounted(() => {
  loadHistory();
  loadEnvVars();
});

async function loadEnvVars() {
  try {
    envVars.value = await GetEnvVariables();
  } catch (e) {
    console.error("Failed to load env vars", e);
  }
}

async function saveEnvVars() {
  try {
    await SaveEnvVariables(envVars.value);
    message.success("环境变量已保存");
    showEnvModal.value = false;
  } catch (e) {
    message.error("保存失败");
  }
}

function addEnvVar() {
  if (newEnvKey.value && newEnvVal.value) {
    envVars.value[newEnvKey.value] = newEnvVal.value;
    newEnvKey.value = "";
    newEnvVal.value = "";
  }
}
</script>

<style scoped>
.api-tester-panel {
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: #1e1e1e;
}

.request-section {
  padding: 12px;
  border-bottom: 1px solid #333;
}

.url-bar {
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
}

.url-input {
  flex: 1;
}

.response-section {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.response-header {
  padding: 8px 12px;
  background: #252526;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.status-badge {
  font-size: 12px;
  font-weight: bold;
  padding: 2px 6px;
  border-radius: 3px;
}

.status-badge.success {
  color: #89d185;
}
.status-badge.error {
  color: #f48771;
}

.duration {
  font-size: 11px;
  color: #858585;
}

.response-body {
  flex: 1;
  overflow: auto;
  padding: 12px;
}

.response-body pre {
  margin: 0;
  font-family: "Fira Code", monospace;
  font-size: 12px;
  color: #cccccc;
}

.history-section {
  padding: 12px;
  border-top: 1px solid #333;
  max-height: 200px;
  overflow-y: auto;
}

.section-title {
  font-size: 11px;
  font-weight: bold;
  color: #bbbbbb;
  margin-bottom: 8px;
  text-transform: uppercase;
}

.history-item {
  background: transparent !important;
  padding: 4px 0;
  cursor: pointer;
}

.history-content {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
}

.history-method {
  font-weight: bold;
  width: 50px;
  text-align: center;
  padding: 2px 4px;
  border-radius: 3px;
  font-size: 10px;
}

.history-method.get {
  color: #61affe;
  background: rgba(97, 175, 254, 0.1);
}
.history-method.post {
  color: #49cc90;
  background: rgba(73, 204, 144, 0.1);
}
.history-method.put {
  color: #fca130;
  background: rgba(252, 161, 48, 0.1);
}
.history-method.delete {
  color: #f93e3e;
  background: rgba(249, 62, 62, 0.1);
}

.history-url {
  color: #cccccc;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.delete-icon {
  color: #858585;
  cursor: pointer;
  opacity: 0;
  transition: opacity 0.2s;
}

.history-item:hover .delete-icon {
  opacity: 1;
}

.env-manager {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.env-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.add-env {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #333;
}
</style>
