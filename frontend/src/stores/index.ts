import { reactive, ref } from "vue";
import { defineStore } from "pinia";

type Conn = WebSocket | null

export const useStore = defineStore("data", () => {
  const isAuthenticated = ref(false);
  const conn = ref<Conn>(null)
  return { isAuthenticated, conn };
});
