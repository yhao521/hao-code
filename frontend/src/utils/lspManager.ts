import {
  GetCompletions,
  InitializeLSP,
  GetDefinition,
  GetDocumentSymbols,
  FindReferences,
  RenameSymbol,
  FormatDocument,
  GetHoverInfo,
  GetSignatureHelp,
  GetCodeActions,
  GetFoldingRanges,
} from "@wails/backend/appservice.js";

export class LSPManager {
  private static instance: LSPManager;
  private initializedLanguages = new Set<string>();

  static getInstance(): LSPManager {
    if (!LSPManager.instance) {
      LSPManager.instance = new LSPManager();
    }
    return LSPManager.instance;
  }

  async ensureInitialized(languageId: string, rootPath: string) {
    if (this.initializedLanguages.has(languageId)) return;

    try {
      await InitializeLSP(languageId, rootPath);
      this.initializedLanguages.add(languageId);
      console.log(`LSP for ${languageId} initialized.`);
    } catch (error) {
      console.error(`Failed to initialize LSP for ${languageId}:`, error);
    }
  }

  async getCompletions(
    languageId: string,
    uri: string,
    line: number,
    col: number,
  ) {
    try {
      const result = await GetCompletions(languageId, uri, line, col);
      return result || [];
    } catch (error) {
      console.error("LSP completion error:", error);
      return [];
    }
  }

  async getDefinition(
    languageId: string,
    uri: string,
    line: number,
    col: number,
  ) {
    try {
      return await GetDefinition(languageId, uri, line, col);
    } catch (error) {
      console.error("LSP definition error:", error);
      return null;
    }
  }

  async getDocumentSymbols(languageId: string, uri: string) {
    try {
      return await GetDocumentSymbols(languageId, uri);
    } catch (error) {
      console.error("LSP document symbols error:", error);
      return [];
    }
  }

  async findReferences(
    languageId: string,
    uri: string,
    line: number,
    col: number,
  ) {
    try {
      return await FindReferences(languageId, uri, line, col);
    } catch (error) {
      console.error("LSP find references error:", error);
      return [];
    }
  }

  async renameSymbol(
    languageId: string,
    uri: string,
    line: number,
    col: number,
    newName: string,
  ) {
    try {
      return await RenameSymbol(languageId, uri, line, col, newName);
    } catch (error) {
      console.error("LSP rename symbol error:", error);
      return null;
    }
  }

  async formatDocument(languageId: string, uri: string, content: string) {
    try {
      return await FormatDocument(languageId, uri, content);
    } catch (error) {
      console.error("LSP format document error:", error);
      return [];
    }
  }

  async getHoverInfo(
    languageId: string,
    uri: string,
    line: number,
    col: number,
  ) {
    try {
      return await GetHoverInfo(languageId, uri, line, col);
    } catch (error) {
      console.error("LSP hover info error:", error);
      return null;
    }
  }

  async getSignatureHelp(
    languageId: string,
    uri: string,
    line: number,
    col: number,
  ) {
    try {
      return await GetSignatureHelp(languageId, uri, line, col);
    } catch (error) {
      console.error("LSP signature help error:", error);
      return null;
    }
  }

  async getCodeActions(
    languageId: string,
    uri: string,
    startLine: number,
    startCol: number,
    endLine: number,
    endCol: number,
    diagnostics: any[]
  ) {
    try {
      return await GetCodeActions(languageId, uri, startLine, startCol, endLine, endCol, diagnostics);
    } catch (error) {
      console.error("LSP code actions error:", error);
      return [];
    }
  }

  async getFoldingRanges(languageId: string, uri: string) {
    try {
      return await GetFoldingRanges(languageId, uri);
    } catch (error) {
      console.error("LSP folding ranges error:", error);
      return [];
    }
  }
}
