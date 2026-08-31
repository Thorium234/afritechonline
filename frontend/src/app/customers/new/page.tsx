'use client'

import { useState } from 'react'

export default function NewCustomerPage() {
  const [form, setForm] = useState({ full_name: '', phone: '', email: '', username: '' })
  const [message, setMessage] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = localStorage.getItem('access_token')
    const res = await fetch('/api/v1/customers', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(form),
    })
    if (res.ok) {
      window.location.href = '/customers'
    } else {
      const data = await res.json()
      setMessage(data.error?.message || 'Failed to create customer')
    }
  }

  return (
    <div className="p-8 max-w-2xl">
      <h1 className="text-2xl font-bold mb-6">New Customer</h1>
      {message && <p className="text-red-600 mb-4">{message}</p>}
      <form onSubmit={handleSubmit} className="space-y-4 bg-white p-6 rounded shadow">
        <div>
          <label className="block text-sm font-medium mb-1">Full Name</label>
          <input className="w-full border rounded p-2" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Phone</label>
          <input className="w-full border rounded p-2" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Email</label>
          <input className="w-full border rounded p-2" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium mb-1">Username</label>
          <input className="w-full border rounded p-2" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} required />
        </div>
        <button type="submit" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          Create Customer
        </button>
      </form>
    </div>
  )
}
