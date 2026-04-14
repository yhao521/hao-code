import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export interface Tab {
  id: string
  path: string
  name: string
  content?: string
  dirty: boolean
  language?: string
}

export type SidebarView = 'explorer' | 'search' | 'git' | 'extensions'

export const useEditorStore = defineStore('editor', () => {
  // State
  const activeEditor = ref<string | null>(null)
  const tabs = ref<Tab[]>([])
  const sidebarVisible = ref(true)
  const sidebarView = ref<SidebarView>('explorer')

  // Getters
  const activeTab = computed(() => 
    tabs.value.find(t => t.id === activeEditor.value)
  )
  
  const dirtyTabs = computed(() => 
    tabs.value.filter(t => t.dirty)
  )

  // Actions
  function openFile(path: string, content: string = '') {
    const existingTab = tabs.value.find(t => t.path === path)
    
    if (existingTab) {
      activeEditor.value = existingTab.id
      return
    }

    const tab: Tab = {
      id: Date.now().toString(),
      path,
      name: path.split('/').pop() || path,
      content,
      dirty: false,
      language: getLanguageFromPath(path)
    }

    tabs.value.push(tab)
    activeEditor.value = tab.id
  }

  function closeTab(id: string) {
    const index = tabs.value.findIndex(t => t.id === id)
    if (index === -1) return

    const tab = tabs.value[index]
    if (tab.dirty) {
      // TODO: 提示保存
      console.warn('File has unsaved changes:', tab.path)
      return
    }

    tabs.value.splice(index, 1)
    
    if (activeEditor.value === id) {
      activeEditor.value = tabs.value[Math.max(0, index - 1)]?.id || null
    }
  }

  function updateContent(id: string, content: string) {
    const tab = tabs.value.find(t => t.id === id)
    if (tab) {
      tab.content = content
      tab.dirty = true
    }
  }

  function saveFile(id: string) {
    const tab = tabs.value.find(t => t.id === id)
    if (tab) {
      tab.dirty = false
    }
  }

  function toggleSidebar() {
    sidebarVisible.value = !sidebarVisible.value
  }

  function setSidebarView(view: SidebarView) {
    sidebarView.value = view
    sidebarVisible.value = true
  }

  function getLanguageFromPath(path: string): string {
    const ext = path.split('.').pop()?.toLowerCase()
    const langMap: Record<string, string> = {
      'ts': 'typescript',
      'js': 'javascript',
      'tsx': 'typescriptreact',
      'jsx': 'javascriptreact',
      'py': 'python',
      'go': 'go',
      'java': 'java',
      'cpp': 'cpp',
      'c': 'c',
      'rs': 'rust',
      'rb': 'ruby',
      'php': 'php',
      'html': 'html',
      'css': 'css',
      'json': 'json',
      'md': 'markdown',
      'yaml': 'yaml',
      'xml': 'xml'
    }
    return langMap[ext || ''] || 'plaintext'
  }

  return {
    // State
    activeEditor,
    tabs,
    sidebarVisible,
    sidebarView,
    // Getters
    activeTab,
    dirtyTabs,
    // Actions
    openFile,
    closeTab,
    updateContent,
    saveFile,
    toggleSidebar,
    setSidebarView
  }
})
