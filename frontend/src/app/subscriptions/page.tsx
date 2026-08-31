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

  if (loading) return <div className="p-8 text-slate">Loading...</div>

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6 text-pearl">Subscriptions</h1>
      <div className="card overflow-hidden">
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
                <td className="px-6 py-4 text-pearl">{sub.id}</td>
                <td className="px-6 py-4 text-pearl">{sub.customer_id}</td>
                <td className="px-6 py-4 text-pearl">{sub.package_id}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    sub.status === 'ACTIVE'
                      ? 'bg-acid-lime/20 text-acid-lime border border-acid-lime/30'
                      : 'bg-white/5 text-slate border border-white/10'
                  }`}>
                    {sub.status}
                  </span>
                </td>
                <td className="px-6 py-4 text-pearl">{new Date(sub.expiry_date).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
