'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import EmptyState from '@/components/EmptyState'
import { api } from '@/lib/api'
import { formatCurrency, formatDate, relativeTime } from '@/lib/utils'
import type { Customer, InternetPackage, Subscription, Pagination } from '@/lib/types'

interface SubsResp { subscriptions: Subscription[]; pagination: Pagination }
interface CustomersResp { customers: Customer[]; pagination: Pagination }
interface PackagesResp { packages: InternetPackage[] }

export default function SubscriptionsPage() {
  const [subs, setSubs] = useState<Subscription[]>([])
  const [pagination, setPagination] = useState<Pagination | null>(null)
  const [customers, setCustomers] = useState<Customer[]>([])
  const [packages, setPackages] = useState<InternetPackage[]>([])
  const [loading, setLoading] = useState(true)
  const [status, setStatus] = useState<string>('')

  useEffect(() => {
    let mounted = true
    setLoading(true)
    Promise.all([
      api.get<SubsResp>('/api/v1/subscriptions', { query: { page_size: 50, status: status || undefined } }),
      api.get<CustomersResp>('/api/v1/customers', { query: { page_size: 200 } }),
      api.get<PackagesResp>('/api/v1/packages', { query: { page_size: 100 } }),
    ]).then(([s, c, p]) => {
      if (!mounted) return
      setSubs(s.subscriptions)
      setPagination(s.pagination)
      setCustomers(c.customers)
      setPackages(p.packages)
    }).finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [status])

  const customerById = (id: number) => customers.find((c) => c.id === id)
  const packageById = (id: number) => packages.find((p) => p.id === id)

  const filters: { label: string; value: string }[] = [
    { label: 'All', value: '' },
    { label: 'Active', value: 'ACTIVE' },
    { label: 'Pending', value: 'PENDING' },
    { label: 'Expired', value: 'EXPIRED' },
    { label: 'Suspended', value: 'SUSPENDED' },
    { label: 'Cancelled', value: 'CANCELLED' },
  ]

  return (
    <div>
      <PageHeader
        title="Subscriptions"
        subtitle="Active and historical customer subscriptions"
        actions={
          <Link href="/subscriptions/new" className="btn btn-primary btn-sm">+ New subscription</Link>
        }
      />

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
      ) : subs.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          }
          title="No subscriptions"
          description="Assign a package to a customer to create a subscription."
          action={<Link href="/subscriptions/new" className="btn btn-primary btn-sm">+ New subscription</Link>}
        />
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Customer</th>
                  <th>Package</th>
                  <th>Amount</th>
                  <th>Status</th>
                  <th>Expires</th>
                  <th className="w-24">Actions</th>
                </tr>
              </thead>
              <tbody>
                {subs.map((s) => {
                  const c = customerById(s.customer_id)
                  const p = packageById(s.package_id)
                  return (
                    <tr key={s.id}>
                      <td className="font-mono text-[var(--text-dim)]">#{s.id}</td>
                      <td>
                        {c ? (
                          <div>
                            <div className="font-medium">{c.full_name}</div>
                            <div className="text-xs text-[var(--text-mute)]">@{c.username}</div>
                          </div>
                        ) : (
                          <span className="text-[var(--text-dim)]">#{s.customer_id}</span>
                        )}
                      </td>
                      <td>
                        {p ? (
                          <div>
                            <div className="font-medium">{p.name}</div>
                            <div className="text-xs text-[var(--text-mute)]">{p.download_mbps}/{p.upload_mbps} Mbps · {p.duration_days}d</div>
                          </div>
                        ) : (
                          <span className="text-[var(--text-dim)]">#{s.package_id}</span>
                        )}
                      </td>
                      <td className="font-semibold">{formatCurrency(s.amount, s.currency)}</td>
                      <td><StatusBadge status={s.status} /></td>
                      <td>
                        <div className="text-sm">{formatDate(s.expiry_date)}</div>
                        <div className="text-xs text-[var(--text-mute)]">{relativeTime(s.expiry_date)}</div>
                      </td>
                      <td className="whitespace-nowrap">
                        <Link href={`/subscriptions/${s.id}`} className="btn btn-secondary btn-sm">Edit</Link>
                      </td>
                    </tr>
                  )
                })}
              </tbody>
            </table>
          </div>
          {pagination && (
            <div className="px-4 py-3 border-t border-[var(--border)] text-xs text-[var(--text-mute)]">
              Showing {subs.length} of {pagination.total}
            </div>
          )}
        </div>
      )}
    </div>
  )
}
