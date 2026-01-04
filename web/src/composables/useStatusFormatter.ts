import { useI18n } from 'vue-i18n'
import type { Task } from '@/types/task'

export function useStatusFormatter() {
  const { t } = useI18n()

  const formatStatus = (
    task: Pick<Task, 'status' | 'attempts' | 'max_attempts'>,
    short = false
  ): string => {
    const { status, attempts = 0, max_attempts = 0 } = task
    const params = { attempt: attempts, max: max_attempts }
    switch (status) {
      case 'pending':
        return attempts > 0
          ? t(short ? 'taskDetail.waitingRetryShort' : 'taskDetail.waitingRetry', params)
          : t('taskDetail.pending')
      case 'running':
        return attempts > 1
          ? t(short ? 'taskDetail.retryingShort' : 'taskDetail.retrying', params)
          : t('taskDetail.processing')
      case 'failed':
        return t('taskDetail.failed')
      case 'completed':
        return t('taskDetail.completed')
      default:
        return status
    }
  }

  return { formatStatus }
}
