'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import EmptyState from '@/components/EmptyState'
import { api } from '@/lib/api'
import { formatCurrency, formatDate, relativeTime } from '@/lib/utils'
import type { Customer, Invoice, Pagination } from '@/lib/types'

interface InvoicesResp { invoices: Invoice[]; pagination: Pagination }
interface CustomersResp { customers: Customer[]; pagination: Pagination }

export default function BillingPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [customers, setCustomers] = useState<Customer[]>([])
  const [pagination, setPagination] = useState<Pagination | null>(null)
  const [loading, setLoading] = useState(true)
  const [status, setStatus] = useState('')

  useEffect(() => {
    let mounted = true
    setLoading(true)
    Promise.all([
      api.get<InvoicesResp>('/api/v1/invoices', { query: { page_size: 50, status: status || undefined } }),
      api.get<CustomersResp>('/api/v1/customers', { query: { page_size: 200 } }),
    ]).then(([i, c]) => {
      if (!mounted) return
      setInvoices(i.invoices)
      setPagination(i.pagination)
      setCustomers(c.customers)
    }).finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [status])

  const customerById = (id: number) => customers.find((c) => c.id === id)

  const totalAmount = invoices.reduce((s, i) => s + (i.amount || 0), 0)
  const currency = invoices[0]?.currency || 'KES'
  const outstanding = invoices.filter((i) => i.status === 'PENDING' || i.status === 'OVERDUE')
  const outstandingAmount = outstanding.reduce((s, i) => s + (i.amount || 0), 0)

  const filters = [
    { label: 'All', value: '' },
    { label: 'Pending', value: 'PENDING' },
    { label: 'Paid', value: 'PAID' },
    { label: 'Overdue', value: 'OVERDUE' },
    { label: 'Cancelled', value: 'CANCELLED' },
  ]

  return (
    <div>
      <PageHeader
        title="Invoices"
        subtitle="Billing documents generated for subscriptions"
        actions={<Link href="/subscriptions/new" className="btn btn-secondary btn-sm">+ From subscription</Link>}
      />

      <div className="grid grid-cols-1 sm:grid-cols-3 gap-4 mb-4">
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Total Invoices</div>
          <div className="mt-1 text-2xl font-semibold">{pagination?.total ?? invoices.length}</div>
        </div>
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Outstanding</div>
          <div className="mt-1 text-2xl font-semibold">{outstanding.length}</div>
        </div>
        <div className="card p-4">
          <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Outstanding value</div>
          <div className="mt-1 text-2xl font-semibold text-[var(--accent)]">{formatCurrency(outstandingAmount, currency)}</div>
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
      ) : invoices.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 7h6m-6 4h6" />
            </svg>
          }
          title="No invoices yet"
          description="Invoices are generated automatically when subscriptions are created."
        />
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Invoice</th>
                  <th>Customer</th>
                  <th>Subscription</th>
                  <th>Amount</th>
                  <th>Status</th>
                  <th>Due</th>
                </tr>
              </thead>
              <tbody>
                {invoices.map((inv) => {
                  const c = customerById(inv.customer_id)
                  return (
                    <tr key={inv.id}>
                      <td className="font-mono text-sm">{inv.invoice_no}</td>
                      <td>
                        {c ? (
                          <div>
                            <div className="font-medium">{c.full_name}</div>
                            <div className="text-xs text-[var(--text-mute)]">@{c.username}</div>
                          </div>
                        ) : (
                          <span className="text-[var(--text-dim)]">#{inv.customer_id}</span>
                        )}
                      </td>
                      <td className="font-mono text-[var(--text-dim)]">#{inv.subscription_id}</td>
                      <td className="font-semibold">{formatCurrency(inv.amount, inv.currency)}</td>
                      <td><StatusBadge status={inv.status} /></td>
                      <td>
                        <div className="text-sm">{formatDate(inv.due_date)}</div>
                        <div className="text-xs text-[var(--text-mute)]">{relativeTime(inv.due_date)}</div>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
          <div className="px-4 py-3 border-t border-[var(--border)] text-xs text-[var(--text-mute)] flex justify-between">
            <span>{invoices.length} shown · total {formatCurrency(totalAmount, currency)}</span>
            {pagination && <span>of {pagination.total}</span>}
          </div>
        </div>
      )}
    </div>
  )
}
