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
  const [loading, setLoading] = useState(false)

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setMessage('')
    setLoading(true)

    try {
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
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="p-8 max-w-2xl">
      <div className="page-header">
        <h1 className="page-title">New Package</h1>
        <p className="page-subtitle">Create a new internet service plan</p>
      </div>

      {message && (
        <div className="mb-6 p-3 rounded-lg bg-red-900/20 border border-red-500/30 text-red-400 text-sm animate-fade-in">
          {message}
        </div>
      )}

      <form onSubmit={handleSubmit} className="card p-6 space-y-5">
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Package Name</label>
          <input
            className="input-field"
            value={form.name}
            onChange={(e) => setForm({ ...form, name: e.target.value })}
            placeholder="Home Basic"
            required
          />
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Description</label>
          <input
            className="input-field"
            value={form.description}
            onChange={(e) => setForm({ ...form, description: e.target.value })}
            placeholder="10 Mbps, 30 days"
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Price</label>
            <input
              type="number"
              className="input-field"
              value={form.price}
              onChange={(e) => setForm({ ...form, price: e.target.value })}
              placeholder="1000"
              required
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Currency</label>
            <select
              className="input-field"
              value={form.currency}
              onChange={(e) => setForm({ ...form, currency: e.target.value })}
            >
              <option value="KES">KES</option>
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
            </select>
          </div>
        </div>
        <div>
          <label className="block text-sm font-medium mb-2 text-pearl">Duration (days)</label>
          <input
            type="number"
            className="input-field"
            value={form.duration_days}
            onChange={(e) => setForm({ ...form, duration_days: e.target.value })}
            required
          />
        </div>
        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Download Mbps</label>
            <input
              type="number"
              className="input-field"
              value={form.download_mbps}
              onChange={(e) => setForm({ ...form, download_mbps: e.target.value })}
            />
          </div>
          <div>
            <label className="block text-sm font-medium mb-2 text-pearl">Upload Mbps</label>
            <input
              type="number"
              className="input-field"
              value={form.upload_mbps}
              onChange={(e) => setForm({ ...form, upload_mbps: e.target.value })}
            />
          </div>
        </div>
        <div className="flex gap-3 pt-2">
          <button type="submit" className="btn-primary" disabled={loading}>
            {loading ? 'Creating...' : 'Create Package'}
          </button>
          <a href="/packages" className="btn-secondary">
            Cancel
          </a>
        </div>
      </form>
    </div>
  )
}
