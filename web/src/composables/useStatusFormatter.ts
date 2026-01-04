import { useI18n } from 'vue-i18n'
import type { Task } from '@/types/task'

export function useStatusFormatter() {
  const { t } = useI18n()

  const formatStatus = (task: Pick<Task, 'status' | 'attempts' | 'max_attempts'>): string => {
    const { status, attempts = 0, max_attempts = 0 } = task
    switch (status) {
      case 'pending':
        return attempts > 0
          ? t('taskDetail.waitingRetry', { attempt: attempts, max: max_attempts })
          : t('taskDetail.pending')
      case 'running':
        return attempts > 1
          ? t('taskDetail.retrying', { attempt: attempts, max: max_attempts })
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
