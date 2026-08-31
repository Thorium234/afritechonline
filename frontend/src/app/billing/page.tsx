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

  if (loading) return <div className="p-8 text-slate">Loading...</div>

  return (
    <div className="p-8">
      <h1 className="text-2xl font-bold mb-6 text-pearl">Billing</h1>
      <div className="card overflow-hidden">
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
                <td className="px-6 py-4 text-pearl">{inv.invoice_no}</td>
                <td className="px-6 py-4 text-pearl">{inv.subscription_id}</td>
                <td className="px-6 py-4 text-pearl">
                  {inv.currency} {inv.amount.toLocaleString()}
                </td>
                <td className="px-6 py-4">
                  <span className={`px-2 py-1 rounded text-xs font-medium ${
                    inv.status === 'PAID'
                      ? 'bg-acid-lime/20 text-acid-lime border border-acid-lime/30'
                      : 'bg-white/5 text-slate border border-white/10'
                  }`}>
                    {inv.status}
                  </span>
                </td>
                <td className="px-6 py-4 text-pearl">{new Date(inv.due_date).toLocaleDateString()}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  )
}
