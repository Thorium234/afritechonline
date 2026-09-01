'use client'

import { useState, type FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { api } from '@/lib/api'
import { ApiError } from '@/lib/types'

export default function NewCustomerPage() {
  const router = useRouter()
  const [form, setForm] = useState({ full_name: '', phone: '', email: '', username: '', status: 'ACTIVE' })
  const [error, setError] = useState<string | null>(null)
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({})
  const [loading, setLoading] = useState(false)

  const onChange = (k: keyof typeof form) => (e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement>) => {
    setForm((f) => ({ ...f, [k]: e.target.value }))
  }

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setFieldErrors({})
    setLoading(true)
    try {
      await api.post('/api/v1/customers', {
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
        setError('Failed to create customer')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title="New customer" subtitle="Add a subscriber to your ISP." />

      {error && <div className="alert alert-error mb-5">{error}</div>}

      <form onSubmit={onSubmit} className="card p-6 space-y-5">
        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Full name</label>
            <input className="input" required value={form.full_name} onChange={onChange('full_name')} placeholder="Jane Wanjiku" />
            {fieldErrors.full_name && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.full_name}</p>}
          </div>
          <div>
            <label className="label">Username</label>
            <input className="input" required value={form.username} onChange={onChange('username')} placeholder="janew" />
            {fieldErrors.username && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.username}</p>}
          </div>
        </div>
        <div className="grid sm:grid-cols-2 gap-4">
          <div>
            <label className="label">Phone</label>
            <input className="input" required value={form.phone} onChange={onChange('phone')} placeholder="+254712345678" />
            {fieldErrors.phone && <p className="text-xs text-[var(--danger)] mt-1">{fieldErrors.phone}</p>}
          </div>
          <div>
            <label className="label">Email (optional)</label>
            <input className="input" type="email" value={form.email} onChange={onChange('email')} placeholder="jane@example.com" />
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
          <button type="submit" className="btn btn-primary" disabled={loading}>
            {loading ? 'Creating…' : 'Create customer'}
          </button>
          <Link href="/customers" className="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  )
}
