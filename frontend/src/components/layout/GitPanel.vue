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
          <div class="section-title">已暂存的更改 ({{ gitStore.stagedChanges.length }})</div>
          <NList hoverable size="small">
            <NListItem v-for="change in gitStore.stagedChanges" :key="change.path">
              <div class="file-item">
                <span class="file-status" :class="change.status">{{ getStatusIcon(change.status) }}</span>
                <span class="file-path">{{ change.path }}</span>
              </div>
            </NListItem>
          </NList>
        </div>
        
        <!-- 未暂存的更改 -->
        <div class="changes-section" v-if="gitStore.changes.length > 0">
          <div class="section-title">更改 ({{ gitStore.changes.length }})</div>
          <NList hoverable size="small">
            <NListItem v-for="change in gitStore.changes" :key="change.path">
              <div class="file-item">
                <span class="file-status" :class="change.status">{{ getStatusIcon(change.status) }}</span>
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
        
        <!-- 最近提交 -->
        <div class="recent-commits" v-if="recentCommits.length > 0">
          <div class="section-title">最近提交</div>
          <NList hoverable size="small">
            <NListItem v-for="commit in recentCommits" :key="commit.hash">
              <div class="commit-item">
                <span class="commit-hash">{{ commit.shortHash }}</span>
                <span class="commit-message">{{ commit.message.split('\n')[0] }}</span>
              </div>
            </NListItem>
          </NList>
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
import { ref, computed, onMounted } from 'vue'
import { 
  NIcon, NList, NListItem, NInput, NButton, NSpin, NEmpty 
} from 'naive-ui'
import { GitBranchOutline, RefreshOutline } from '@vicons/ionicons5'
import { useGitStore } from '@/stores/git'
import { 
  OpenRepository,
  GetGitStatus, 
  GitCommit,
  GitGetBranches,
  GitGetLog 
} from '@wails/go/backend/App'

const gitStore = useGitStore()
const recentCommits = ref<any[]>([])

const hasChanges = computed(() => 
  gitStore.changes.length > 0 || gitStore.stagedChanges.length > 0
)

async function loadGitInfo() {
  const projectRoot = await import('@wails/go/backend/App').then(m => m.GetProjectRoot())
  
  try {
    gitStore.isLoading = true
    
    // 打开仓库
    const repoInfo = await OpenRepository(projectRoot)
    if (repoInfo) {
      gitStore.repository = {
        path: repoInfo.path,
        currentBranch: repoInfo.currentBranch
      }
      gitStore.currentBranch = repoInfo.currentBranch
      
      // 获取状态
      await fetchGitStatus()
      
      // 获取分支列表
      const branches = await GitGetBranches(projectRoot)
      if (branches) {
        gitStore.branches = branches.local.map((name: string) => ({
          name,
          fullName: `refs/heads/${name}`,
          isRemote: false,
          isCurrent: name === repoInfo.currentBranch,
          ahead: 0,
          behind: 0
        }))
      }
      
      // 获取提交日志
      const commits = await GitGetLog(projectRoot, 10)
      recentCommits.value = commits || []
    }
  } catch (error) {
    console.log('Not a git repository or failed to open:', error)
    gitStore.repository = null
  } finally {
    gitStore.isLoading = false
  }
}

async function fetchGitStatus() {
  try {
    const status = await GetGitStatus(gitStore.repository!.path)
    if (status) {
      // 将后端的 string 类型转换为前端的字面量类型
      gitStore.changes = (status.changes || []).map((c: any) => ({
        ...c,
        status: c.status as 'modified' | 'added' | 'deleted' | 'renamed'
      }))
      gitStore.stagedChanges = (status.stagedChanges || []).map((c: any) => ({
        ...c,
        status: c.status as 'modified' | 'added' | 'deleted' | 'renamed'
      }))
    }
  } catch (error) {
    console.error('Failed to get git status:', error)
  }
}

async function handleCommit() {
  if (!gitStore.commitMessage.trim()) return
  
  try {
    await gitStore.commit(gitStore.commitMessage)
    // 刷新状态
    await fetchGitStatus()
    // 刷新提交日志
    const commits = await GitGetLog(gitStore.repository!.path, 10)
    recentCommits.value = commits || []
  } catch (error) {
    console.error('Commit failed:', error)
  }
}

async function refreshGit() {
  await loadGitInfo()
}

async function initRepository() {
  // TODO: 实现初始化仓库功能
  console.log('Init repository')
}

function getStatusIcon(status: string): string {
  const icons: Record<string, string> = {
    'modified': 'M',
    'added': 'A',
    'deleted': 'D',
    'renamed': 'R'
  }
  return icons[status] || '?'
}

onMounted(() => {
  loadGitInfo()
})
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
  border-bottom: 1px solid #3E3E42;
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
  padding: 6px;
  background-color: #2D2D30;
  border-radius: 4px;
}

.branch-name {
  font-size: 13px;
  font-weight: 500;
  color: #4EC9B0;
}

.section-title {
  font-size: 11px;
  font-weight: bold;
  color: #BBBBBB;
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
  color: #E2C08D;
}

.file-status.added {
  color: #81B88B;
}

.file-status.deleted {
  color: #C74E39;
}

.file-path {
  font-size: 12px;
  color: #CCCCCC;
}

.commit-section {
  margin-top: 16px;
  padding-top: 12px;
  border-top: 1px solid #3E3E42;
}

.recent-commits {
  margin-top: 16px;
}

.commit-item {
  display: flex;
  gap: 8px;
  align-items: center;
}

.commit-hash {
  font-family: monospace;
  font-size: 11px;
  color: #569CD6;
  min-width: 50px;
}

.commit-message {
  font-size: 12px;
  color: #CCCCCC;
  flex: 1;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.no-repo {
  text-align: center;
  padding: 40px 20px;
}
</style>
