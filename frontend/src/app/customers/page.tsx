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

  return (
    <div className="p-8">
      <div className="page-header">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="page-title">Customers</h1>
            <p className="page-subtitle">Manage your ISP subscribers</p>
          </div>
          <a href="/customers/new" className="btn-primary">
            New Customer
          </a>
        </div>
      </div>

      {loading ? (
        <div className="card overflow-hidden">
          <div className="p-8 space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-12 bg-obsidian-light/50 rounded animate-pulse" />
            ))}
          </div>
        </div>
      ) : customers.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
            <svg className="w-6 h-6 text-slate" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 4.354a4 4 0 014 4V8a4 4 0 01-4 4 4 4 0 01-4-4V8.354a4 4 0 014-4zM4 14.354A4 4 0 018 14.354V16a4 4 0 01-4 4 4 4 0 01-4-4v-.646zM16 14.354A4 4 0 0120 14.354V16a4 4 0 01-4 4 4 4 0 01-4-4v-.646z" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-pearl mb-1">No customers yet</h3>
          <p className="text-slate text-sm mb-4">Get started by creating your first customer.</p>
          <a href="/customers/new" className="btn-primary inline-block">
            Add Customer
          </a>
        </div>
      ) : (
        <div className="card overflow-hidden animate-fade-in">
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
                  <td className="px-6 py-4">
                    <div className="text-sm font-medium text-pearl">{c.full_name}</div>
                    <div className="text-xs text-slate">@{c.username}</div>
                  </td>
                  <td className="px-6 py-4 text-sm text-pearl">{c.phone}</td>
                  <td className="px-6 py-4 text-sm text-pearl">{c.email || '-'}</td>
                  <td className="px-6 py-4">
                    <span className={`badge ${c.status === 'ACTIVE' ? 'badge-success' : 'badge-neutral'}`}>
                      {c.status}
                    </span>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
