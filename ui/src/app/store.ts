import { writable, derived } from 'svelte/store'
import { getStoredAuth, setStoredAuth, clearStoredAuth, type Auth } from '../lib/api'

function createAuthStore() {
  const { subscribe, set } = writable<Auth | null>(getStoredAuth())
  return {
    subscribe,
    setAuth: (auth: Auth) => {
      setStoredAuth(auth)
      set(auth)
    },
    clear: () => {
      clearStoredAuth()
      set(null)
    },
  }
}

export interface ToastItem {
  id: number
  type: 'error' | 'success' | 'info'
  message: string
  details?: Record<string, unknown>
}

function createToastStore() {
  const { subscribe, update } = writable<ToastItem[]>([])
  let nextId = 0
  return {
    subscribe,
    push: (type: ToastItem['type'], message: string, details?: Record<string, unknown>) => {
      const id = ++nextId
      update((list) => [...list, { id, type, message, details }])
      setTimeout(() => update((list) => list.filter((t) => t.id !== id)), 8000)
    },
    remove: (id: number) => update((list) => list.filter((t) => t.id !== id)),
  }
}

export const authStore = createAuthStore()
export const isConnected = derived(authStore, ($a) => $a !== null)
export const toastStore = createToastStore()
