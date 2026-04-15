declare module "@wails/backend/appservice.js" {
  export function StartTerminal(): Promise<string>;
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
  export function WriteFile(path: string, content: string): Promise<void>;
  export function OpenFileDialog(): Promise<string>;
  export function OpenFolderDialog(): Promise<string>;
  export function GetProjectRoot(): Promise<string>;
  // 可以根据需要添加更多导出函数的声明
}
