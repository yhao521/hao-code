import { defineStore } from "pinia";
import { ref } from "vue";
import * as monaco from "monaco-editor";

export const useDiagnosticsStore = defineStore("diagnostics", () => {
  // State: Map of file path to markers
  const markers = ref<Record<string, monaco.editor.IMarkerData[]>>({});

  // Actions
  function setMarkers(path: string, newMarkers: monaco.editor.IMarkerData[]) {
    markers.value[path] = newMarkers;
  }

  function clearMarkers(path: string) {
    delete markers.value[path];
  }

  return {
    markers,
    setMarkers,
    clearMarkers,
  };
});
