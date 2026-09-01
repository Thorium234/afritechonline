'use client'

import { useEffect, useState } from 'react'
import PageHeader from '@/components/PageHeader'
import StatCard from '@/components/StatCard'
import { SkeletonTable } from '@/components/Skeleton'
import { api } from '@/lib/api'
import { formatCurrency, formatDate } from '@/lib/utils'

interface RevenuePoint { date: string; total: number; count: number }
interface RevenueResp { summary: RevenuePoint[]; total: number; currency?: string }
interface CustomerStatsResp { total: number; active: number; inactive: number; suspended: number }
interface RoutersResp { routers: { id: number; status: string }[] }

export default function ReportsPage() {
  const [days, setDays] = useState(30)
  const [revenue, setRevenue] = useState<RevenueResp | null>(null)
  const [customers, setCustomers] = useState<CustomerStatsResp | null>(null)
  const [routers, setRouters] = useState<RoutersResp | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    Promise.all([
      api.get<RevenueResp>('/api/v1/reports/revenue', { query: { days } }),
      api.get<CustomerStatsResp>('/api/v1/reports/customers'),
      api.get<RoutersResp>('/api/v1/routers', { query: { page_size: 100 } }),
    ]).then(([r, c, rt]) => {
      if (!mounted) return
      setRevenue(r)
      setCustomers(c)
      setRouters(rt)
    }).finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [days])

  const max = revenue?.summary?.reduce((m, p) => Math.max(m, p.total), 0) || 1
  const onlineRouters = routers?.routers.filter((r) => r.status === 'ONLINE').length || 0

  return (
    <div>
      <PageHeader
        title="Reports"
        subtitle="Operational and financial insights"
        actions={
          <div className="flex items-center gap-1 rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] p-0.5 text-xs">
            {[7, 30, 90].map((d) => (
              <button
                key={d}
                onClick={() => setDays(d)}
                className={`px-3 py-1.5 rounded-md transition ${days === d ? 'bg-[var(--accent-soft)] text-[var(--accent)]' : 'text-[var(--text-dim)] hover:text-[var(--text)]'}`}
              >{d}d</button>
            ))}
          </div>
        }
      />

      <div className="grid grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
        <StatCard label="Revenue" value={loading ? '—' : formatCurrency(revenue?.total || 0)} hint={`Last ${days} days`} accent="lime" />
        <StatCard label="Customers" value={loading ? '—' : customers?.total ?? '—'} hint={customers ? `${customers.active} active` : ''} accent="blue" />
        <StatCard label="Active routers" value={loading ? '—' : `${onlineRouters} / ${routers?.routers.length || 0}`} accent="amber" />
        <StatCard label="Transactions" value={loading ? '—' : revenue?.summary?.reduce((s, p) => s + p.count, 0) ?? 0} hint={`Last ${days} days`} accent="violet" />
      </div>

      <div className="card p-5">
        <div className="flex items-center justify-between mb-4">
          <div>
            <div className="text-sm font-semibold">Revenue trend</div>
            <div className="text-xs text-[var(--text-dim)]">Daily revenue (last {days} days)</div>
          </div>
        </div>
        {loading ? (
          <SkeletonTable rows={5} />
        ) : !revenue?.summary?.length ? (
          <div className="text-sm text-[var(--text-dim)] py-8 text-center">No revenue data for this period</div>
        ) : (
          <div className="space-y-2">
            {revenue.summary.slice(-14).map((p) => (
              <div key={p.date} className="grid grid-cols-12 items-center gap-3 text-sm">
                <div className="col-span-3 sm:col-span-2 text-[var(--text-dim)]">{formatDate(p.date)}</div>
                <div className="col-span-7 sm:col-span-8">
                  <div className="h-2 rounded-full bg-[var(--bg-elev-2)] overflow-hidden">
                    <div
                      className="h-full bg-gradient-to-r from-[var(--accent)] to-[#79b8ff]"
                      style={{ width: `${Math.max(2, (p.total / max) * 100)}%` }}
                    />
                  </div>
                </div>
                <div className="col-span-2 text-right font-semibold">{formatCurrency(p.total)}</div>
              </div>
            ))}
          </div>
        )}
      </div>
    </div>
  )
}
