'use client'

import { useState, type FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { api } from '@/lib/api'
import { ApiError } from '@/lib/types'

export default function NewPackagePage() {
  const router = useRouter()
  const [form, setForm] = useState({
    name: '',
    description: '',
    price: '',
    currency: 'KES',
    duration_days: '30',
    download_mbps: '10',
    upload_mbps: '5',
    data_limit_gb: '',
  })
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(false)

  const set = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) =>
    setForm((f) => ({ ...f, [k]: e.target.value }))

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setFieldErrors({})
    setLoading(true)
    try {
      const payload = {
        name: form.name.trim(),
        description: form.description.trim(),
        price: parseFloat(form.price),
        currency: form.currency,
        duration_days: parseInt(form.duration_days, 10),
        download_mbps: parseInt(form.download_mbps, 10),
        upload_mbps: parseInt(form.upload_mbps, 10),
        data_limit_gb: form.data_limit_gb ? parseInt(form.data_limit_gb, 10) : undefined,
        is_active: true,
      }
      await api.post('/api/v1/packages', payload)
      router.push('/packages')
      router.refresh()
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message)
        setFieldErrors(err.fields || {})
      } else {
        setError('Failed to create package')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title="New package" subtitle="Create an internet service plan." />
      {error && <div className="alert alert-error mb-5">{error}</div>}
      <form onSubmit={onSubmit} className="card p-6 space-y-5">
        <div>
          <label className="label">Name</label>
          <input className="input" required value={form.name} onChange={set('name')} placeholder="Home Basic 10Mbps" />
        </div>
        <div>
          <label className="label">Description</label>
          <input className="input" value={form.description} onChange={set('description')} placeholder="10 Mbps download, 30 days unlimited" />
        </div>
        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Price</label>
            <input className="input" type="number" min="0" step="0.01" required value={form.price} onChange={set('price')} placeholder="1500" />
          </div>
          <div>
            <label className="label">Currency</label>
            <select className="input" value={form.currency} onChange={set('currency')}>
              <option value="KES">KES</option>
              <option value="USD">USD</option>
              <option value="EUR">EUR</option>
              <option value="UGX">UGX</option>
              <option value="TZS">TZS</option>
            </select>
          </div>
        </div>
        <div className="grid sm:grid-cols-3 gap-4">
          <div>
            <label className="label">Duration (days)</label>
            <input className="input" type="number" min="1" required value={form.duration_days} onChange={set('duration_days')} />
          </div>
          <div>
            <label className="label">Download (Mbps)</label>
            <input className="input" type="number" min="0" required value={form.download_mbps} onChange={set('download_mbps')} />
          </div>
          <div>
            <label className="label">Upload (Mbps)</label>
            <input className="input" type="number" min="0" required value={form.upload_mbps} onChange={set('upload_mbps')} />
          </div>
        </div>
        <div>
          <label className="label">Data limit (GB, optional)</label>
          <input className="input" type="number" min="0" value={form.data_limit_gb} onChange={set('data_limit_gb')} placeholder="Leave empty for unlimited" />
        </div>
        <div className="flex items-center gap-2 pt-2">
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Creating…' : 'Create package'}
          </button>
          <Link href="/packages" className="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  )
}
