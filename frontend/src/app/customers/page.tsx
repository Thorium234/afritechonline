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

  if (loading) return <div className="p-8">Loading...</div>

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Customers</h1>
        <a href="/customers/new" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          New Customer
        </a>
      </div>
      <div className="bg-white shadow rounded">
        <table className="min-w-full">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Phone</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Email</th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase">Status</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-gray-200">
            {customers.map((c) => (
              <tr key={c.id}>
                <td className="px-6 py-4">{c.full_name}</td>
                <td className="px-6 py-4">{c.phone}</td>
                <td className="px-6 py-4">{c.email}</td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs ${c.status === 'ACTIVE' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}`}>
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
