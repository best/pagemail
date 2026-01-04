export interface Task {
  id: string
  url: string
  status: 'pending' | 'running' | 'completed' | 'failed'
  formats: string[]
  created_at: string
  updated_at: string
  attempts: number
  max_attempts: number
  error_message?: string
  outputs?: TaskOutput[]
  delivery_history?: DeliveryAttempt[]
}

export interface TaskOutput {
  id: string
  format: string
  size: number
  path: string
}

export interface TaskCreatePayload {
  url: string
  formats: string[]
  cookies?: string
  delivery_config?: { type: 'email' | 'webhook'; id: string }
}

export interface DeliveryAttempt {
  channel: string
  status: string
  attempt_time: string
  error?: string
}
