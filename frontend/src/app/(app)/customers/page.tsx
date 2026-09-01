'use client'

import { useEffect, useMemo, useState } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import EmptyState from '@/components/EmptyState'
import { api } from '@/lib/api'
import { formatDate, relativeTime, initials } from '@/lib/utils'
import type { Customer, Pagination } from '@/lib/types'

interface CustomersResp { customers: Customer[]; pagination: Pagination }

export default function CustomersPage() {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [pagination, setPagination] = useState<Pagination | null>(null)
  const [loading, setLoading] = useState(true)
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    api
      .get<CustomersResp>('/api/v1/customers', { query: { page, page_size: 20, search } })
      .then((res) => {
        if (!mounted) return
        setCustomers(res.customers)
        setPagination(res.pagination)
      })
      .finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [page, search])

  const totalPages = pagination ? Math.max(1, Math.ceil(pagination.total / pagination.page_size)) : 1

  return (
    <div>
      <PageHeader
        title="Customers"
        subtitle="Subscribers on your network"
        actions={
          <>
            <div className="relative">
              <input
                value={search}
                onChange={(e) => { setSearch(e.target.value); setPage(1) }}
                placeholder="Search name, phone, email…"
                className="input pl-9 w-64"
              />
              <svg className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M21 21l-4.35-4.35M11 19a8 8 0 100-16 8 8 0 000 16z" />
              </svg>
            </div>
            <Link href="/customers/new" className="btn btn-primary btn-sm">+ New customer</Link>
          </>
        }
      />

      {loading ? (
        <SkeletonTable rows={6} />
      ) : customers.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          }
          title={search ? 'No customers match your search' : 'No customers yet'}
          description={search ? 'Try a different name, phone, email or username.' : 'Get started by creating your first subscriber.'}
          action={!search && <Link href="/customers/new" className="btn btn-primary btn-sm">+ Add customer</Link>}
        />
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Customer</th>
                  <th>Phone</th>
                  <th>Email</th>
                  <th>Status</th>
                  <th>Joined</th>
                </tr>
              </thead>
              <tbody>
                {customers.map((c) => (
                  <tr key={c.id}>
                    <td>
                      <div className="flex items-center gap-3">
                        <div className="grid h-9 w-9 place-items-center rounded-full bg-[var(--bg-elev-2)] border border-[var(--border)] text-xs font-semibold text-[var(--text-dim)]">
                          {initials(c.full_name || c.username)}
                        </div>
                        <div className="min-w-0">
                          <div className="font-medium truncate">{c.full_name}</div>
                          <div className="text-xs text-[var(--text-mute)] truncate">@{c.username}</div>
                        </div>
                      </div>
                    </td>
                    <td className="text-[var(--text-dim)]">{c.phone || '—'}</td>
                    <td className="text-[var(--text-dim)]">{c.email || '—'}</td>
                    <td><StatusBadge status={c.status} /></td>
                    <td>
                      <div className="text-sm">{formatDate(c.created_at)}</div>
                      <div className="text-xs text-[var(--text-mute)]">{relativeTime(c.created_at)}</div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
          {pagination && pagination.total > pagination.page_size && (
            <div className="flex items-center justify-between border-t border-[var(--border)] px-4 py-3 text-sm text-[var(--text-dim)]">
              <div>Page {pagination.page} of {totalPages} · {pagination.total} total</div>
              <div className="flex gap-1">
                <button className="btn btn-secondary btn-sm" disabled={page <= 1} onClick={() => setPage((p) => Math.max(1, p - 1))}>Prev</button>
                <button className="btn btn-secondary btn-sm" disabled={page >= totalPages} onClick={() => setPage((p) => Math.min(totalPages, p + 1))}>Next</button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
