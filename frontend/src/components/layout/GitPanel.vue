<template>
  <div class="git-panel">
    <div class="git-header">
      <span class="git-title">源代码管理</span>
    </div>
    
    <div class="git-content" v-if="gitStore.repository">
      <div class="branch-info">
        <NIcon><GitBranchOutline /></NIcon>
        <span>{{ gitStore.currentBranch }}</span>
      </div>
      
      <div class="changes-section" v-if="gitStore.changes.length > 0">
        <div class="section-title">更改</div>
        <NList hoverable>
          <NListItem v-for="change in gitStore.changes" :key="change.path">
            <span class="file-status" :class="change.status">{{ change.status[0].toUpperCase() }}</span>
            {{ change.path }}
          </NListItem>
        </NList>
      </div>
      
      <div class="commit-section">
        <NInput
          v-model:value="gitStore.commitMessage"
          type="textarea"
          placeholder="Message (press Ctrl+Enter to commit)"
          :autosize="{ minRows: 3 }"
        />
        <NButton
          block
          type="primary"
          :loading="gitStore.isCommitting"
          :disabled="!gitStore.commitMessage.trim()"
          @click="handleCommit"
        >
          提交
        </NButton>
      </div>
    </div>
    
    <div class="no-repo" v-else>
      <p>未打开 Git 仓库</p>
      <NButton size="small" @click="openRepository">打开文件夹</NButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { NIcon, NList, NListItem, NInput, NButton } from 'naive-ui'
import { GitBranchOutline } from '@vicons/ionicons5'
import { useGitStore } from '@/stores/git'

const gitStore = useGitStore()

async function handleCommit() {
  try {
    await gitStore.commit(gitStore.commitMessage)
  } catch (error) {
    console.error('Commit failed:', error)
  }
}

function openRepository() {
  // TODO: 调用 Wails API 选择文件夹
  console.log('Open repository')
}
</script>

<style scoped>
.git-panel {
  padding: 8px;
}

.git-header {
  margin-bottom: 12px;
}

.git-title {
  font-size: 11px;
  font-weight: bold;
  text-transform: uppercase;
  color: #BBBBBB;
}

.branch-info {
  display: flex;
  align-items: center;
  gap: 6px;
  margin-bottom: 12px;
  font-size: 12px;
}

.file-status {
  font-weight: bold;
  margin-right: 8px;
}

.file-status.modified {
  color: #E2C08D;
}

.file-status.added {
  color: #81B88B;
}

.file-status.deleted {
  color: #C74E39;
}

.commit-section {
  margin-top: 12px;
}

.no-repo {
  text-align: center;
  padding: 20px;
  color: #888;
}
</style>
