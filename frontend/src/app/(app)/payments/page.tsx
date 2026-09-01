'use client'

import { useEffect, useState } from 'react'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import EmptyState from '@/components/EmptyState'
import { api } from '@/lib/api'
import { formatCurrency, formatDateTime, relativeTime } from '@/lib/utils'
import type { Customer, Payment, Pagination } from '@/lib/types'

interface PaymentsResp { payments: Payment[]; pagination: Pagination }
interface CustomersResp { customers: Customer[]; pagination: Pagination }

export default function PaymentsPage() {
  const [payments, setPayments] = useState<Payment[]>([])
  const [customers, setCustomers] = useState<Customer[]>([])
  const [pagination, setPagination] = useState<Pagination | null>(null)
  const [loading, setLoading] = useState(true)
  const [status, setStatus] = useState('')

  useEffect(() => {
    let mounted = true
    setLoading(true)
    Promise.all([
      api.get<PaymentsResp>('/api/v1/payments', { query: { page_size: 50, status: status || undefined } }),
      api.get<CustomersResp>('/api/v1/customers', { query: { page_size: 200 } }),
    ]).then(([p, c]) => {
      if (!mounted) return
      setPayments(p.payments)
      setPagination(p.pagination)
      setCustomers(c.customers)
    }).finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [status])

  const customerById = (id: number) => customers.find((c) => c.id === id)

  const totalCompleted = payments.filter((p) => p.status === 'COMPLETED').reduce((s, p) => s + (p.amount || 0), 0)
  const currency = payments[0]?.currency || 'KES'

  const filters = [
    { label: 'All', value: '' },
    { label: 'Completed', value: 'COMPLETED' },
    { label: 'Pending', value: 'PENDING' },
    { label: 'Failed', value: 'FAILED' },
    { label: 'Cancelled', value: 'CANCELLED' },
  ]

  return (
    <div>
      <PageHeader
        title="Payments"
        subtitle="M-Pesa, manual and card payments"
      />

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Total</div>
          <div className="mt-1 text-2xl font-semibold">{pagination?.total ?? payments.length}</div>
        </div>
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Completed</div>
          <div className="mt-1 text-2xl font-semibold text-[#6fdc8c]">
            {payments.filter((p) => p.status === 'COMPLETED').length}
          </div>
        </div>
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Collected</div>
          <div className="mt-1 text-2xl font-semibold text-[var(--accent)]">{formatCurrency(totalCompleted, currency)}</div>
        </div>
      </div>

      <div className="mb-4 flex flex-wrap items-center gap-1.5">
        {filters.map((f) => (
          <button
            key={f.label}
            onClick={() => setStatus(f.value)}
            className={`px-3 py-1.5 rounded-md text-xs font-medium border transition ${
              status === f.value
                ? 'bg-[var(--accent-soft)] text-[var(--accent)] border-[var(--accent)]/30'
                : 'bg-[var(--bg-elev)] text-[var(--text-dim)] border-[var(--border)] hover:text-[var(--text)]'
            }`}
          >{f.label}</button>
        ))}
      </div>

      {loading ? (
        <SkeletonTable rows={6} />
      ) : payments.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3 10h18M5 6h14a2 2 0 012 2v8a2 2 0 01-2 2H5a2 2 0 01-2-2V8a2 2 0 012-2z" />
            </svg>
          }
          title="No payments yet"
          description="Payments will appear here as customers pay for subscriptions."
        />
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Customer</th>
                  <th>Method</th>
                  <th>Reference</th>
                  <th>Amount</th>
                  <th>Status</th>
                  <th>When</th>
                </tr>
              </thead>
              <tbody>
                {payments.map((p) => {
                  const c = customerById(p.customer_id)
                  return (
                    <tr key={p.id}>
                      <td className="font-mono text-[var(--text-dim)]">#{p.id}</td>
                      <td>
                        {c ? (
                          <div>
                            <div className="font-medium">{c.full_name}</div>
                            <div className="text-xs text-[var(--text-mute)]">@{c.username}</div>
                          </div>
                        ) : (
                          <span className="text-[var(--text-dim)]">#{p.customer_id}</span>
                        )}
                      </td>
                      <td><span className="badge badge-neutral">{p.method}</span></td>
                      <td className="font-mono text-xs text-[var(--text-dim)]">{p.reference || '—'}</td>
                      <td className="font-semibold">{formatCurrency(p.amount, p.currency)}</td>
                      <td><StatusBadge status={p.status} /></td>
                      <td>
                        <div className="text-sm">{formatDateTime(p.paid_at || p.created_at)}</div>
                        <div className="text-xs text-[var(--text-mute)]">{relativeTime(p.paid_at || p.created_at)}</div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
