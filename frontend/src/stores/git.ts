import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Repository {
  path: string
  currentBranch: string
}

export interface Branch {
  name: string
  fullName: string
  isRemote: boolean
  isCurrent: boolean
  lastCommit?: Commit
  ahead: number
  behind: number
}

export interface Commit {
  hash: string
  shortHash: string
  author: Signature
  message: string
  timestamp: string
}

export interface Signature {
  name: string
  email: string
}

export interface Change {
  path: string
  status: 'modified' | 'added' | 'deleted' | 'renamed'
  oldPath?: string
}

export const useGitStore = defineStore('git', () => {
  // State
  const repository = ref<Repository | null>(null)
  const branches = ref<Branch[]>([])
  const currentBranch = ref('')
  const changes = ref<Change[]>([])
  const stagedChanges = ref<Change[]>([])
  const commitMessage = ref('')
  const isCommitting = ref(false)
  const isLoading = ref(false)

  // Actions
  async function loadRepository(path: string) {
    isLoading.value = true
    try {
      // TODO: 调用后端 API
      console.log('Loading repository:', path)
      repository.value = {
        path,
        currentBranch: 'main'
      }
      currentBranch.value = 'main'
    } catch (error) {
      console.error('Failed to load repository:', error)
    } finally {
      isLoading.value = false
    }
  }

  async function fetchChanges() {
    // TODO: 调用后端 API
    changes.value = []
    stagedChanges.value = []
  }

  async function commit(message: string) {
    if (!message.trim()) {
      throw new Error('Commit message cannot be empty')
    }

    isCommitting.value = true
    try {
      // TODO: 调用后端 API
      console.log('Committing:', message)
      commitMessage.value = ''
      await fetchChanges()
    } catch (error) {
      console.error('Commit failed:', error)
      throw error
    } finally {
      isCommitting.value = false
    }
  }

  async function switchBranch(branchName: string) {
    isLoading.value = true
    try {
      // TODO: 调用后端 API
      console.log('Switching to branch:', branchName)
      currentBranch.value = branchName
    } catch (error) {
      console.error('Failed to switch branch:', error)
      throw error
    } finally {
      isLoading.value = false
    }
  }

  function stageFile(path: string) {
    // TODO: 实现暂存逻辑
    console.log('Staging:', path)
  }

  function unstageFile(path: string) {
    // TODO: 实现取消暂存逻辑
    console.log('Unstaging:', path)
  }

  return {
    // State
    repository,
    branches,
    currentBranch,
    changes,
    stagedChanges,
    commitMessage,
    isCommitting,
    isLoading,
    // Actions
    loadRepository,
    fetchChanges,
    commit,
    switchBranch,
    stageFile,
    unstageFile
  }
})
