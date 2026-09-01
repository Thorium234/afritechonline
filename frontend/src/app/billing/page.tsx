'use client'

import { useEffect, useState } from 'react'

interface Invoice {
  id: number
  invoice_no: string
  subscription_id: number
  amount: number
  currency: string
  status: string
  due_date: string
}

export default function BillingPage() {
  const [invoices, setInvoices] = useState<Invoice[]>([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    fetch('/api/v1/invoices')
      .then((r) => r.json())
      .then((data) => {
        setInvoices(data.data.invoices)
        setLoading(false)
      })
  }, [])

  const getStatusBadge = (status: string) => {
    switch (status) {
      case 'PAID':
        return 'badge-success'
      case 'PENDING':
        return 'badge-warning'
      case 'OVERDUE':
      case 'CANCELLED':
        return 'badge-danger'
      default:
        return 'badge-neutral'
    }
  }

  return (
    <div className="p-8">
      <div className="page-header">
        <h1 className="page-title">Billing</h1>
        <p className="page-subtitle">Invoices and payment tracking</p>
      </div>

      {loading ? (
        <div className="card overflow-hidden">
          <div className="p-8 space-y-4">
            {[1, 2, 3, 4, 5].map((i) => (
              <div key={i} className="h-12 bg-obsidian-light/50 rounded animate-pulse" />
            ))}
          </div>
        </div>
      ) : invoices.length === 0 ? (
        <div className="card p-12 text-center">
          <div className="w-12 h-12 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
            <svg className="w-6 h-6 text-slate" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
            </svg>
          </div>
          <h3 className="text-lg font-medium text-pearl mb-1">No invoices yet</h3>
          <p className="text-slate text-sm">Invoices are generated automatically when subscriptions are created.</p>
        </div>
      ) : (
        <div className="card overflow-hidden animate-fade-in">
          <table className="min-w-full">
            <thead className="bg-obsidian-dark">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Invoice</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Subscription</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Amount</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Status</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-slate uppercase tracking-wider">Due Date</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-white/5">
              {invoices.map((inv) => (
                <tr key={inv.id} className="hover:bg-white/5 transition-colors">
                  <td className="px-6 py-4 text-sm font-mono text-pearl">{inv.invoice_no}</td>
                  <td className="px-6 py-4 text-sm text-pearl">#{inv.subscription_id}</td>
                  <td className="px-6 py-4 text-sm text-pearl">
                    {inv.currency} {inv.amount.toLocaleString()}
                  </td>
                  <td className="px-6 py-4">
                    <span className={`badge ${getStatusBadge(inv.status)}`}>
                      {inv.status}
                    </span>
                  </td>
                  <td className="px-6 py-4 text-sm text-pearl">{new Date(inv.due_date).toLocaleDateString()}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
