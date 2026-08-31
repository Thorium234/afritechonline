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

  if (loading) return <div className="p-8">Loading...</div>

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6">Subscriptions</h1>
      <div className="bg-white shadow rounded">
        <table className="min-w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">ID</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Customer</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Package</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Expires</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {subscriptions.map((sub) => (
              <tr key={sub.id}>
                <td className="px-6 py-4">{sub.id}</td>
                <td className="px-6 py-4">{sub.customer_id}</td>
                <td className="px-6 py-4">{sub.package_id}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs ${sub.status === 'ACTIVE' ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'}`}>
                    {sub.status}
                  </span>
                </td>
                <td className="px-6 py-4">{new Date(sub.expiry_date).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
