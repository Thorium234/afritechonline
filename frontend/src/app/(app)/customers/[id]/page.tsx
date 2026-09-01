'use client'

import { useEffect, useState, type FormEvent } from 'react'
import { useParams, useRouter } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { api } from '@/lib/api'
import { ApiError, type Customer } from '@/lib/types'

export default function EditCustomerPage() {
  const params = useParams<{ id: string }>()
  const router = useRouter()
  const [form, setForm] = useState({ full_name: '', phone: '', email: '', username: '', status: 'ACTIVE' })
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(true)
  const [saving, setSaving] = useState(false)

  useEffect(() => {
    let mounted = true
    api
      .get<{ customer: Customer }>(`/api/v1/customers/${params.id}`)
      .then((res) => {
        if (!mounted) return
        const customer = res.customer
        setForm({
          full_name: customer.full_name || '',
          phone: customer.phone || '',
          email: customer.email || '',
          username: customer.username || '',
          status: customer.status || 'ACTIVE',
        })
      })
      .catch((err) => {
        if (mounted) setError(err instanceof Error ? err.message : 'Failed to load customer')
      })
      .finally(() => mounted && setLoading(false))
    return () => { mounted = false }
  }, [params.id])

  const onChange = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setForm((f) => ({ ...f, [k]: e.target.value }))
  }

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setFieldErrors({})
    setSaving(true)
    try {
      await api.put(`/api/v1/customers/${params.id}`, {
        ...form,
        email: form.email.trim() || undefined,
      })
      router.push('/customers')
      router.refresh()
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message)
        setFieldErrors(err.fields || {})
      } else {
        setError('Failed to update customer')
      }
    } finally {
      setSaving(false)
    }
  }

  if (loading) {
    return (
      <div className="max-w-2xl">
        <div className="card p-6 text-sm text-[var(--text-dim)]">Loading customer…</div>
      </div>
    )
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title="Edit customer" subtitle="Update subscriber details." />

      {error && <div className="alert alert-error mb-5">{error}</div>}

      <form onSubmit={onSubmit} className="card p-6 space-y-5">
        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Full name</label>
            <input className="input" required value={form.full_name} onChange={onChange('full_name')} />
            {fieldErrors.full_name && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.full_name}</p>}
          </div>
          <div>
            <label className="label">Username</label>
            <input className="input" required value={form.username} onChange={onChange('username')} />
            {fieldErrors.username && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.username}</p>}
          </div>
        </div>

        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Phone</label>
            <input className="input" required value={form.phone} onChange={onChange('phone')} />
            {fieldErrors.phone && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.phone}</p>}
          </div>
          <div>
            <label className="label">Email (optional)</label>
            <input className="input" type="email" value={form.email} onChange={onChange('email')} />
            {fieldErrors.email && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.email}</p>}
          </div>
        </div>

        <div>
          <label className="label">Status</label>
          <select className="input" value={form.status} onChange={onChange('status')}>
            <option value="ACTIVE">Active</option>
            <option value="INACTIVE">Inactive</option>
            <option value="SUSPENDED">Suspended</option>
          </select>
        </div>

        <div className="flex items-center gap-2 pt-2">
          <button type="submit" className="btn btn-primary" disabled={saving}>
            {saving ? 'Saving…' : 'Save changes'}
          </button>
          <Link href="/customers" className="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  )
}
