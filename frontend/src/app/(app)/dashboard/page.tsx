'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatCard from '@/components/StatCard'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import { api } from '@/lib/api'
import { formatCurrency, formatDate, relativeTime } from '@/lib/utils'
import type { Customer, InternetPackage, Subscription, Invoice, Payment, Pagination, Router } from '@/lib/types'

interface CustomersResp { customers: Customer[]; pagination: Pagination }
interface PackagesResp { packages: InternetPackage[] }
interface SubscriptionsResp { subscriptions: Subscription[] }
interface InvoicesResp { invoices: Invoice[] }
interface PaymentsResp { payments: Payment[] }
interface RoutersResp { routers: Router[] }

interface RevenuePoint { date: string; total: number; count: number }
interface RevenueResp { summary: RevenuePoint[]; total: number; currency?: string }

export default function DashboardPage() {
  const [loading, setLoading] = useState(true)
  const [customers, setCustomers] = useState<Customer[]>([])
  const [packages, setPackages] = useState<InternetPackage[]>([])
  const [subs, setSubs] = useState<Subscription[]>([])
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [payments, setPayments] = useState<Payment[]>([])
  const [routers, setRouters] = useState<Router[]>([])
  const [revenue, setRevenue] = useState<RevenueResp | null>(null)

  useEffect(() => {
    let mounted = true
    const all = <T,>(p: Promise<T>) => p.then((v) => v as T | null).catch(() => null)
    Promise.all([
      all(api.get<CustomersResp>('/api/v1/customers', { query: { page_size: 100 } })),
      all(api.get<PackagesResp>('/api/v1/packages', { query: { page_size: 100 } })),
      all(api.get<SubscriptionsResp>('/api/v1/subscriptions', { query: { page_size: 100 } })),
      all(api.get<InvoicesResp>('/api/v1/invoices', { query: { page_size: 100 } })),
      all(api.get<PaymentsResp>('/api/v1/payments', { query: { page_size: 100 } })),
      all(api.get<RoutersResp>('/api/v1/routers', { query: { page_size: 100 } })),
      all(api.get<RevenueResp>('/api/v1/reports/revenue', { query: { days: 30 } })),
    ]).then(([c, p, s, i, pay, r, rev]) => {
      if (!mounted) return
      setCustomers(c?.customers || [])
      setPackages(p?.packages || [])
      setSubs(s?.subscriptions || [])
      setInvoices(i?.invoices || [])
      setPayments(pay?.payments || [])
      setRouters(r?.routers || [])
      setRevenue(rev || null)
      setLoading(false)
    })
    return () => { mounted = false }
  }, [])

  const activeCustomers = customers.filter((c) => c.status === 'ACTIVE').length
  const activeSubs = subs.filter((s) => s.status === 'ACTIVE').length
  const pendingInvoices = invoices.filter((i) => i.status === 'PENDING' || i.status === 'OVERDUE')
  const onlineRouters = routers.filter((r) => r.status === 'ONLINE').length
  const totalRevenue = revenue?.total ?? payments.filter((p) => p.status === 'COMPLETED').reduce((s, p) => s + (p.amount || 0), 0)

  return (
    <div>
      <PageHeader
        title="Overview"
        subtitle="Realtime snapshot of your ISP operations"
        actions={
          <>
            <Link href="/customers/new" className="btn btn-secondary btn-sm">+ Customer</Link>
            <Link href="/subscriptions/new" className="btn btn-primary btn-sm">+ Subscription</Link>
          </>
        }
      />

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4">
        <StatCard
          label="Total Customers"
          value={loading ? '—' : customers.length}
          hint={loading ? '' : `${activeCustomers} active`}
          accent="lime"
          icon={
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          }
        />
        <StatCard
          label="Active Subscriptions"
          value={loading ? '—' : activeSubs}
          hint={loading ? '' : `${subs.length} total`}
          accent="blue"
          icon={
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          }
        />
        <StatCard
          label="Pending Invoices"
          value={loading ? '—' : pendingInvoices.length}
          hint={loading ? '' : `${invoices.length} total`}
          accent="amber"
          icon={
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-6 7h6m-6 4h6" />
            </svg>
          }
        />
        <StatCard
          label="Revenue (30d)"
          value={loading ? '—' : formatCurrency(totalRevenue || 0)}
          hint={loading ? '' : `${payments.filter((p) => p.status === 'COMPLETED').length} completed payments`}
          accent="violet"
          icon={
            <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8v8m0 0v2m0-10V6" />
            </svg>
          }
        />
      </div>

      <div className="mt-6 grid grid-cols-1 lg:grid-cols-3 gap-4">
        <div className="lg:col-span-2 card p-5">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="text-sm font-semibold">Recent Subscriptions</div>
              <div className="text-xs text-[var(--text-dim)]">Latest customer subscriptions</div>
            </div>
            <Link href="/subscriptions" className="text-xs text-[var(--accent)] hover:underline">View all</Link>
          </div>
          {loading ? (
            <SkeletonTable rows={4} />
          ) : subs.length === 0 ? (
            <div className="text-sm text-[var(--text-dim)] py-8 text-center">No subscriptions yet</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="table">
                <thead>
                  <tr>
                    <th>ID</th>
                    <th>Customer</th>
                    <th>Package</th>
                    <th>Status</th>
                    <th>Expires</th>
                  </tr>
                </thead>
                <tbody>
                  {subs.slice(0, 6).map((s) => (
                    <tr key={s.id}>
                      <td className="font-mono text-[var(--text-dim)]">#{s.id}</td>
                      <td className="font-medium">#{s.customer_id}</td>
                      <td className="font-medium">#{s.package_id}</td>
                      <td><StatusBadge status={s.status} /></td>
                      <td className="text-[var(--text-dim)]">{formatDate(s.expiry_date)}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </div>

        <div className="card p-5">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="text-sm font-semibold">Recent Payments</div>
              <div className="text-xs text-[var(--text-dim)]">Last activity</div>
            </div>
            <Link href="/payments" className="text-xs text-[var(--accent)] hover:underline">View all</Link>
          </div>
          {loading ? (
            <SkeletonTable rows={4} />
          ) : payments.length === 0 ? (
            <div className="text-sm text-[var(--text-dim)] py-8 text-center">No payments yet</div>
          ) : (
            <div className="space-y-2">
              {payments.slice(0, 6).map((p) => (
                <div key={p.id} className="flex items-center justify-between rounded-lg border border-[var(--border)] p-3">
                  <div className="min-w-0">
                    <div className="text-sm font-medium truncate">{formatCurrency(p.amount, p.currency)}</div>
                    <div className="text-[11px] text-[var(--text-mute)] truncate">
                      {p.method} · {p.reference || '—'}
                    </div>
                  </div>
                  <div className="text-right shrink-0 ml-3">
                    <StatusBadge status={p.status} />
                    <div className="mt-1 text-[11px] text-[var(--text-mute)]">{relativeTime(p.paid_at || p.created_at)}</div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>

      <div className="mt-4 grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div className="card p-5">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="text-sm font-semibold">Active Packages</div>
              <div className="text-xs text-[var(--text-dim)]">Plans you offer to customers</div>
            </div>
            <Link href="/packages" className="text-xs text-[var(--accent)] hover:underline">Manage</Link>
          </div>
          {loading ? (
            <SkeletonTable rows={3} />
          ) : packages.length === 0 ? (
            <div className="text-sm text-[var(--text-dim)] py-6 text-center">No packages yet</div>
          ) : (
            <div className="space-y-2">
              {packages.filter((p) => p.is_active).slice(0, 5).map((p) => (
                <div key={p.id} className="flex items-center justify-between rounded-lg border border-[var(--border)] p-3">
                  <div className="min-w-0">
                    <div className="text-sm font-medium truncate">{p.name}</div>
                    <div className="text-[11px] text-[var(--text-mute)]">
                      {p.download_mbps}/{p.upload_mbps} Mbps · {p.duration_days} days
                    </div>
                  </div>
                  <div className="text-right shrink-0 ml-3">
                    <div className="text-sm font-semibold text-[var(--accent)]">
                      {formatCurrency(p.price, p.currency)}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        <div className="card p-5">
          <div className="flex items-center justify-between mb-4">
            <div>
              <div className="text-sm font-semibold">Network Status</div>
              <div className="text-xs text-[var(--text-dim)]">MikroTik routers</div>
            </div>
            <Link href="/routers" className="text-xs text-[var(--accent)] hover:underline">Manage</Link>
          </div>
          {loading ? (
            <SkeletonTable rows={3} />
          ) : routers.length === 0 ? (
            <div className="text-sm text-[var(--text-dim)] py-6 text-center">No routers registered</div>
          ) : (
            <div className="space-y-2">
              <div className="flex items-center justify-between rounded-lg border border-[var(--border)] p-3 bg-[var(--bg-elev-2)]">
                <div className="text-sm text-[var(--text-dim)]">Online routers</div>
                <div className="text-sm font-semibold">
                  {onlineRouters} / {routers.length}
                </div>
              </div>
              {routers.slice(0, 4).map((r) => (
                <div key={r.id} className="flex items-center justify-between rounded-lg border border-[var(--border)] p-3">
                  <div className="min-w-0">
                    <div className="text-sm font-medium truncate">{r.name}</div>
                    <div className="text-[11px] text-[var(--text-mute)] truncate">{r.host}:{r.api_port}</div>
                  </div>
                  <StatusBadge status={r.status} />
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
