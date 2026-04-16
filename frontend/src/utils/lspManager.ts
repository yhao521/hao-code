import {
  GetCompletions,
  InitializeLSP,
  GetDefinition,
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
}
