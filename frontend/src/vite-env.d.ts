/// <reference types="vite/client" />

declare module '*.vue' {
    import type {DefineComponent} from 'vue'
    const component: DefineComponent<{}, {}, any>
    export default component
}

// Wails 全局类型声明
interface WailsGoBackendApp {
  Greet(name: string): Promise<string>
  ReadFile(path: string): Promise<string>
  WriteFile(path: string, content: string): Promise<void>
  ListDir(path: string): Promise<any[]>
  GetProjectRoot(): Promise<string>
  OpenFolderDialog(): Promise<string>
  OpenFileDialog(): Promise<string>
  SaveFileDialog(): Promise<string>
  SetProjectRoot(path: string): Promise<void>
  CreateFile(path: string): Promise<void>
  CreateDirectory(path: string): Promise<void>
  DeleteFileOrDirectory(path: string): Promise<void>
  RenameFileOrDirectory(oldPath: string, newPath: string): Promise<void>
  MoveFileOrDirectory(sourcePath: string, targetPath: string): Promise<void>
  GetFileStats(path: string): Promise<any>
  SearchFiles(rootPath: string, keyword: string, maxResults: number): Promise<any[]>
  CopyFileOrDirectory(sourcePath: string, targetPath: string): Promise<void>
  IsTextFile(path: string): Promise<boolean>
  GetFileExtension(path: string): Promise<string>
  GetDirectoryTree(path: string, depth: number): Promise<any[]>
  BackupFile(path: string): Promise<void>
  TouchFile(path: string): Promise<void>
  OpenRepository(path: string): Promise<any>
  GetGitStatus(path: string): Promise<any>
  GitCommit(path: string, message: string): Promise<string>
  GitGetBranches(path: string): Promise<any>
  GitGetLog(path: string, maxCommits: number): Promise<any[]>
  GitRebase(path: string, upstream: string): Promise<string>
  GitCherryPick(path: string, commit: string): Promise<string>
  GitMerge(path: string, branch: string): Promise<string>
  GitReset(path: string, mode: string, target: string): Promise<string>
  GitStash(path: string, action: string, message: string): Promise<string>
  GitCheckout(path: string, branch: string): Promise<string>
  GitPull(path: string, remote: string, branch: string): Promise<string>
  GitPush(path: string, remote: string, branch: string): Promise<string>
}

interface WailsGoBackend {
  App: WailsGoBackendApp
}

interface WailsGo {
  backend: WailsGoBackend
}

declare interface Window {
  go: WailsGo
}
