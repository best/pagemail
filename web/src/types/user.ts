export interface User {
  id: string
  email: string
  role: 'admin' | 'user'
  is_active: boolean
  avatar_url?: string
  created_at: string
  updated_at: string
  last_login_at?: string
}
