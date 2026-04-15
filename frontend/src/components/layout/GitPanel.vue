<template>
  <div class="git-panel">
    <div class="git-header">
      <span class="git-title">源代码管理</span>
      <NButton text circle size="tiny" @click="refreshGit">
        <template #icon>
          <NIcon><RefreshOutline /></NIcon>
        </template>
      </NButton>
    </div>

    <NSpin :show="gitStore.isLoading">
      <div class="git-content" v-if="gitStore.repository">
        <!-- 分支信息 -->
        <div class="branch-info">
          <NIcon><GitBranchOutline /></NIcon>
          <span class="branch-name">{{ gitStore.currentBranch }}</span>
        </div>

        <!-- 暂存的更改 -->
        <div class="changes-section" v-if="gitStore.stagedChanges.length > 0">
          <div class="section-title">
            已暂存的更改 ({{ gitStore.stagedChanges.length }})
          </div>
          <NList hoverable size="small">
            <NListItem
              v-for="change in gitStore.stagedChanges"
              :key="change.path"
            >
              <div class="file-item">
                <span class="file-status" :class="change.status">{{
                  getStatusIcon(change.status)
                }}</span>
                <span class="file-path">{{ change.path }}</span>
              </div>
            </NListItem>
          </NList>
        </div>

        <!-- 未暂存的更改 -->
        <div class="changes-section" v-if="gitStore.changes.length > 0">
          <div class="section-title">更改 ({{ gitStore.changes.length }})</div>
          <NList hoverable size="small">
            <NListItem
              v-for="change in gitStore.changes"
              :key="change.path"
              @click="handleFileClick(change.path)"
            >
              <div class="file-item">
                <span class="file-status" :class="change.status">{{
                  getStatusIcon(change.status)
                }}</span>
                <span class="file-path">{{ change.path }}</span>
              </div>
            </NListItem>
          </NList>
        </div>

        <!-- 提交区域 -->
        <div class="commit-section" v-if="hasChanges">
          <NInput
            v-model:value="gitStore.commitMessage"
            type="textarea"
            placeholder="输入提交消息 (Ctrl+Enter 提交)"
            :autosize="{ minRows: 3, maxRows: 6 }"
            @keydown.ctrl.enter="handleCommit"
          />
          <NButton
            block
            type="primary"
            :loading="gitStore.isCommitting"
            :disabled="!gitStore.commitMessage.trim()"
            @click="handleCommit"
            style="margin-top: 8px"
          >
            提交
          </NButton>
        </div>

        <!-- 分支图谱 -->
        <div class="graph-section" v-if="graphNodes.length > 0">
          <div class="section-title">分支图谱</div>
          <GitGraph :nodes="graphNodes" />
        </div>

        <!-- 提交历史时间线 -->
        <div class="history-timeline" v-if="recentCommits.length > 0">
          <div class="section-title">提交历史 ({{ recentCommits.length }})</div>
          <div class="timeline-list">
            <div
              v-for="(commit, index) in recentCommits"
              :key="commit.hash"
              class="timeline-item"
              @click="handleCommitClick(commit)"
            >
              <div
                class="timeline-line"
                :style="{ backgroundColor: getCommitColor(index) }"
              ></div>
              <div class="timeline-content">
                <div class="commit-header">
                  <span class="commit-hash">{{ commit.shortHash }}</span>
                  <span class="commit-time">{{
                    formatTime(commit.timestamp)
                  }}</span>
                </div>
                <div class="commit-message">
                  {{ commit.message.split("\n")[0] }}
                </div>
                <div class="commit-author">{{ commit.author }}</div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="no-repo" v-else>
        <NEmpty description="未检测到 Git 仓库">
          <template #extra>
            <NButton size="small" @click="initRepository">初始化仓库</NButton>
          </template>
        </NEmpty>
      </div>
    </NSpin>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from "vue";
import { NIcon, NInput, NButton, NSpin, NEmpty, useMessage } from "naive-ui";
import { GitBranchOutline, RefreshOutline } from "@vicons/ionicons5";
import { useGitStore } from "@/stores/git";
import { useEditorStore } from "@/stores/editor";
import {
  OpenRepository,
  GetGitStatus,
  GitCommit,
  GitGetBranches,
  GitGetLog,
  GetGitGraph,
  GetFileDiff,
} from "@wails/backend/appservice";
import GitGraph from "./GitGraph.vue";

const gitStore = useGitStore();
const editorStore = useEditorStore();
const recentCommits = ref<any[]>([]);
const graphNodes = ref<any[]>([]);
const message = useMessage();

const hasChanges = computed(
  () => gitStore.changes.length > 0 || gitStore.stagedChanges.length > 0,
);

async function loadGitInfo() {
  const projectRoot = await import("@wails/backend/appservice").then((m) =>
    m.GetProjectRoot(),
  );

  try {
    gitStore.isLoading = true;

    // 打开仓库
    const repoInfo = await OpenRepository(projectRoot);
    if (repoInfo) {
      gitStore.repository = {
        path: repoInfo.path,
        currentBranch: repoInfo.currentBranch,
      };
      gitStore.currentBranch = repoInfo.currentBranch;

      // 获取状态
      await fetchGitStatus();

      // 获取分支列表
      const branches = await GitGetBranches(projectRoot);
      if (branches) {
        gitStore.branches = branches.local.map((name: string) => ({
          name,
          fullName: `refs/heads/${name}`,
          isRemote: false,
          isCurrent: name === repoInfo.currentBranch,
          ahead: 0,
          behind: 0,
        }));
      }

      // 获取提交日志
      const commits = await GitGetLog(projectRoot, 10);
      recentCommits.value = commits || [];

      // 获取图谱数据
      const nodes = await GetGitGraph(projectRoot, 50);
      graphNodes.value = nodes || [];
    }
  } catch (error) {
    console.log("Not a git repository or failed to open:", error);
    gitStore.repository = null;
  } finally {
    gitStore.isLoading = false;
  }
}

async function fetchGitStatus() {
  try {
    const status = await GetGitStatus(gitStore.repository!.path);
    if (status) {
      // 将后端的 string 类型转换为前端的字面量类型
      gitStore.changes = (status.changes || []).map((c: any) => ({
        ...c,
        status: c.status as "modified" | "added" | "deleted" | "renamed",
      }));
      gitStore.stagedChanges = (status.stagedChanges || []).map((c: any) => ({
        ...c,
        status: c.status as "modified" | "added" | "deleted" | "renamed",
      }));
    }
  } catch (error) {
    console.error("Failed to get git status:", error);
  }
}

async function handleCommit() {
  if (!gitStore.commitMessage.trim()) return;

  try {
    await gitStore.commit(gitStore.commitMessage);
    // 刷新状态
    await fetchGitStatus();
    // 刷新提交日志
    const commits = await GitGetLog(gitStore.repository!.path, 10);
    recentCommits.value = commits || [];
  } catch (error) {
    console.error("Commit failed:", error);
  }
}

async function refreshGit() {
  await loadGitInfo();
}

async function initRepository() {
  // TODO: 实现初始化仓库功能
  console.log("Init repository");
}

function getStatusIcon(status: string): string {
  const icons: Record<string, string> = {
    modified: "M",
    added: "A",
    deleted: "D",
    renamed: "R",
  };
  return icons[status] || "?";
}

function formatTime(timestamp: number): string {
  if (!timestamp) return "";
  const date = new Date(timestamp * 1000);
  const now = new Date();
  const diff = now.getTime() - date.getTime();

  if (diff < 60000) return "刚刚";
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`;
  return date.toLocaleDateString();
}

function getCommitColor(index: number): string {
  const colors = ["#E57373", "#64B5F6", "#81C784", "#FFD54F", "#BA68C8"];
  return colors[index % colors.length];
}

function handleCommitClick(commit: any) {
  message.info(`查看提交详情: ${commit.shortHash}`);
  // TODO: 实现点击跳转至 Diff 视图或 Commit 详情页
}

async function handleFileClick(filePath: string) {
  try {
    const projectRoot = await import("@wails/backend/appservice").then((m) =>
      m.GetProjectRoot(),
    );
    const diff = await GetFileDiff(projectRoot, filePath);
    if (diff) {
      editorStore.setDiffMode(true, {
        path: diff.path,
        oldContent: diff.oldContent,
        newContent: diff.newContent,
      });
    }
  } catch (error) {
    console.error("Failed to get file diff:", error);
    message.error("获取差异失败");
  }
}

onMounted(() => {
  loadGitInfo();
});
</script>

<style scoped>
.git-panel {
  padding: 8px;
  height: 100%;
  overflow-y: auto;
}

.git-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
  padding-bottom: 8px;
  border-bottom: 1px solid #3e3e42;
}

.git-title {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  color: #bbbbbb;
}

.branch-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
  padding: 6px;
  background-color: #2d2d30;
  border-radius: 4px;
}

.branch-name {
  font-size: 13px;
  font-weight: 500;
  color: #4ec9b0;
}

.section-title {
  font-size: 11px;
  font-weight: bold;
  color: #bbbbbb;
  margin: 12px 0 6px 0;
  text-transform: uppercase;
}

.file-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.file-status {
  font-weight: bold;
  font-size: 12px;
  width: 20px;
  text-align: center;
}

.file-status.modified {
  color: #e2c08d;
}

.file-status.added {
  color: #81b88b;
}

.file-status.deleted {
  color: #c74e39;
}

.file-path {
  font-size: 12px;
  color: #cccccc;
}

.commit-section {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #3e3e42;
}

.history-timeline {
  margin-top: 16px;
}

.timeline-list {
  position: relative;
  padding-left: 12px;
}

.timeline-item {
  position: relative;
  padding: 8px 0 8px 20px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.timeline-item:hover {
  background-color: #2a2d2e;
}

.timeline-line {
  position: absolute;
  left: 0;
  top: 0;
  bottom: 0;
  width: 2px;
}

.timeline-content {
  padding-left: 8px;
}

.commit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2px;
}

.commit-time {
  font-size: 11px;
  color: #858585;
}

.commit-hash {
  font-family: monospace;
  font-size: 11px;
  color: #569cd6;
  min-width: 50px;
}

.commit-message {
  font-size: 12px;
  color: #cccccc;
  margin: 2px 0;
}

.commit-author {
  font-size: 11px;
  color: #858585;
  margin-top: 2px;
}

.no-repo {
  text-align: center;
  padding: 40px 20px;
}
</style>
