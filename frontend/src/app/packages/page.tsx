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

  if (loading) return <div className="p-8">Loading...</div>

  return (
    <div className="p-8">
      <div className="flex justify-between items-center mb-6">
        <h1 className="text-2xl font-bold">Packages</h1>
        <a href="/packages/new" className="bg-blue-600 text-white px-4 py-2 rounded hover:bg-blue-700">
          New Package
        </a>
      </div>
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        {packages.map((pkg) => (
          <div key={pkg.id} className="bg-white p-6 rounded shadow">
            <h3 className="text-lg font-semibold">{pkg.name}</h3>
            <p className="text-gray-600 text-sm mb-4">{pkg.description}</p>
            <div className="text-2xl font-bold text-blue-600 mb-2">
              {pkg.currency} {pkg.price.toLocaleString()}
            </div>
            <div className="text-sm text-gray-500">
              {pkg.duration_days} days | {pkg.download_mbps} Mbps / {pkg.upload_mbps} Mbps
            </div>
            {pkg.is_active && (
              <span className="inline-block mt-3 px-2 py-1 bg-green-100 text-green-800 text-xs rounded">Active</span>
            )}
          </div>
        ))}
      </div>
    </div>
  )
}
