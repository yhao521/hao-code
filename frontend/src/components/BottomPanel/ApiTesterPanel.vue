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
  </div>
</template>

<script setup lang="ts">
import { ref } from "vue";
import { NSelect, NInput, NButton, NTabs, NTabPane } from "naive-ui";
import { SendHTTPRequest } from "@wails/backend/appservice";

const method = ref("GET");
const url = ref("");
const headersText = ref("");
const body = ref("");
const isLoading = ref(false);
const response = ref<any>(null);

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
    // 解析 Headers
    const headers: Record<string, string> = {};
    headersText.value.split("\n").forEach((line) => {
      const parts = line.split(":");
      if (parts.length >= 2) {
        headers[parts[0].trim()] = parts.slice(1).join(":").trim();
      }
    });

    const req = {
      method: method.value,
      url: url.value,
      headers: headers,
      body: body.value,
    };

    response.value = await SendHTTPRequest(req);
  } catch (error) {
    console.error("Request failed:", error);
  } finally {
    isLoading.value = false;
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
</style>
