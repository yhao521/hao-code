import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface LayoutConfig {
  sidebarVisible: boolean
  rightPanelVisible: boolean
  bottomPanelVisible: boolean
  sidebarWidth: number
  rightPanelWidth: number
  bottomPanelHeight: number
  activeRightPanel: string
  activeBottomPanel: string
}

export const useLayoutStore = defineStore('layout', () => {
  const sidebarVisible = ref(true)
  const rightPanelVisible = ref(false)
  const bottomPanelVisible = ref(false)
  const sidebarWidth = ref(240)
  const rightPanelWidth = ref(300)
  const bottomPanelHeight = ref(200)
  
  // 右侧面板标签
  const activeRightPanel = ref('properties')
  const rightPanelTabs = [
    { id: 'properties', label: '属性', icon: '📋' },
    { id: 'preview', label: '预览', icon: '👁️' },
    { id: 'outline', label: '大纲', icon: '📑' }
  ]
  
  // 下面板标签
  const activeBottomPanel = ref('terminal')
  const bottomPanelTabs = [
    { id: 'terminal', label: '终端', icon: '💻' },
    { id: 'output', label: '输出', icon: '📄' },
    { id: 'problems', label: '问题', icon: '⚠️' },
    { id: 'debug', label: '调试控制台', icon: '🐛' }
  ]
  
  // 布局模式
  const layoutMode = ref<'default' | 'side-by-side' | 'bottom-focused'>('default')
  
  // 切换侧边栏
  function toggleSidebar() {
    sidebarVisible.value = !sidebarVisible.value
  }
  
  // 切换右侧面板
  function toggleRightPanel() {
    rightPanelVisible.value = !rightPanelVisible.value
  }
  
  // 切换下面板
  function toggleBottomPanel() {
    bottomPanelVisible.value = !bottomPanelVisible.value
  }
  
  // 设置右侧面板活动标签
  function setActiveRightPanel(tabId: string) {
    activeRightPanel.value = tabId
    if (!rightPanelVisible.value) {
      rightPanelVisible.value = true
    }
  }
  
  // 设置下面板活动标签
  function setActiveBottomPanel(tabId: string) {
    activeBottomPanel.value = tabId
    if (!bottomPanelVisible.value) {
      bottomPanelVisible.value = true
    }
  }
  
  // 设置布局模式
  function setLayoutMode(mode: 'default' | 'side-by-side' | 'bottom-focused') {
    layoutMode.value = mode
    
    // 根据不同模式调整面板可见性
    switch (mode) {
      case 'default':
        sidebarVisible.value = true
        rightPanelVisible.value = false
        bottomPanelVisible.value = false
        break
      case 'side-by-side':
        sidebarVisible.value = true
        rightPanelVisible.value = true
        bottomPanelVisible.value = false
        break
      case 'bottom-focused':
        sidebarVisible.value = true
        rightPanelVisible.value = false
        bottomPanelVisible.value = true
        break
    }
  }
  
  return {
    sidebarVisible,
    rightPanelVisible,
    bottomPanelVisible,
    sidebarWidth,
    rightPanelWidth,
    bottomPanelHeight,
    activeRightPanel,
    rightPanelTabs,
    activeBottomPanel,
    bottomPanelTabs,
    layoutMode,
    toggleSidebar,
    toggleRightPanel,
    toggleBottomPanel,
    setActiveRightPanel,
    setActiveBottomPanel,
    setLayoutMode
  }
})
