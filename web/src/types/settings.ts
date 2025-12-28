export interface SmtpProfile {
  id: string
  name: string
  host: string
  port: number
  username: string
  password?: string
  from_email: string
  from_name: string
  use_tls: boolean
  is_default: boolean
  created_at: string
  updated_at: string
}

export interface WebhookConfig {
  id: string
  name: string
  url: string
  secret?: string
  headers?: Record<string, string>
  is_active: boolean
  created_at: string
  updated_at: string
}
