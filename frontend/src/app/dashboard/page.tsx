'use client'

import { useEffect, useState } from 'react'

interface Stats {
  totalCustomers: number
  activeCustomers: number
  activePackages: number
  pendingPayments: number
}

export default function DashboardPage() {
  const [stats, setStats] = useState<Stats | null>(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    Promise.all([
      fetch('/api/v1/customers').then((r) => r.json()),
      fetch('/api/v1/packages?active=true').then((r) => r.json()),
      fetch('/api/v1/payments').then((r) => r.json()),
    ])
      .then(([customersData, packagesData, paymentsData]) => {
        setStats({
          totalCustomers: customersData.data?.total || customersData.data?.customers?.length || 0,
          activeCustomers: customersData.data?.customers?.filter((c: any) => c.status === 'ACTIVE').length || 0,
          activePackages: packagesData.data?.packages?.length || 0,
          pendingPayments: paymentsData.data?.payments?.filter((p: any) => p.status === 'PENDING').length || 0,
        })
        setLoading(false)
      })
  }, [])

  if (loading) return <div className="p-8 text-slate">Loading...</div>

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-8 text-pearl">Dashboard</h1>
      <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
        <div className="card p-6">
          <div className="text-sm text-slate mb-1">Total Customers</div>
          <div className="text-4xl font-bold text-pearl">{stats?.totalCustomers || 0}</div>
        </div>
        <div className="card p-6">
          <div className="text-sm text-slate mb-1">Active Customers</div>
          <div className="text-4xl font-bold text-acid-lime">{stats?.activeCustomers || 0}</div>
        </div>
        <div className="card p-6">
          <div className="text-sm text-slate mb-1">Active Packages</div>
          <div className="text-4xl font-bold text-pearl">{stats?.activePackages || 0}</div>
        </div>
        <div className="card p-6">
          <div className="text-sm text-slate mb-1">Pending Payments</div>
          <div className="text-4xl font-bold text-acid-lime">{stats?.pendingPayments || 0}</div>
        </div>
      </div>
    </div>
  )
}
