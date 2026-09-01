'use client'

import { useEffect, useState } from 'react'
import PageHeader from '@/components/PageHeader'
import StatusBadge from '@/components/StatusBadge'
import { SkeletonTable } from '@/components/Skeleton'
import EmptyState from '@/components/EmptyState'
import { api } from '@/lib/api'
import { formatDate } from '@/lib/utils'
import type { Router } from '@/lib/types'

interface RoutersResp { routers: Router[] }

export default function RoutersPage() {
  const [routers, setRouters] = useState<Router[]>([])
  const [loading, setLoading] = useState(true)
  const [testing, setTesting] = useState<number | null>(null)

  const load = () => {
    setLoading(true)
    api.get<RoutersResp>('/api/v1/routers', { query: { page_size: 100 } })
      .then((r) => setRouters(r.routers))
      .finally(() => setLoading(false))
  }

  useEffect(() => { load() }, [])

  const test = async (id: number) => {
    setTesting(id)
    try {
      await api.post(`/api/v1/routers/${id}/test`)
    } catch {
      // ignore — UI will reflect on next refresh
    } finally {
      setTesting(null)
      load()
    }
  }

  return (
    <div>
      <PageHeader title="Routers" subtitle="MikroTik devices managed by the platform" />
      {loading ? (
        <SkeletonTable rows={4} />
      ) : routers.length === 0 ? (
        <EmptyState
          icon={
            <svg className="w-6 h-6 text-[var(--text-mute)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={1.6}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M8.111 16.404a5.5 5.5 0 017.778 0M12 20h.01M2 8.82a15 15 0 0120 0M5 12.859a10 10 0 0114 0" />
            </svg>
          }
          title="No routers registered"
          description="Register MikroTik routers from the backend API to manage them here."
        />
      ) : (
        <div className="card overflow-hidden">
          <div className="overflow-x-auto">
            <table className="table">
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Host</th>
                  <th>Location</th>
                  <th>Status</th>
                  <th>Created</th>
                  <th />
                </tr>
              </thead>
              <tbody>
                {routers.map((r) => (
                  <tr key={r.id}>
                    <td className="font-medium">{r.name}</td>
                    <td className="font-mono text-sm">{r.host}:{r.api_port}</td>
                    <td className="text-[var(--text-dim)]">{r.location || '—'}</td>
                    <td><StatusBadge status={r.status} /></td>
                    <td className="text-[var(--text-dim)] text-sm">{formatDate(r.created_at)}</td>
                    <td>
                      <button
                        className="btn btn-secondary btn-sm"
                        disabled={testing === r.id}
                        onClick={() => test(r.id)}
                      >
                        {testing === r.id ? 'Testing…' : 'Test'}
                      </button>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </div>
      )}
    </div>
  )
}
