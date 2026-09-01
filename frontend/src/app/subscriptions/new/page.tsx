'use client'

import { useEffect, useState } from 'react'

interface Customer {
  id: number
  full_name: string
  username: string
}

interface Package {
  id: number
  name: string
  price: number
  duration_days: number
}

export default function NewSubscriptionPage() {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [packages, setPackages] = useState<Package[]>([])
  const [customerId, setCustomerId] = useState('')
  const [packageId, setPackageId] = useState('')
  const [message, setMessage] = useState('')
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    const token = localStorage.getItem('access_token')
    Promise.all([
      fetch('/api/v1/customers', { headers: { Authorization: `Bearer ${token}` } }).then((r) => r.json()),
      fetch('/api/v1/packages?active=true').then((r) => r.json()),
    ]).then(([customersData, packagesData]) => {
      setCustomers(customersData.data?.customers || [])
      setPackages(packagesData.data?.packages || [])
    })
  }, [])

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage('')
    setLoading(true)

    try {
      const token = localStorage.getItem('access_token')
      const res = await fetch('/api/v1/subscriptions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({ customer_id: parseInt(customerId), package_id: parseInt(packageId) }),
      })
      if (res.ok) {
        window.location.href = '/subscriptions'
      } else {
        const data = await res.json()
        setMessage(data.error?.message || 'Failed to create subscription')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-8 max-w-2xl">
      <div className="page-header">
        <h1 className="page-title">New Subscription</h1>
        <p className="page-subtitle">Assign a package to a customer</p>
      </div>

      {message && (
        <div className="mb-6 p-3 rounded-lg bg-red-900/20 border border-red-500/30 text-red-400 text-sm animate-fade-in">
          {message}
        </div>
      )}

      <form onSubmit={handleSubmit} className="card p-6 space-y-5">
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Customer</label>
          <select
            className="input-field"
            value={customerId}
            onChange={(e) => setCustomerId(e.target.value)}
            required
          >
            <option value="">Select a customer</option>
            {customers.map((c) => (
              <option key={c.id} value={c.id}>{c.full_name} (@{c.username})</option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Package</label>
          <select
            className="input-field"
            value={packageId}
            onChange={(e) => setPackageId(e.target.value)}
            required
          >
            <option value="">Select a package</option>
            {packages.map((p) => (
              <option key={p.id} value={p.id}>{p.name} - {p.currency} {p.price} ({p.duration_days} days)</option>
            ))}
          </select>
        </div>
        <div className="flex gap-3 pt-2">
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Creating...' : 'Create Subscription'}
          </button>
          <a href="/subscriptions" className="btn-secondary">
            Cancel
          </a>
        </div>
      </form>
    </div>
  )
}
