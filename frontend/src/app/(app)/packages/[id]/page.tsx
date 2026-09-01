'use client'

import { useEffect, useState, type FormEvent } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { api } from '@/lib/api'
import { ApiError, type InternetPackage } from '@/lib/types'

export default function EditPackagePage() {
  const params = useParams<{ id: string }>()
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
    is_active: true,
  })
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    let mounted = true
    api
      .get<{ package: InternetPackage }>(`/api/v1/packages/${params.id}`)
      .then((res) => {
        if (!mounted) return
        const pkg = res.package
        setForm({
          name: pkg.name || '',
          description: pkg.description || '',
          price: String(pkg.price),
          currency: pkg.currency || 'KES',
          duration_days: String(pkg.duration_days),
          download_mbps: String(pkg.download_mbps),
          upload_mbps: String(pkg.upload_mbps),
          data_limit_gb: pkg.data_limit_gb ? String(pkg.data_limit_gb) : '',
          is_active: pkg.is_active,
        })
      })
      .catch((err) => {
        if (mounted) setError(err instanceof Error ? err.message : 'Failed to load package')
      })
      .finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [params.id])

  const set = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    if (k === 'is_active') {
      setForm((f) => ({ ...f, [k]: (e.target as HTMLInputElement).checked }))
    } else {
      setForm((f) => ({ ...f, [k]: e.target.value }))
    }
  }

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setFieldErrors({})
    setSaving(true)
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
        is_active: form.is_active,
      }
      await api.put(`/api/v1/packages/${params.id}`, payload)
      router.push('/packages')
      router.refresh()
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message)
        setFieldErrors(err.fields || {})
      } else {
        setError('Failed to update package')
      }
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl">
        <div className="card p-6 text-sm text-[var(--text-dim)]">Loading package…</div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title="Edit package" subtitle="Update internet service plan details." />

      {error && <div className="alert alert-error mb-5">{error}</div>}

      <form onSubmit={onSubmit} className="card p-6 space-y-5">
        <div>
          <label className="label">Name</label>
          <input className="input" required value={form.name} onChange={set('name')} />
        </div>

        <div>
          <label className="label">Description</label>
          <input className="input" value={form.description} onChange={set('description')} />
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Price</label>
            <input className="input" type="number" min="0" step="0.01" required value={form.price} onChange={set('price')} />
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
          <input className="input" type="number" min="0" value={form.data_limit_gb} onChange={set('data_limit_gb')} />
        </div>

        <div className="flex items-center gap-2">
          <input type="checkbox" checked={form.is_active} onChange={set('is_active')} id="is_active" className="checkbox" />
          <label htmlFor="is_active" className="cursor-pointer text-sm">Active</label>
        </div>

        <div className="flex items-center gap-2 pt-2">
          <button type="submit" className="btn btn-primary" disabled={saving}>
            {saving ? 'Saving…' : 'Save changes'}
          </button>
          <Link href="/packages" className="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  )
}
