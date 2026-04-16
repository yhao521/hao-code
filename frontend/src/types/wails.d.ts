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
  export function GetDefinition(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any>;
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
  ): Promise<any>;
  export function FormatDocument(
    languageID: string,
    uri: string,
    content: string,
  ): Promise<any[]>;
  export function GetHoverInfo(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any>;
  export function GetSignatureHelp(
    languageID: string,
    uri: string,
    line: number,
    col: number,
  ): Promise<any>;
  export interface PluginManifest {
    name: string;
    version: string;
    description: string;
    main: string;
    author: string;
    license: string;
    contributes: any;
  }
  export function GetInstalledPlugins(): Promise<PluginManifest[]>;
  export function ActivatePlugin(name: string): Promise<void>;
  export function GetCodeActions(
    languageID: string,
    uri: string,
    startLine: number,
    startCol: number,
    endLine: number,
    endCol: number,
    diagnostics: any[]
  ): Promise<any[]>;
  export function GetFoldingRanges(languageID: string, uri: string): Promise<any[]>;
  export function WriteFile(path: string, content: string): Promise<void>;
  export function OpenFileDialog(): Promise<string>;
  export function OpenFolderDialog(): Promise<string>;
  export function GetProjectRoot(): Promise<string>;
  // 可以根据需要添加更多导出函数的声明
}
