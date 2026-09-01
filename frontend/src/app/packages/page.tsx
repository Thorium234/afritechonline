'use client'

import { useEffect, useState } from 'react'

interface Package {
  id: number
  name: string
  description: string
  price: number
  currency: string
  duration_days: number
  download_mbps: number
  upload_mbps: number
  is_active: boolean
}

export default function PackagesPage() {
  const [packages, setPackages] = useState<Package[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/packages')
      .then((r) => r.json())
      .then((data) => {
        setPackages(data.data.packages)
        setLoading(false)
      })
  }, [])

  return (
    <div className="p-8">
      <div className="page-header">
        <div className="flex justify-between items-start">
          <div>
            <h1 className="page-title">Packages</h1>
            <p className="page-subtitle">Internet service plans and pricing</p>
          </div>
          <a href="/packages/new" className="btn-primary">
            New Package
          </a>
        </div>
      </div>

      {loading ? (
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {[1, 2, 3].map((i) => (
            <div key={i} className="h-48 bg-obsidian-light/50 rounded-xl animate-pulse" />
          ))}
        </div>
      ) : packages.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
            <svg className="w-6 h-6 text-slate" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-pearl mb-1">No packages yet</h3>
          <p className="text-slate text-sm mb-4">Create your first internet package.</p>
          <a href="/packages/new" className="btn-primary inline-block">
            Add Package
          </a>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {packages.map((pkg, index) => (
            <div
              key={pkg.id}
              className="card p-6 animate-slide-up"
              style={{ animationDelay: `${index * 0.05}s` }}
            >
              <div className="flex items-start justify-between mb-4">
                <div>
                  <h3 className="text-lg font-semibold text-pearl">{pkg.name}</h3>
                  <p className="text-slate text-sm mt-1">{pkg.description}</p>
                </div>
                {pkg.is_active && (
                  <span className="badge badge-success">Active</span>
                )}
              </div>

              <div className="mb-4">
                <span className="text-3xl font-bold text-acid-lime">
                  {pkg.currency} {pkg.price.toLocaleString()}
                </span>
                <span className="text-slate text-sm ml-1">/ {pkg.duration_days} days</span>
              </div>

              <div className="flex items-center gap-4 text-sm text-slate">
                <div className="flex items-center gap-1.5">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  <span>{pkg.download_mbps} Mbps</span>
                </div>
                <div className="w-px h-4 bg-white/10" />
                <div className="flex items-center gap-1.5">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M13 10V3L4 14h7v7l9-11h-7z" />
                  </svg>
                  <span>{pkg.upload_mbps} Mbps</span>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  )
}
