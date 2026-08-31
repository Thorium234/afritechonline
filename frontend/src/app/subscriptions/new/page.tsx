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
  }

  return (
    <div className="p-8 max-w-2xl">
      <h1 className="text-2xl font-bold mb-6">New Subscription</h1>
      {message && <p className="text-red-600 mb-4">{message}</p>}
      <form onSubmit={handleSubmit} className="space-y-4 bg-white p-6 rounded shadow">
        <div>
          <label className="block text-sm font-medium mb-1">Customer</label>
          <select className="w-full border rounded p-2" value={customerId} onChange={(e) => setCustomerId(e.target.value)} required>
            <option value="">Select customer</option>
            {customers.map((c) => (
              <option key={c.id} value={c.id}>{c.full_name} ({c.username})</option>
            ))}
          </select>
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Package</label>
          <select className="w-full border rounded p-2" value={packageId} onChange={(e) => setPackageId(e.target.value)} required>
            <option value="">Select package</option>
            {packages.map((p) => (
              <option key={p.id} value={p.id}>{p.name} - {p.currency} {p.price} ({p.duration_days} days)</option>
            ))}
          </select>
        </div>
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          Create Subscription
        </button>
      </form>
    </div>
  )
}
