import * as monaco from "monaco-editor";

export class DiagnosticsManager {
  private markers = new Map<string, monaco.editor.IMarkerData[]>();

  /**
   * 更新指定 URI 的诊断信息
   */
  updateDiagnostics(uri: string, diagnostics: any[]) {
    const markers: monaco.editor.IMarkerData[] = diagnostics.map((d) => ({
      severity: this.getSeverity(d.severity),
      startLineNumber: d.range.start.line + 1,
      startColumn: d.range.start.character + 1,
      endLineNumber: d.range.end.line + 1,
      endColumn: d.range.end.character + 1,
      message: d.message,
      source: d.source || "LSP",
    }));

    this.markers.set(uri, markers);
    const model = this.getModelForUri(uri);
    if (model) {
      monaco.editor.setModelMarkers(model, "lsp-owner", markers);
    }
  }

  /**
   * 清除指定 URI 的诊断信息
   */
  clearDiagnostics(uri: string) {
    this.markers.delete(uri);
    const model = this.getModelForUri(uri);
    if (model) {
      monaco.editor.setModelMarkers(model, "lsp-owner", []);
    }
  }

  private getSeverity(lspSeverity: number): monaco.MarkerSeverity {
    switch (lspSeverity) {
      case 1:
        return monaco.MarkerSeverity.Error;
      case 2:
        return monaco.MarkerSeverity.Warning;
      case 3:
        return monaco.MarkerSeverity.Info;
      default:
        return monaco.MarkerSeverity.Hint;
    }
  }

  private getModelForUri(uri: string): monaco.editor.ITextModel | null {
    return monaco.editor.getModel(monaco.Uri.parse(uri));
  }
}
