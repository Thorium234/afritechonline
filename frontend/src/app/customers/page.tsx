'use client'

import { useEffect, useState } from 'react'

interface Customer {
  id: number
  full_name: string
  phone: string
  email: string
  username: string
  status: string
}

export default function CustomersPage() {
  const [customers, setCustomers] = useState<Customer[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetchCustomers()
  }, [])

  const fetchCustomers = async () => {
    const token = localStorage.getItem('access_token')
    const res = await fetch('/api/v1/customers', {
      headers: { Authorization: `Bearer ${token}` },
    })
    if (res.ok) {
      const data = await res.json()
      setCustomers(data.data.customers)
    }
    setLoading(false)
  }

  if (loading) return <div className="p-8 text-slate">Loading...</div>

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold text-pearl">Customers</h1>
        <a href="/customers/new" className="btn-primary">
          New Customer
        </a>
      </div>
      <div className="card overflow-hidden">
        <table className="min-w-full">
          <thead className="bg-obsidian-dark">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Phone</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-white/5">
            {customers.map((c) => (
              <tr key={c.id} className="hover:bg-white/5 transition-colors">
                <td className="px-6 py-4 text-pearl">{c.full_name}</td>
                <td className="px-6 py-4 text-pearl">{c.phone}</td>
                <td className="px-6 py-4 text-pearl">{c.email}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    c.status === 'ACTIVE'
                      ? 'bg-acid-lime/20 text-acid-lime border border-acid-lime/30'
                      : 'bg-white/5 text-slate border border-white/10'
                  }`}>
                    {c.status}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
