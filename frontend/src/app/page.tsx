'use client'

import { useEffect } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/lib/auth-context'

export default function HomePage() {
  const router = useRouter()
  const { user, loading } = useAuth()
  useEffect(() => {
    if (loading) return
    router.replace(user ? '/dashboard' : '/login')
  }, [user, loading, router])
  return (
    <div className="min-h-screen grid place-items-center text-[var(--text-dim)]">
      <div className="flex items-center gap-3">
        <span className="inline-block h-2 w-2 rounded-full bg-[var(--accent)] animate-pulse" />
        Loading Afritech Online…
      </div>
    </div>
  )
}
