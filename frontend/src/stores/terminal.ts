import { defineStore } from "pinia";
import { ref, computed } from "vue";

export interface TerminalInstance {
  id: string;
  name: string;
}

export const useTerminalStore = defineStore("terminal", () => {
  const instances = ref<TerminalInstance[]>([]);
  const activeId = ref<string | null>(null);
  const splitDirection = ref<"horizontal" | "vertical">("horizontal");

  const activeInstance = computed(() =>
    instances.value.find((t) => t.id === activeId.value),
  );

  function addInstance(id: string) {
    const name = `Terminal ${instances.value.length + 1}`;
    instances.value.push({ id, name });
    if (!activeId.value) {
      activeId.value = id;
    }
  }

  function removeInstance(id: string) {
    const index = instances.value.findIndex((t) => t.id === id);
    if (index !== -1) {
      instances.value.splice(index, 1);
      if (activeId.value === id) {
        activeId.value =
          instances.value.length > 0 ? instances.value[0].id : null;
      }
    }
  }

  function setActive(id: string) {
    activeId.value = id;
  }

  function toggleSplitDirection() {
    splitDirection.value =
      splitDirection.value === "horizontal" ? "vertical" : "horizontal";
  }

  return {
    instances,
    activeId,
    activeInstance,
    splitDirection,
    addInstance,
    removeInstance,
    setActive,
    toggleSplitDirection,
  };
});
