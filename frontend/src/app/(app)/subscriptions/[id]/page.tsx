'use client'

import { useEffect, useState } from 'react'
import { useParams } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { api } from '@/lib/api'
import { formatCurrency, formatDate, relativeTime } from '@/lib/utils'
import type { Customer, InternetPackage, Subscription } from '@/lib/types'

interface SubscriptionDetail extends Subscription {
  customer?: Customer
  package?: InternetPackage
}

export default function SubscriptionDetailPage() {
  const params = useParams<{ id: string }>()
  const [sub, setSub] = useState<SubscriptionDetail | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    let mounted = true
    api
      .get<{ subscription: Subscription }>(`/api/v1/subscriptions/${params.id}`)
      .then((res) => {
        if (!mounted) return
        setSub(res.subscription as SubscriptionDetail)
      })
      .catch((err) => {
        if (mounted) setError(err instanceof Error ? err.message : 'Failed to load subscription')
      })
      .finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [params.id])

  if (loading) {
    return (
      <div className="max-w-2xl">
        <div className="card p-6 text-sm text-[var(--text-dim)]">Loading subscription…</div>
      </div>
    )
  }

  if (error || !sub) {
    return (
      <div className="max-w-2xl">
        <PageHeader title="Subscription" subtitle="View subscription details." />
        <div className="alert alert-error">{error || 'Subscription not found'}</div>
        <Link href="/subscriptions" className="btn btn-secondary btn-sm mt-4">← Back to subscriptions</Link>
      </div>
    )
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title={`Subscription #${sub.id}`} subtitle="View subscription details." />

      <div className="grid sm:grid-cols-2 gap-4">
        <div className="card p-6">
          <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-1">Subscription ID</div>
          <div className="font-mono text-lg">#{sub.id}</div>
        </div>
        <div className="card p-6">
          <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-1">Status</div>
          <StatusBadge status={sub.status} />
        </div>
      </div>

      <div className="mt-4 card p-6 space-y-4">
        <div>
          <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-2">Amount</div>
          <div className="text-2xl font-semibold">{formatCurrency(sub.amount, sub.currency)}</div>
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-1">Start Date</div>
            <div className="text-sm font-medium">{formatDate(sub.start_date)}</div>
            <div className="text-xs text-[var(--text-mute)]">{relativeTime(sub.start_date)}</div>
          </div>
          <div>
            <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-1">Expiry Date</div>
            <div className="text-sm font-medium">{formatDate(sub.expiry_date)}</div>
            <div className="text-xs text-[var(--text-mute)]">{relativeTime(sub.expiry_date)}</div>
          </div>
        </div>

        <div className="pt-4 border-t border-[var(--border)]">
          <div className="text-xs uppercase tracking-wider text-[var(--text-mute)] mb-2">Created</div>
          <div className="text-sm text-[var(--text-dim)]">{formatDate(sub.created_at)}</div>
        </div>
      </div>

      <div className="mt-4">
        <Link href="/subscriptions" className="btn btn-secondary btn-sm">← Back to subscriptions</Link>
      </div>
    </div>
  )
}
