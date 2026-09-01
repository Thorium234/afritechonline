'use client'

import { useState, type FormEvent } from 'react'
import { useRouter } from 'next/navigation'
import { useAuth } from '@/lib/auth-context'
import { ApiError } from '@/lib/types'

export default function LoginPage() {
  const { login } = useAuth()
  const router = useRouter()
  const [identifier, setIdentifier] = useState('')
  const [password, setPassword] = useState('')
  const [showPwd, setShowPwd] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setLoading(true)
    try {
      await login(identifier.trim(), password)
      router.replace('/dashboard')
    } catch (err) {
      if (err instanceof ApiError) {
        setError(err.message)
      } else {
        setError('Unable to sign in. Please try again.')
      }
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen grid lg:grid-cols-2">
      {/* Brand panel */}
      <div className="hidden lg:flex flex-col justify-between p-12 border-r border-[var(--border)] relative overflow-hidden">
        <div className="absolute -top-32 -right-32 h-96 w-96 rounded-full bg-[var(--accent)]/10 blur-3xl" />
        <div className="absolute -bottom-32 -left-32 h-96 w-96 rounded-full bg-[#58a6ff]/10 blur-3xl" />
        <div className="relative flex items-center gap-3">
          <div className="grid h-10 w-10 place-items-center rounded-xl bg-[var(--accent-soft)] border border-[var(--accent)]/30">
            <svg className="w-5 h-5 text-[var(--accent)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
          </div>
          <div>
            <div className="text-lg font-semibold">Afritech Online</div>
            <div className="text-xs uppercase tracking-[0.18em] text-[var(--text-mute)]">ISP Management Platform</div>
          </div>
        </div>
        <div className="relative">
          <h2 className="text-3xl xl:text-4xl font-semibold tracking-tight leading-tight">
            Run your ISP with <span className="text-[var(--accent)]">clarity</span>.
          </h2>
          <p className="mt-4 text-[var(--text-dim)] max-w-md">
            Manage customers, packages, subscriptions, billing, M-Pesa payments and MikroTik routers from a single dashboard.
          </p>
          <div className="mt-8 grid grid-cols-2 gap-4 max-w-md">
            {[
              { k: 'M-Pesa', v: 'Daraja API' },
              { k: 'Network', v: 'MikroTik + RADIUS' },
              { k: 'Billing', v: 'Auto invoices' },
              { k: 'Reports', v: 'Realtime metrics' },
            ].map((it) => (
              <div key={it.k} className="rounded-xl border border-[var(--border)] bg-[var(--bg-elev)] p-3">
                <div className="text-[11px] uppercase tracking-wider text-[var(--text-mute)]">{it.k}</div>
                <div className="mt-0.5 text-sm font-medium">{it.v}</div>
              </div>
            ))}
          </div>
        </div>
        <div className="relative text-xs text-[var(--text-mute)]">
          © {new Date().getFullYear()} Afritech Online · v0.1.0
        </div>
      </div>

      {/* Form panel */}
      <div className="flex items-center justify-center p-6 sm:p-10">
        <div className="w-full max-w-sm">
          <div className="mb-8 flex items-center gap-3 lg:hidden">
            <div className="grid h-9 w-9 place-items-center rounded-xl bg-[var(--accent-soft)] border border-[var(--accent)]/30">
              <svg className="w-4 h-4 text-[var(--accent)]" fill="none" stroke="currentColor" viewBox="0 0 24 24" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
            </div>
            <div className="leading-tight">
              <div className="text-base font-semibold">Afritech Online</div>
              <div className="text-[11px] uppercase tracking-[0.18em] text-[var(--text-mute)]">ISP Management</div>
            </div>
          </div>

          <h1 className="text-2xl font-semibold tracking-tight">Welcome back</h1>
          <p className="mt-1.5 text-sm text-[var(--text-dim)]">Sign in to your admin account.</p>

          <form onSubmit={onSubmit} className="mt-8 space-y-4">
            {error && <div className="alert alert-error">{error}</div>}
            <div>
              <label htmlFor="identifier" className="label">Username or email</label>
              <input
                id="identifier"
                type="text"
                className="input"
                autoComplete="username"
                required
                value={identifier}
                onChange={(e) => setIdentifier(e.target.value)}
                placeholder="you@example.com"
              />
            </div>
            <div>
              <div className="flex items-center justify-between">
                <label htmlFor="password" className="label">Password</label>
                <button
                  type="button"
                  onClick={() => setShowPwd((v) => !v)}
                  className="text-xs text-[var(--text-dim)] hover:text-[var(--text)]"
                >
                  {showPwd ? 'Hide' : 'Show'}
                </button>
              </div>
              <input
                id="password"
                type={showPwd ? 'text' : 'password'}
                className="input"
                autoComplete="current-password"
                required
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                placeholder="••••••••"
              />
            </div>
            <button type="submit" className="btn btn-primary w-full" disabled={loading}>
              {loading ? (
                <>
                  <span className="inline-block h-2 w-2 rounded-full bg-current animate-pulse" />
                  Signing in…
                </>
              ) : (
                'Sign in'
              )}
            </button>
          </form>

          <p className="mt-6 text-center text-xs text-[var(--text-mute)]">
            Protected area. Authorized access only.
          </p>
        </div>
      </div>
    </div>
  )
}
