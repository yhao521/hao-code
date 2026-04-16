// Type declarations for Wails bindings
// Support for @wails/runtime/runtime path (Wails v3 runtime)
declare module "@wails/runtime/runtime" {
  export * from "@wailsio/runtime";
}

// Support for @wails/go/backend/App path (used in most files)
declare module "@wails/go/backend/App" {
  import { SearchOptions, SearchResult } from "@wails/go/backend/models";

  export function AddRecentFile(path: string): Promise<void>;
  export function AddRecentFolder(path: string): Promise<void>;
  export function BackupFile(path: string): Promise<void>;
  export function ClearRecentFiles(): Promise<void>;
  export function ClearRecentFolders(): Promise<void>;
  export function CopyFileOrDirectory(
    sourcePath: string,
    targetPath: string,
  ): Promise<void>;
  export function CreateDirectory(path: string): Promise<void>;
  export function CreateFile(path: string): Promise<void>;
  export function DeleteFileOrDirectory(path: string): Promise<void>;
  export function GetDirectoryTree(path: string, depth: number): Promise<any[]>;
  export function GetFileBlame(path: string, filePath: string): Promise<any[]>;
  export function GetFileDiff(
    path: string,
    filePath: string,
  ): Promise<any | null>;
  export function GetFileExtension(path: string): Promise<string>;
  export function GetFileStats(path: string): Promise<any | null>;
  export function GetGitGraph(path: string, maxCommits: number): Promise<any[]>;
  export function GetGitStatus(path: string): Promise<any | null>;
  export function GetProjectRoot(): Promise<string>;
  export function GetRecentFiles(): Promise<any[]>;
  export function GetRecentFolders(): Promise<any[]>;
  export function GitCommit(path: string, message: string): Promise<string>;
  export function GitGetBranches(path: string): Promise<any | null>;
  export function GitGetLog(path: string, maxCommits: number): Promise<any[]>;
  export function Greet(name: string): Promise<string>;
  export function IsTextFile(path: string): Promise<boolean>;
  export function ListDir(path: string): Promise<any[]>;
  export function MoveFileOrDirectory(
    sourcePath: string,
    targetPath: string,
  ): Promise<void>;
  export function OpenFileDialog(): Promise<string>;
  export function OpenFolderDialog(): Promise<string>;
  export function OpenRepository(path: string): Promise<any | null>;
  export function ReadFile(path: string): Promise<string>;
  export function RemoveRecentFile(path: string): Promise<void>;
  export function RemoveRecentFolder(path: string): Promise<void>;
  export function RenameFileOrDirectory(
    oldPath: string,
    newPath: string,
  ): Promise<void>;
  export function RunTask(rootPath: string, command: string): Promise<void>;
  export function GetTasks(rootPath: string): Promise<any[]>;
  export function SaveFileDialog(): Promise<string>;
  export function SearchFiles(
    rootPath: string,
    keyword: string,
    maxResults: number,
  ): Promise<any[]>;
  export function SearchInFiles(opts: SearchOptions): Promise<SearchResult[]>;
  export function ScanTodos(rootPath: string): Promise<SearchResult[]>;
  export function SetProjectRoot(path: string): Promise<void>;
  export function TouchFile(path: string): Promise<void>;
  export function WriteFile(path: string, content: string): Promise<void>;
}

declare module "@wails/go/backend/models" {
  export class BranchInfo {
    local: string[];
    remote: string[];
    currentBranch: string;

    constructor(source?: Partial<BranchInfo>);
    static createFrom(source?: any): BranchInfo;
  }

  export class Change {
    path: string;
    status: string;
    oldPath?: string;

    constructor(source?: Partial<Change>);
    static createFrom(source?: any): Change;
  }

  export class CommitInfo {
    hash: string;
    shortHash: string;
    message: string;
    author: string;
    email: string;
    timestamp: string;

    constructor(source?: Partial<CommitInfo>);
    static createFrom(source?: any): CommitInfo;
  }

  export class FileInfo {
    name: string;
    path: string;
    size?: number;
    isDir: boolean;
    modTime?: number;

    constructor(source?: Partial<FileInfo>);
    static createFrom(source?: any): FileInfo;
  }

  export class GitGraphNode {
    hash: string;
    shortHash: string;
    message: string;
    author: string;
    timestamp: number;
    branches: string[];
    parents: string[];
    color: string;

    constructor(source?: Partial<GitGraphNode>);
    static createFrom(source?: any): GitGraphNode;
  }

  export class GitStatus {
    stagedChanges: Change[];
    changes: Change[];

    constructor(source?: Partial<GitStatus>);
    static createFrom(source?: any): GitStatus;
  }

  export class RecentItem {
    path: string;
    name: string;
    openedAt: string;

    constructor(source?: Partial<RecentItem>);
    static createFrom(source?: any): RecentItem;
  }

  export class RepoInfo {
    path: string;
    currentBranch: string;
    isRepository: boolean;

    constructor(source?: Partial<RepoInfo>);
    static createFrom(source?: any): RepoInfo;
  }

  export class SearchOptions {
    rootPath: string;
    query: string;
    caseSensitive: boolean;
    matchWholeWord: boolean;
    useRegex: boolean;
    exclude: string;

    constructor(source?: Partial<SearchOptions>);
    static createFrom(source?: any): SearchOptions;
  }

  export class SearchResult {
    filePath: string;
    lineNumber: number;
    lineContent: string;

    constructor(source?: Partial<SearchResult>);
    static createFrom(source?: any): SearchResult;
  }
}

// New path structure
declare module "@wails/backend/appservice" {
  import { SearchOptions, SearchResult } from "@wails/backend/models";

  export function AddRecentFile(path: string): Promise<void>;
  export function AddRecentFolder(path: string): Promise<void>;
  export function BackupFile(path: string): Promise<void>;
  export function ClearRecentFiles(): Promise<void>;
  export function ClearRecentFolders(): Promise<void>;
  export function CopyFileOrDirectory(
    sourcePath: string,
    targetPath: string,
  ): Promise<void>;
  export function CreateDirectory(path: string): Promise<void>;
  export function CreateFile(path: string): Promise<void>;
  export function DeleteFileOrDirectory(path: string): Promise<void>;
  export function GetDirectoryTree(path: string, depth: number): Promise<any[]>;
  export function GetFileBlame(path: string, filePath: string): Promise<any[]>;
  export function GetFileDiff(
    path: string,
    filePath: string,
  ): Promise<any | null>;
  export function GetFileExtension(path: string): Promise<string>;
  export function GetFileStats(path: string): Promise<any | null>;
  export function GetGitGraph(path: string, maxCommits: number): Promise<any[]>;
  export function GetGitStatus(path: string): Promise<any | null>;
  export function GetProjectRoot(): Promise<string>;
  export function GetRecentFiles(): Promise<any[]>;
  export function GetRecentFolders(): Promise<any[]>;
  export function GitCommit(path: string, message: string): Promise<string>;
  export function GitGetBranches(path: string): Promise<any | null>;
  export function GitGetLog(path: string, maxCommits: number): Promise<any[]>;
  export function Greet(name: string): Promise<string>;
  export function IsTextFile(path: string): Promise<boolean>;
  export function ListDir(path: string): Promise<any[]>;
  export function MoveFileOrDirectory(
    sourcePath: string,
    targetPath: string,
  ): Promise<void>;
  export function OpenFileDialog(): Promise<string>;
  export function OpenFolderDialog(): Promise<string>;
  export function OpenRepository(path: string): Promise<any | null>;
  export function ReadFile(path: string): Promise<string>;
  export function RemoveRecentFile(path: string): Promise<void>;
  export function RemoveRecentFolder(path: string): Promise<void>;
  export function RenameFileOrDirectory(
    oldPath: string,
    newPath: string,
  ): Promise<void>;
  export function RunTask(rootPath: string, command: string): Promise<void>;
  export function GetTasks(rootPath: string): Promise<any[]>;
  export function SaveFileDialog(): Promise<string>;
  export function SearchFiles(
    rootPath: string,
    keyword: string,
    maxResults: number,
  ): Promise<any[]>;
  export function SearchInFiles(opts: SearchOptions): Promise<SearchResult[]>;
  export function ScanTodos(rootPath: string): Promise<SearchResult[]>;
  export function SetProjectRoot(path: string): Promise<void>;
  export function TouchFile(path: string): Promise<void>;
  export function WriteFile(path: string, content: string): Promise<void>;
  export function ScanTodos(rootPath: string): Promise<SearchResult[]>;

  // AI Assistant methods
  export interface ChatMessage {
    role: string;
    content: string;
  }

  export interface ChatResponse {
    reply: string;
  }

  export interface AIConfig {
    apiKey: string;
    baseURL: string;
    model: string;
    maxTokens: number;
  }

  export function GetGhostText(
    prefix: string,
    suffix: string,
    language: string,
    filePath: string,
  ): Promise<{ text: string }>;

  // LSP Methods
  export function InitializeLSP(
    languageID: string,
    rootPath: string,
  ): Promise<void>;
  export function GetCompletions(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any[]>;
  export function GetDefinition(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any | null>;
  export function GetDocumentSymbols(
    languageID: string,
    uri: string,
  ): Promise<any[]>;
  export function FindReferences(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any[]>;
  export function RenameSymbol(
    languageID: string,
    uri: string,
    line: number,
    col: number,
    newName: string,
  ): Promise<any | null>;
  export function FormatDocument(
    languageID: string,
    uri: string,
    content: string,
  ): Promise<any[]>;
  export function GetDiagnostics(
    languageID: string,
    uri: string,
  ): Promise<any[]>;
  export function GetDiagnosticsCount(
    languageID: string,
    uri: string,
  ): Promise<{ errors: number; warnings: number }>;
  export function GetHoverInfo(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any | null>;
  export function GetSignatureHelp(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any | null>;
  export function GetCodeActions(
    languageID: string,
    uri: string,
    startLine: number,
    startCol: number,
    endLine: number,
    endCol: number,
    diagnostics: any[],
  ): Promise<any[]>;
  export function GetFoldingRanges(
    languageID: string,
    uri: string,
  ): Promise<any[]>;
  export function GetSemanticTokens(
    languageID: string,
    uri: string,
  ): Promise<any | null>;
  export function GetDocumentLinks(
    languageID: string,
    uri: string,
  ): Promise<any[]>;
  export function GetCodeLenses(
    languageID: string,
    uri: string,
  ): Promise<any[]>;
  export function PrepareCallHierarchy(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any[]>;
  export function GetIncomingCalls(
    languageID: string,
    item: any,
  ): Promise<any[]>;
  export function GetTypeHierarchy(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any[]>;
  export function GetImplementations(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any[]>;
  export function GetWorkspaceSymbols(query: string): Promise<any[]>;
  export function ResolveCodeAction(
    languageID: string,
    action: any,
  ): Promise<any | null>;

  export function SetAIConfig(
    apiKey: string,
    baseURL: string,
    model: string,
  ): Promise<void>;

  export function ChatWithAI(
    messages: ChatMessage[],
    context: string,
  ): Promise<ChatResponse>;

  export function GetAIConfig(): Promise<AIConfig>;

  // Git Staging methods
  export interface LineRange {
    start: number;
    end: number;
  }

  export function StageSelectedRanges(
    path: string,
    filePath: string,
    ranges: LineRange[],
  ): Promise<void>;

  export function UnstageFile(path: string, filePath: string): Promise<void>;

  // API Tester methods
  export interface APIRequest {
    method: string;
    url: string;
    headers: Record<string, string>;
    body: string;
  }

  export interface APIResponse {
    status: number;
    statusText: string;
    headers: Record<string, string>;
    body: string;
    duration: number;
  }

  export function SendHTTPRequest(req: APIRequest): Promise<APIResponse>;

  // API History methods
  export interface APIHistoryItem {
    id: string;
    timestamp: number;
    method: string;
    url: string;
    headers: Record<string, string>;
    body: string;
  }

  export function SaveApiHistory(req: APIRequest): Promise<void>;
  export function GetApiHistory(): Promise<APIHistoryItem[]>;
  export function DeleteApiHistory(id: string): Promise<void>;
}

declare module "@wails/backend/models" {
  export class BranchInfo {
    local: string[];
    remote: string[];
    currentBranch: string;

    constructor(source?: Partial<BranchInfo>);
    static createFrom(source?: any): BranchInfo;
  }

  export class Change {
    path: string;
    status: string;
    oldPath?: string;

    constructor(source?: Partial<Change>);
    static createFrom(source?: any): Change;
  }

  export class CommitInfo {
    hash: string;
    shortHash: string;
    message: string;
    author: string;
    email: string;
    timestamp: string;

    constructor(source?: Partial<CommitInfo>);
    static createFrom(source?: any): CommitInfo;
  }

  export class FileInfo {
    name: string;
    path: string;
    size?: number;
    isDir: boolean;
    modTime?: number;

    constructor(source?: Partial<FileInfo>);
    static createFrom(source?: any): FileInfo;
  }

  export class GitGraphNode {
    hash: string;
    shortHash: string;
    message: string;
    author: string;
    timestamp: number;
    branches: string[];
    parents: string[];
    color: string;

    constructor(source?: Partial<GitGraphNode>);
    static createFrom(source?: any): GitGraphNode;
  }

  export class GitStatus {
    stagedChanges: Change[];
    changes: Change[];

    constructor(source?: Partial<GitStatus>);
    static createFrom(source?: any): GitStatus;
  }

  export class RecentItem {
    path: string;
    name: string;
    openedAt: string;

    constructor(source?: Partial<RecentItem>);
    static createFrom(source?: any): RecentItem;
  }

  export class RepoInfo {
    path: string;
    currentBranch: string;
    isRepository: boolean;

    constructor(source?: Partial<RepoInfo>);
    static createFrom(source?: any): RepoInfo;
  }

  export class SearchOptions {
    rootPath: string;
    query: string;
    caseSensitive: boolean;
    matchWholeWord: boolean;
    useRegex: boolean;
    exclude: string;

    constructor(source?: Partial<SearchOptions>);
    static createFrom(source?: any): SearchOptions;
  }

  export class SearchResult {
    filePath: string;
    lineNumber: number;
    lineContent: string;

    constructor(source?: Partial<SearchResult>);
    static createFrom(source?: any): SearchResult;
  }
}
