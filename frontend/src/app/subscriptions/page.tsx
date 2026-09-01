'use client'

import { useEffect, useState } from 'react'

interface Subscription {
  id: number
  customer_id: number
  package_id: number
  start_date: string
  expiry_date: string
  status: string
  amount: number
  currency: string
}

export default function SubscriptionsPage() {
  const [subscriptions, setSubscriptions] = useState<Subscription[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/subscriptions')
      .then((r) => r.json())
      .then((data) => {
        setSubscriptions(data.data.subscriptions)
        setLoading(false)
      })
  }, [])

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'ACTIVE':
        return 'badge-success'
      case 'PENDING':
        return 'badge-warning'
      case 'EXPIRED':
      case 'CANCELLED':
        return 'badge-danger'
      default:
        return 'badge-neutral'
    }
  }

  return (
    <div className="p-8">
      <div className="page-header">
        <h1 className="page-title">Subscriptions</h1>
        <p className="page-subtitle">Customer subscription management</p>
      </div>

      {loading ? (
        <div className="card overflow-hidden">
          <div className="p-8 space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-12 bg-obsidian-light/50 rounded animate-pulse" />
            ))}
          </div>
        </div>
      ) : subscriptions.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
            <svg className="w-6 h-6 text-slate" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-pearl mb-1">No subscriptions yet</h3>
          <p className="text-slate text-sm">Create a subscription by assigning a package to a customer.</p>
        </div>
      ) : (
        <div className="card overflow-hidden animate-fade-in">
          <table className="min-w-full">
            <thead className="bg-obsidian-dark">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">ID</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Customer</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Package</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Expires</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {subscriptions.map((sub) => (
                <tr key={sub.id} className="hover:bg-white/5 transition-colors">
                  <td className="px-6 py-4 text-sm text-pearl font-mono">#{sub.id}</td>
                  <td className="px-6 py-4 text-sm text-pearl">{sub.customer_id}</td>
                  <td className="px-6 py-4 text-sm text-pearl">{sub.package_id}</td>
                  <td className="px-6 py-4">
                    <span className={`badge ${getStatusBadge(sub.status)}`}>
                      {sub.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-pearl">{new Date(sub.expiry_date).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
