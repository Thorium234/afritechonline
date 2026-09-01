export type Role = 'SUPER_ADMIN' | 'ADMIN' | 'STAFF' | 'CUSTOMER'

export interface User {
  id: number
  username: string
  email: string
  role: Role
  is_active: boolean
  last_login_at?: string
  created_at: string
  updated_at: string
}

export interface AuthTokens {
  access_token: string
  refresh_token: string
}

export interface Customer {
  id: number
  user_id?: number
  full_name: string
  phone: string
  email: string
  username: string
  status: 'ACTIVE' | 'INACTIVE' | 'SUSPENDED'
  created_at: string
  updated_at: string
}

export interface InternetPackage {
  id: number
  name: string
  description: string
  price: number
  currency: string
  duration_days: number
  download_mbps: number
  upload_mbps: number
  data_limit_gb?: number | null
  is_active: boolean
  created_at: string
  updated_at: string
}

export type SubscriptionStatus =
  | 'PENDING'
  | 'ACTIVE'
  | 'EXPIRED'
  | 'SUSPENDED'
  | 'CANCELLED'

export interface Subscription {
  id: number
  customer_id: number
  package_id: number
  start_date: string
  expiry_date: string
  status: SubscriptionStatus
  amount: number
  currency: string
  created_at: string
  updated_at: string
}

export type InvoiceStatus = 'PENDING' | 'PAID' | 'OVERDUE' | 'CANCELLED'

export interface Invoice {
  id: number
  invoice_no: string
  subscription_id: number
  customer_id: number
  amount: number
  currency: string
  status: InvoiceStatus
  due_date: string
  created_at: string
  updated_at: string
}

export type PaymentStatus = 'PENDING' | 'COMPLETED' | 'FAILED' | 'CANCELLED'
export type PaymentMethod = 'MANUAL' | 'MPESA' | 'CARD' | 'OTHER'

export interface Payment {
  id: number
  invoice_id: number
  customer_id: number
  amount: number
  currency: string
  method: PaymentMethod
  reference: string
  status: PaymentStatus
  paid_at?: string
  created_at: string
  updated_at: string
}

export interface Router {
  id: number
  name: string
  host: string
  api_port: number
  username: string
  location: string
  status: 'OFFLINE' | 'ONLINE' | 'UNKNOWN'
  created_at: string
  updated_at: string
}

export interface Pagination {
  page: number
  page_size: number
  total: number
}

export interface ApiSuccess<T> {
  data: T
}

export interface ApiErrorBody {
  error: {
    status: number
    message: string
    fields?: Record<string, string>
  }
}

export class ApiError extends Error {
  status: number
  fields?: Record<string, string>
  constructor(message: string, status: number, fields?: Record<string, string>) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.fields = fields
  }
}
