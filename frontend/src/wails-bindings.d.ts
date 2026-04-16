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
  export function SaveFileDialog(): Promise<string>;
  export function SearchFiles(
    rootPath: string,
    keyword: string,
    maxResults: number,
  ): Promise<any[]>;
  export function SearchInFiles(opts: SearchOptions): Promise<SearchResult[]>;
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
  export function SaveFileDialog(): Promise<string>;
  export function SearchFiles(
    rootPath: string,
    keyword: string,
    maxResults: number,
  ): Promise<any[]>;
  export function SearchInFiles(opts: SearchOptions): Promise<SearchResult[]>;
  export function SetProjectRoot(path: string): Promise<void>;
  export function TouchFile(path: string): Promise<void>;
  export function WriteFile(path: string, content: string): Promise<void>;
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
