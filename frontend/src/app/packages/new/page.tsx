'use client'

import { useState } from 'react'

export default function NewPackagePage() {
  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    currency: 'KES',
    duration_days: '30',
    download_mbps: '0',
    upload_mbps: '0',
  })
  const [message, setMessage] = useState('')

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    const token = localStorage.getItem('access_token')
    const res = await fetch('/api/v1/packages', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify({
        ...form,
        price: parseFloat(form.price),
        duration_days: parseInt(form.duration_days),
        download_mbps: parseInt(form.download_mbps),
        upload_mbps: parseInt(form.upload_mbps),
      }),
    })
    if (res.ok) {
      window.location.href = '/packages'
    } else {
      const data = await res.json()
      setMessage(data.error?.message || 'Failed to create package')
    }
  }

  return (
    <div className="p-8 max-w-2xl">
      <h1 className="text-2xl font-bold mb-6 text-pearl">New Package</h1>
      {message && (
        <div className="mb-4 p-3 rounded bg-red-900/20 border border-red-500/30 text-red-400 text-sm">
          {message}
        </div>
      )}
      <form onSubmit={handleSubmit} className="card p-6 space-y-4">
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Name</label>
          <input className="input-field" value={form.name} onChange={(e) => setForm({ ...form, name: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Description</label>
          <input className="input-field" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Price</label>
          <input type="number" className="input-field" value={form.price} onChange={(e) => setForm({ ...form, price: e.target.value })} required />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Duration (days)</label>
          <input type="number" className="input-field" value={form.duration_days} onChange={(e) => setForm({ ...form, duration_days: e.target.value })} required />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Download Mbps</label>
            <input type="number" className="input-field" value={form.download_mbps} onChange={(e) => setForm({ ...form, download_mbps: e.target.value })} />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Upload Mbps</label>
            <input type="number" className="input-field" value={form.upload_mbps} onChange={(e) => setForm({ ...form, upload_mbps: e.target.value })} />
          </div>
        </div>
        <button type="submit" className="btn-primary">
          Create Package
        </button>
      </form>
    </div>
  )
}
