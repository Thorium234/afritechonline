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
      .catch(() => setLoading(false))
  }, [])

  if (loading) {
    return (
      <div className="p-8">
        <div className="animate-pulse space-y-6">
          <div className="h-8 bg-obsidian-light rounded w-1/4" />
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            {[1, 2, 3, 4].map((i) => (
              <div key={i} className="h-32 bg-obsidian-light rounded-xl" />
            ))}
          </div>
        </div>
      </div>
    )
  }

  const statCards = [
    { label: 'Total Customers', value: stats?.totalCustomers || 0, color: 'text-pearl', icon: 'M12 4.354a4 4 0 014 4V8a4 4 0 01-4 4 4 4 0 01-4-4V8.354a4 4 0 014-4zM4 14.354A4 4 0 018 14.354V16a4 4 0 01-4 4 4 4 0 01-4-4v-.646zM16 14.354A4 4 0 0120 14.354V16a4 4 0 01-4 4 4 4 0 01-4-4v-.646z' },
    { label: 'Active Customers', value: stats?.activeCustomers || 0, color: 'text-acid-lime', icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z' },
    { label: 'Active Packages', value: stats?.activePackages || 0, color: 'text-pearl', icon: 'M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4' },
    { label: 'Pending Payments', value: stats?.pendingPayments || 0, color: 'text-acid-lime', icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z' },
  ]

  return (
    <div className="p-8">
      <div className="page-header">
        <h1 className="page-title">Dashboard</h1>
        <p className="page-subtitle">Overview of your ISP operations</p>
      </div>

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {statCards.map((stat, index) => (
          <div
            key={stat.label}
            className="stat-card animate-slide-up"
            style={{ animationDelay: `${index * 0.05}s` }}
          >
            <div className="flex items-start justify-between">
              <div>
                <p className="text-sm text-slate mb-1">{stat.label}</p>
                <p className={`text-3xl font-bold ${stat.color}`}>{stat.value}</p>
              </div>
              <div className="p-2 rounded-lg bg-white/5">
                <svg className="w-5 h-5 text-slate" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d={stat.icon} />
                </svg>
              </div>
            </div>
          </div>
        ))}
      </div>
    </div>
  )
}
