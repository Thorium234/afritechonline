'use client'

import { useEffect, useState, type FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import Link from 'next/link'
import PageHeader from '@/components/PageHeader'
import { api } from '@/lib/api'
import { ApiError } from '@/lib/types'
import { formatCurrency } from '@/lib/utils'
import type { Customer, InternetPackage } from '@/lib/types'

interface CustomersResp { customers: Customer[] }
interface PackagesResp { packages: InternetPackage[] }

export default function NewSubscriptionPage() {
  const router = useRouter()
  const [customers, setCustomers] = useState<Customer[]>([])
  const [packages, setPackages] = useState<InternetPackage[]>([])
  const [customerId, setCustomerId] = useState('')
  const [packageId, setPackageId] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)
  const [loadingData, setLoadingData] = useState(true)

  useEffect(() => {
    let mounted = true
    Promise.all([
      api.get<CustomersResp>('/api/v1/customers', { query: { page_size: 200, status: 'ACTIVE' } }),
      api.get<PackagesResp>('/api/v1/packages', { query: { is_active: 'true', page_size: 100 } }),
    ]).then(([c, p]) => {
      if (!mounted) return
      setCustomers(c.customers)
      setPackages(p.packages)
    }).finally(() => mounted && setLoadingData(false))
    return () => { mounted = false }
  }, [])

  const selectedPkg = packages.find((p) => String(p.id) === packageId)
  const selectedCust = customers.find((c) => String(c.id) === customerId)

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await api.post('/api/v1/subscriptions', {
        customer_id: parseInt(customerId, 10),
        package_id: parseInt(packageId, 10),
      })
      router.push('/subscriptions')
      router.refresh()
    } catch (err) {
      if (err instanceof ApiError) setError(err.message)
      else setError('Failed to create subscription')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="max-w-2xl">
      <PageHeader title="New subscription" subtitle="Assign a package to a customer." />
      {error && <div className="alert alert-error mb-5">{error}</div>}

      <form onSubmit={onSubmit} className="card p-6 space-y-5">
        <div>
          <label className="label">Customer</label>
          <select className="input" required value={customerId} onChange={(e) => setCustomerId(e.target.value)} disabled={loadingData}>
            <option value="">Select a customer</option>
            {customers.map((c) => (
              <option key={c.id} value={c.id}>
                {c.full_name} · @{c.username} · {c.phone}
              </option>
            ))}
          </select>
        </div>
        <div>
          <label className="label">Package</label>
          <select className="input" required value={packageId} onChange={(e) => setPackageId(e.target.value)} disabled={loadingData}>
            <option value="">Select a package</option>
            {packages.filter((p) => p.is_active).map((p) => (
              <option key={p.id} value={p.id}>
                {p.name} — {formatCurrency(p.price, p.currency)} / {p.duration_days}d
              </option>
            ))}
          </select>
        </div>

        {(selectedCust || selectedPkg) && (
          <div className="rounded-xl border border-[var(--border)] bg-[var(--bg-elev-2)] p-4 grid sm:grid-cols-2 gap-4 text-sm">
            {selectedCust && (
              <div>
                <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)] mb-1">Customer</div>
                <div className="font-medium">{selectedCust.full_name}</div>
                <div className="text-xs text-[var(--text-dim)]">{selectedCust.phone}</div>
              </div>
            )}
            {selectedPkg && (
              <div>
                <div className="text-[10px] uppercase tracking-wider text-[var(--text-mute)] mb-1">Package</div>
                <div className="font-medium">{selectedPkg.name}</div>
                <div className="text-xs text-[var(--text-dim)]">
                  {selectedPkg.download_mbps}/{selectedPkg.upload_mbps} Mbps · {selectedPkg.duration_days} days · {formatCurrency(selectedPkg.price, selectedPkg.currency)}
                </div>
              </div>
            )}
          </div>
        )}

        <div className="flex items-center gap-2 pt-2">
          <button type="submit" className="btn btn-primary" disabled={loading || loadingData}>
            {loading ? 'Creating…' : 'Create subscription'}
          </button>
          <Link href="/subscriptions" className="btn btn-secondary">Cancel</Link>
        </div>
      </form>
    </div>
  )
}
