'use client'

import { useEffect, useState } from 'react'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import EmptyState from '@/components/EmptyState'
import { SkeletonCards } from '@/components/Skeleton'
import { api } from '@/lib/api'
import { formatCurrency, formatDate, relativeTime } from '@/lib/utils'
import type { InternetPackage } from '@/lib/types'

interface PackagesResp { packages: InternetPackage[] }

export default function PackagesPage() {
  const [packages, setPackages] = useState<InternetPackage[]>([])
  const [loading, setLoading] = useState(true)
  const [activeOnly, setActiveOnly] = useState(false)

  useEffect(() => {
    let mounted = true
    setLoading(true)
    api
      .get<PackagesResp>('/api/v1/packages', { query: { page_size: 100, is_active: activeOnly ? 'true' : undefined } })
      .then((res) => mounted && setPackages(res.packages))
      .finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [activeOnly])

  return (
    <div>
      <PageHeader
        title="Packages"
        subtitle="Internet service plans and pricing"
        actions={
          <>
            <div className="flex items-center gap-1 rounded-lg border border-[var(--border)] bg-[var(--bg-elev)] p-0.5 text-xs">
              <button
                onClick={() => setActiveOnly(false)}
                className={`px-3 py-1.5 rounded-md transition ${!activeOnly ? 'bg-[var(--accent-soft)] text-[var(--accent)]' : 'text-[var(--text-dim)] hover:text-[var(--text)]'}`}
              >All</button>
              <button
                onClick={() => setActiveOnly(true)}
                className={`px-3 py-1.5 rounded-md transition ${activeOnly ? 'bg-[var(--accent-soft)] text-[var(--accent)]' : 'text-[var(--text-dim)] hover:text-[var(--text)]'}`}
              >Active only</button>
            </div>
            <Link href="/packages/new" className="btn btn-primary btn-sm">+ New package</Link>
          </>
        }
      />

      {loading ? (
        <SkeletonCards count={6} />
      ) : packages.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          }
          title="No packages yet"
          description="Create your first internet service plan."
          action={<Link href="/packages/new" className="btn btn-primary btn-sm">+ Add package</Link>}
        />
      ) : (
        <div className="grid-cards">
          {packages.map((p) => (
            <div key={p.id} className="card surface-hover p-5">
              <div className="flex items-start justify-between gap-2">
                <div className="min-w-0">
                  <h3 className="font-semibold truncate">{p.name}</h3>
                  {p.description && (
                    <p className="text-xs text-[var(--text-dim)] mt-0.5 line-clamp-2">{p.description}</p>
                  )}
                </div>
                <StatusBadge status={p.is_active ? 'ACTIVE' : 'INACTIVE'} />
              </div>
              <div className="mt-4">
                <div className="text-2xl font-semibold tracking-tight">
                  {formatCurrency(p.price, p.currency)}
                </div>
                <div className="text-xs text-[var(--text-mute)]">/ {p.duration_days} days</div>
              </div>
              <div className="mt-4 grid grid-cols-2 gap-2 text-sm">
                <div className="rounded-lg border border-[var(--border)] p-2.5">
                  <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Download</div>
                  <div className="font-semibold">{p.download_mbps} Mbps</div>
                </div>
                <div className="rounded-lg border border-[var(--border)] p-2.5">
                  <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)]">Upload</div>
                  <div className="font-semibold">{p.upload_mbps} Mbps</div>
                </div>
              </div>
              {p.data_limit_gb ? (
                <div className="mt-3 text-xs text-[var(--text-dim)]">Data cap: {p.data_limit_gb} GB</div>
              ) : (
                <div className="mt-3 text-xs text-[var(--text-dim)]">Unlimited data</div>
              )}
              <div className="mt-3 pt-3 border-t border-[var(--border)] text-xs text-[var(--text-mute)] flex justify-between">
                <span>Created {relativeTime(p.created_at)}</span>
                <span>{formatDate(p.created_at)}</span>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
