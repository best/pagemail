export interface ApiResponse<T> {
  data: T
  message?: string
}

export interface ProblemDetail {
  type: string
  title: string
  status: number
  detail?: string
  instance?: string
  errors?: Array<{ field: string; message: string }>
}

export interface PaginatedResponse<T> {
  data: T[]
  meta: {
    page: number
    per_page: number
    total: number
    total_pages: number
  }
}

export interface SiteConfig {
  site_name: string
  site_slogan: string
}
