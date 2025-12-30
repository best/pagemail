import { ref, onMounted, onUnmounted, unref, watch, type Ref } from 'vue'

type MaybeRef<T> = T | Ref<T>

const DEFAULT_INTERVAL_MS = 10000
const PENDING_INTERVAL_MS = 5000

export interface UsePollingOptions {
  intervalMs?: MaybeRef<number>
  pendingIntervalMs?: MaybeRef<number>
  isPending?: Ref<boolean> | (() => boolean)
  immediate?: boolean
  autoStart?: boolean
}

export interface UsePollingControls {
  start: () => void
  stop: () => void
  pause: () => void
  resume: () => void
  isActive: Ref<boolean>
  isPaused: Ref<boolean>
  isRunning: Ref<boolean>
}

export function usePolling(
  callback: () => void | Promise<void>,
  options: UsePollingOptions = {}
): UsePollingControls {
  const {
    intervalMs = DEFAULT_INTERVAL_MS,
    pendingIntervalMs = PENDING_INTERVAL_MS,
    isPending,
    immediate = true,
    autoStart = true
  } = options

  const isActive = ref(false)
  const isPaused = ref(false)
  const isRunning = ref(false)

  let timerId: ReturnType<typeof setTimeout> | null = null

  const clearTimer = () => {
    if (timerId !== null) {
      clearTimeout(timerId)
      timerId = null
    }
  }

  const isDocumentHidden = () =>
    typeof document !== 'undefined' && document.visibilityState === 'hidden'

  const getIsPending = () => {
    if (typeof isPending === 'function') return isPending()
    return Boolean(isPending?.value)
  }

  const resolveInterval = () => {
    const base = getIsPending() ? pendingIntervalMs : intervalMs
    const value = unref(base)
    return Number.isFinite(value) && value > 0 ? value : DEFAULT_INTERVAL_MS
  }

  const scheduleNext = (delay?: number) => {
    clearTimer()
    if (!isActive.value || isPaused.value) return
    const ms = typeof delay === 'number' ? delay : resolveInterval()
    timerId = setTimeout(() => void tick(), ms)
  }

  const tick = async () => {
    clearTimer()
    if (!isActive.value || isPaused.value) return
    if (isDocumentHidden()) {
      scheduleNext()
      return
    }
    if (isRunning.value) {
      scheduleNext()
      return
    }
    isRunning.value = true
    try {
      await callback()
    } catch {
      // Callback errors are swallowed to keep polling alive
      // The callback should handle its own error state
    } finally {
      isRunning.value = false
      scheduleNext()
    }
  }

  const start = () => {
    if (isActive.value) return
    isActive.value = true
    isPaused.value = false
    if (immediate) {
      void tick()
    } else {
      scheduleNext()
    }
  }

  const stop = () => {
    isActive.value = false
    isPaused.value = false
    clearTimer()
  }

  const pause = () => {
    if (!isActive.value || isPaused.value) return
    isPaused.value = true
    clearTimer()
  }

  const resume = () => {
    if (!isActive.value || !isPaused.value) return
    isPaused.value = false
    scheduleNext(immediate ? 0 : undefined)
  }

  const handleVisibilityChange = () => {
    if (!isActive.value || isPaused.value) return
    if (!isDocumentHidden()) {
      scheduleNext(0)
    }
  }

  if (typeof document !== 'undefined') {
    document.addEventListener('visibilitychange', handleVisibilityChange)
  }

  watch(
    () => [unref(intervalMs), unref(pendingIntervalMs), getIsPending()],
    () => {
      if (isActive.value && !isPaused.value) {
        scheduleNext()
      }
    }
  )

  if (autoStart) {
    onMounted(start)
  }

  onUnmounted(() => {
    stop()
    if (typeof document !== 'undefined') {
      document.removeEventListener('visibilitychange', handleVisibilityChange)
    }
  })

  return { start, stop, pause, resume, isActive, isPaused, isRunning }
}
