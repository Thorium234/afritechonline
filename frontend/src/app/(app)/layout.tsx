'use client'

import { useState } from 'react'
import Sidebar from '@/components/Sidebar'
import Topbar from '@/components/Topbar'
import { useRequireAuth } from '@/lib/auth-context'

export default function AppLayout({ children }: { children: React.ReactNode }) {
  const auth = useRequireAuth()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  if (auth.loading) {
    return (
      <div className="min-h-screen grid place-items-center text-[var(--text-dim)]">
        <div className="flex items-center gap-3">
          <span className="inline-block h-2 w-2 rounded-full bg-[var(--accent)] animate-pulse" />
          Loading…
        </div>
      </div>
    )
  }

  if (!auth.user) {
    return null
  }

  return (
    <div className="min-h-screen flex">
      <Sidebar open={sidebarOpen} onClose={() => setSidebarOpen(false)} />
      <div className="flex-1 min-w-0 flex flex-col">
        <Topbar onMenu={() => setSidebarOpen(true)} />
        <main className="flex-1 p-4 sm:p-6 lg:p-8 max-w-[1400px] w-full mx-auto animate-fade-in">
          {children}
        </main>
      </div>
    </div>
  )
}
