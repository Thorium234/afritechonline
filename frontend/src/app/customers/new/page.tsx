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
      <h1 className="text-2xl font-bold mb-6 text-pearl">New Customer</h1>
      {message && (
        <div className="mb-4 p-3 rounded bg-red-900/20 border border-red-500/30 text-red-400 text-sm">
          {message}
        </div>
      )}
      <form onSubmit={handleSubmit} className="card p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Full Name</label>
          <input className="input-field" value={form.full_name} onChange={(e) => setForm({ ...form, full_name: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Phone</label>
          <input className="input-field" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Email</label>
          <input className="input-field" value={form.email} onChange={(e) => setForm({ ...form, email: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Username</label>
          <input className="input-field" value={form.username} onChange={(e) => setForm({ ...form, username: e.target.value })} required />
        </div>
        <button type="submit" className="btn-primary">
          Create Customer
        </button>
      </form>
    </div>
  )
}
